package member

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"exam/internal/consts"
	"exam/internal/model/bo"

	"github.com/gogf/gf/v2/errors/gerror"
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

func (s *sMember) MemberImport(ctx context.Context, r io.Reader, creator string) (*bo.MemberImportResult, error) {
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

	out := &bo.MemberImportResult{Errors: make([]string, 0, 8)}
	for i := 1; i < len(records); i++ {
		lineNo := i + 1
		row := records[i]
		if memberImportRowEmpty(row) {
			continue
		}
		out.Total++
		username := memberImportCell(row, colIdx["username"])
		password := memberImportCell(row, colIdx["password"])
		nickname := memberImportCell(row, colIdx["nickname"])
		email := memberImportCell(row, colIdx["email"])
		mobile := memberImportCell(row, colIdx["mobile"])
		statusStr := memberImportCell(row, colIdx["status"])
		status := memberImportParseStatus(statusStr)

		if username == "" || password == "" {
			out.Failed++
			memberImportAppendErr(out, lineNo, "用户名与密码不能为空")
			continue
		}

		_, err := s.MemberCreate(ctx, username, password, nickname, email, mobile, creator, status)
		if err != nil {
			out.Failed++
			memberImportAppendErr(out, lineNo, memberImportFormatRowErr(username, err))
			continue
		}
		out.Success++
	}
	return out, nil
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
	if _, ok := idx["username"]; !ok {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if _, ok := idx["password"]; !ok {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	return idx, nil
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

func memberImportParseStatus(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return consts.StatusNormal
	}
	switch strings.ToLower(s) {
	case "0", "正常", "启用":
		return consts.StatusNormal
	case "1", "停用", "禁用":
		return consts.StatusDisabled
	default:
		return consts.StatusNormal
	}
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
