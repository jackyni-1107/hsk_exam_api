package member

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	sysentity "exam/internal/model/entity/sys"
	secsvc "exam/internal/service/security"
	notisvc "exam/internal/service/sysnotification"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const memberImportMaxRows = 2000
const memberImportMaxErrLines = 100

var memberImportHeaderKeys = map[string]string{
	"username": "username",
	"用户名":      "username",
	"password": "password",
	"密码":       "password",
	"nickname": "nickname",
	"昵称":       "nickname",
	"email":    "email",
	"邮箱":       "email",
	"mobile":   "mobile",
	"手机":       "mobile",
	"status":   "status",
	"状态":       "status",
}

func (s *sMember) MemberImport(ctx context.Context, r io.Reader, creator string, country, year string, seqDigits int, opts bo.MemberImportOptions) (*bo.MemberImportResult, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if bytes.HasPrefix(raw, []byte{0xEF, 0xBB, 0xBF}) {
		raw = raw[3:]
	}

	cr := csv.NewReader(bytes.NewReader(raw))
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true
	records, err := cr.ReadAll()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "CSV 解析失败")
	}
	if len(records) < 2 {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	dataRowCount := len(records) - 1
	if dataRowCount > memberImportMaxRows {
		return nil, gerror.NewCodef(consts.CodeInvalidParams, "单次最多导入 %d 条", memberImportMaxRows)
	}

	colIdx, err := memberImportParseHeader(records[0])
	if err != nil {
		return nil, err
	}
	pwdCol := memberImportCol(colIdx, "password")
	nickCol := memberImportCol(colIdx, "nickname")
	emailCol := memberImportCol(colIdx, "email")
	mobileCol := memberImportCol(colIdx, "mobile")

	out := &bo.MemberImportResult{Errors: make([]string, 0, 8)}
	opts = memberImportNormalizeOptions(opts)
	usernameRule, err := memberImportValidateUsernameParams(country, year, seqDigits)
	if err != nil {
		return nil, err
	}
	nextSeq, err := s.memberImportLoadNextSeq(ctx, usernameRule)
	if err != nil {
		return nil, err
	}
	seenEmails := make(map[string]int, len(records))
	for i := 1; i < len(records); i++ {
		lineNo := i + 1
		row := records[i]
		if memberImportRowEmpty(row) {
			continue
		}
		out.Total++
		password := memberImportCell(row, pwdCol)
		nickname := memberImportCell(row, nickCol)
		email := memberImportCell(row, emailCol)
		mobile := memberImportCell(row, mobileCol)

		if nickname == "" || email == "" {
			out.Failed++
			memberImportAppendErr(out, lineNo, "昵称与邮箱不能为空")
			continue
		}
		if password == "" {
			pwd, genErr := memberImportBuildDefaultPassword(ctx, email, opts)
			if genErr != nil {
				out.Failed++
				memberImportAppendErr(out, lineNo, genErr.Error())
				continue
			}
			password = pwd
		}
		emailKey := strings.ToLower(email)
		if prevLine, ok := seenEmails[emailKey]; ok {
			out.Failed++
			memberImportAppendErr(out, lineNo, fmt.Sprintf("邮箱「%s」在第%d行已出现", email, prevLine))
			continue
		}
		exists, err := s.memberImportEmailExists(ctx, email)
		if err != nil {
			out.Failed++
			memberImportAppendErr(out, lineNo, memberImportFormatRowErr("", err))
			continue
		}
		if exists {
			out.Failed++
			memberImportAppendErr(out, lineNo, fmt.Sprintf("邮箱「%s」已存在", email))
			continue
		}
		username := memberImportBuildUsername(usernameRule, nextSeq)
		_, err = s.MemberCreate(ctx, username, password, nickname, email, mobile, creator, consts.StatusNormal)
		if err != nil {
			out.Failed++
			memberImportAppendErr(out, lineNo, memberImportFormatRowErr(username, err))
			continue
		}
		seenEmails[emailKey] = lineNo
		nextSeq++
		out.Success++

		if opts.SendPasswordNotice {
			memberImportNotifyPasswordAsync(ctx, username, email, password)
		}
	}
	return out, nil
}

func memberImportNormalizeOptions(opts bo.MemberImportOptions) bo.MemberImportOptions {
	opts.EmailPasswordPickPositions = strings.TrimSpace(opts.EmailPasswordPickPositions)
	if opts.EmailPasswordPickPositions == "" {
		opts.EmailPasswordPickPositions = "1,3,5"
	}
	opts.FixedPasswordSuffix = strings.TrimSpace(opts.FixedPasswordSuffix)
	if opts.FixedPasswordSuffix == "" {
		opts.FixedPasswordSuffix = "hskmock"
	}
	return opts
}

func memberImportBuildDefaultPassword(ctx context.Context, email string, opts bo.MemberImportOptions) (string, error) {
	if opts.UseRandomPassword {
		cfg := secsvc.Security().LoadPasswordCfg(ctx)
		password, err := utility.GeneratePasswordByPolicy(cfg)
		if err != nil {
			return "", err
		}
		if err = secsvc.Security().ValidatePasswordPolicy(ctx, password); err != nil {
			return "", gerror.NewCode(consts.CodePasswordWeak)
		}
		return password, nil
	}
	picks, err := memberImportParsePickPositions(opts.EmailPasswordPickPositions)
	if err != nil {
		return "", err
	}
	pwd, err := memberImportPasswordFromEmail(email, opts.FixedPasswordSuffix, picks)
	if err != nil {
		return "", gerror.NewCodef(consts.CodeInvalidParams, "邮箱长度不足，无法按固定规则位次生成默认密码，请填写密码列")
	}
	return pwd, nil
}

func memberImportParseHeader(header []string) (map[string]int, error) {
	idx := make(map[string]int)
	for i, h := range header {
		key := strings.TrimSpace(h)
		if key == "" {
			continue
		}
		if canon, ok := memberImportHeaderKeys[strings.ToLower(key)]; ok {
			idx[canon] = i
			continue
		}
		if canon, ok := memberImportHeaderKeys[key]; ok {
			idx[canon] = i
		}
	}
	if _, ok := idx["nickname"]; !ok {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if _, ok := idx["email"]; !ok {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	return idx, nil
}

// memberImportCol 读取列下标；表头未包含该列时返回 -1，避免 map 缺键时误用 0。
func memberImportCol(idx map[string]int, key string) int {
	if i, ok := idx[key]; ok {
		return i
	}
	return -1
}

// memberImportPasswordFromEmail 按配置位次从邮箱提取字符（1-based）并与固定后缀组成默认密码。
func memberImportPasswordFromEmail(email, fixedSuffix string, picks []int) (string, error) {
	runes := []rune(strings.TrimSpace(email))
	if len(runes) == 0 {
		return "", gerror.NewCode(consts.CodeInvalidParams)
	}
	if len(picks) == 0 {
		return "", gerror.NewCode(consts.CodeInvalidParams)
	}
	fixedSuffix = strings.TrimSpace(fixedSuffix)
	if fixedSuffix == "" {
		fixedSuffix = "hskmock"
	}
	var b strings.Builder
	for _, p := range picks {
		i := p - 1
		if i < 0 || i >= len(runes) {
			return "", gerror.NewCode(consts.CodeInvalidParams)
		}
		b.WriteRune(runes[i])
	}
	b.WriteString("@")
	b.WriteString(fixedSuffix)
	return b.String(), nil
}

func memberImportParsePickPositions(raw string) ([]int, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return []int{1, 3, 5}, nil
	}
	replacer := strings.NewReplacer("，", ",", "；", ",", ";", ",", " ", ",")
	s = replacer.Replace(s)
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	seen := make(map[int]struct{}, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 {
			return nil, gerror.NewCodef(consts.CodeInvalidParams, "邮箱取位规则非法：%s", raw)
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	if len(out) == 0 {
		return nil, gerror.NewCodef(consts.CodeInvalidParams, "邮箱取位规则不能为空")
	}
	return out, nil
}

func memberImportNotifyPasswordAsync(ctx context.Context, username, email, password string) {
	recipient := strings.TrimSpace(email)
	if recipient == "" {
		return
	}
	go func(user, rcpt, pwd string) {
		notifyCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		vars, _ := json.Marshal(map[string]string{
			"username": user,
			"password": pwd,
		})
		if _, err := notisvc.SysNotification().Send(notifyCtx, "forget_password", "email", rcpt, string(vars)); err != nil {
			g.Log().Warningf(notifyCtx, "member import async notify failed username=%s email=%s: %v", user, rcpt, err)
		}
	}(username, recipient, password)
}

func memberImportRowEmpty(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

func memberImportCell(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[col])
}

func memberImportAppendErr(out *bo.MemberImportResult, line int, msg string) {
	if len(out.Errors) >= memberImportMaxErrLines {
		return
	}
	out.Errors = append(out.Errors, fmt.Sprintf("第%d行：%s", line, msg))
}

func memberImportFormatRowErr(username string, err error) string {
	if err == nil {
		return "未知错误"
	}
	if c := gerror.Code(err); c != nil {
		switch c.Code() {
		case consts.CodeMemberExists.Code():
			return fmt.Sprintf("用户名「%s」已存在", username)
		case consts.CodePasswordWeak.Code():
			return "密码不符合安全策略"
		case consts.CodeInvalidParams.Code():
			return "参数无效（请检查用户名等）"
		default:
			if m := c.Message(); m != "" {
				return m
			}
		}
	}
	return err.Error()
}

type memberImportUsernameRule struct {
	country string
	year    string
	digits  int
}

func memberImportValidateUsernameParams(country, year string, seqDigits int) (*memberImportUsernameRule, error) {
	country = strings.TrimSpace(country)
	year = strings.TrimSpace(year)
	if country == "" || year == "" {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if seqDigits < 1 {
		return nil, gerror.NewCodef(consts.CodeInvalidParams, "序号位数必须大于等于1")
	}
	return &memberImportUsernameRule{
		country: strings.ToUpper(country),
		year:    year,
		digits:  seqDigits,
	}, nil
}

func (s *sMember) memberImportLoadNextSeq(ctx context.Context, rule *memberImportUsernameRule) (int, error) {
	// 单次查询：取当前规则前缀下已有用户名，在内存中解析最大序号（勿用 Scan(&[]string)，多行结果在 ORM 中不可靠）
	prefix := fmt.Sprintf("%s%s-", rule.country, rule.year)
	var rows []sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).
		Fields("username").
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		WhereLike("username", prefix+"%").
		Scan(&rows); err != nil {
		return 0, err
	}
	maxSeq := 0
	for _, row := range rows {
		seq, ok := memberImportParseUsernameSeq(row.Username, rule)
		if ok && seq > maxSeq {
			maxSeq = seq
		}
	}
	return maxSeq + 1, nil
}

func memberImportBuildUsername(rule *memberImportUsernameRule, seq int) string {
	return fmt.Sprintf("%s%s-%0*d", rule.country, rule.year, rule.digits, seq)
}

// memberImportParseUsernameSeq 仅识别「连字符后恰好 rule.digits 位数字」的用户名，使不同位数互不影响序号。
// 例如 5 位 TH2026-00004 与 6 位 TH2026-000004 分属两套序列；切换位数且库中无同位数记录时从 1 起号。
func memberImportParseUsernameSeq(username string, rule *memberImportUsernameRule) (int, bool) {
	if rule == nil || rule.digits < 1 {
		return 0, false
	}
	pattern := fmt.Sprintf(
		`^%s%s-(\d{%d})$`,
		regexp.QuoteMeta(rule.country),
		regexp.QuoteMeta(rule.year),
		rule.digits,
	)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(strings.TrimSpace(username))
	if len(matches) != 2 {
		return 0, false
	}
	seq, err := strconv.Atoi(matches[1])
	if err != nil || seq < 1 {
		return 0, false
	}
	return seq, true
}

func (s *sMember) memberImportEmailExists(ctx context.Context, email string) (bool, error) {
	cnt, err := dao.SysMember.Ctx(ctx).
		Where("email", email).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
