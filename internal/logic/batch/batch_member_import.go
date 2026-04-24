package batch

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"exam/internal/consts"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	batchMemberImportMaxRows     = 5000
	batchMemberImportMaxErrLines = 100
)

var batchMemberImportHeaderKeys = map[string]string{
	"username":  "username",
	"用户名":       "username",
	"member_id": "member_id",
	"会员id":      "member_id",
	"会员ID":      "member_id",
}

type batchMemberImportRow struct {
	lineNo   int
	memberID int64
	username string
}

// ExamBatchMembersImportByCSV 按模板 CSV 向指定批次和试卷导入学员。
func (s *sBatch) ExamBatchMembersImportByCSV(ctx context.Context, batchID int64, examPaperID int64, r io.Reader, creator string) (total, success, failed, inserted int, errors []string, err error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, 0, 0, nil, err
	}
	if len(raw) == 0 {
		return 0, 0, 0, 0, nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if bytes.HasPrefix(raw, []byte{0xEF, 0xBB, 0xBF}) {
		raw = raw[3:]
	}
	cr := csv.NewReader(bytes.NewReader(raw))
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true
	records, err := cr.ReadAll()
	if err != nil {
		return 0, 0, 0, 0, nil, gerror.WrapCode(consts.CodeInvalidParams, err, "CSV 解析失败")
	}
	if len(records) < 2 {
		return 0, 0, 0, 0, nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if len(records)-1 > batchMemberImportMaxRows {
		return 0, 0, 0, 0, nil, gerror.NewCodef(consts.CodeInvalidParams, "单次最多导入 %d 条", batchMemberImportMaxRows)
	}
	colIdx, err := batchMemberImportParseHeader(records[0])
	if err != nil {
		return 0, 0, 0, 0, nil, err
	}

	appendErr := func(line int, msg string) {
		if len(errors) >= batchMemberImportMaxErrLines {
			return
		}
		errors = append(errors, fmt.Sprintf("第%d行：%s", line, msg))
	}

	parsed := make([]batchMemberImportRow, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		lineNo := i + 1
		row := records[i]
		if batchMemberImportRowEmpty(row) {
			continue
		}
		total++
		item := batchMemberImportRow{lineNo: lineNo}
		if idx, ok := colIdx["member_id"]; ok {
			v := batchMemberImportCell(row, idx)
			if v != "" {
				n, e := strconv.ParseInt(v, 10, 64)
				if e != nil || n <= 0 {
					failed++
					appendErr(lineNo, "会员ID 格式无效")
					continue
				}
				item.memberID = n
			}
		}
		if item.memberID == 0 {
			item.username = batchMemberImportCell(row, colIdx["username"])
			if item.username == "" {
				failed++
				appendErr(lineNo, "用户名不能为空")
				continue
			}
		}
		parsed = append(parsed, item)
	}
	if len(parsed) == 0 {
		return total, 0, failed, 0, errors, nil
	}

	usernameRows := make([]batchMemberImportRow, 0, len(parsed))
	idSet := make(map[int64]struct{}, len(parsed))
	usernameSet := make(map[string]struct{}, len(parsed))
	for _, p := range parsed {
		if p.memberID > 0 {
			idSet[p.memberID] = struct{}{}
			continue
		}
		key := strings.TrimSpace(p.username)
		usernameRows = append(usernameRows, p)
		usernameSet[key] = struct{}{}
	}

	usernameToID := map[string]int64{}
	if len(usernameSet) > 0 {
		names := make([]interface{}, 0, len(usernameSet))
		for name := range usernameSet {
			names = append(names, name)
		}
		var users []sysentity.SysMember
		if err := dao.SysMember.Ctx(ctx).
			Fields("id", "username").
			WhereIn("username", names).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&users); err != nil {
			return total, 0, failed, 0, errors, err
		}
		for _, u := range users {
			usernameToID[strings.TrimSpace(u.Username)] = u.Id
		}
	}

	resolvedIDs := make([]int64, 0, len(parsed))
	lineByID := make(map[int64]int, len(parsed))
	for _, p := range parsed {
		id := p.memberID
		if id == 0 {
			id = usernameToID[strings.TrimSpace(p.username)]
			if id <= 0 {
				failed++
				appendErr(p.lineNo, fmt.Sprintf("用户名「%s」不存在", p.username))
				continue
			}
		}
		if prevLine, ok := lineByID[id]; ok {
			failed++
			appendErr(p.lineNo, fmt.Sprintf("与第%d行重复", prevLine))
			continue
		}
		lineByID[id] = p.lineNo
		resolvedIDs = append(resolvedIDs, id)
	}
	if len(resolvedIDs) == 0 {
		return total, 0, failed, 0, errors, nil
	}

	inserted, err = s.ExamBatchMembersAdd(ctx, batchID, examPaperID, resolvedIDs, creator)
	if err != nil {
		return total, 0, failed, 0, errors, err
	}
	success = inserted
	already := len(resolvedIDs) - inserted
	if already > 0 {
		failed += already
		if len(errors) < batchMemberImportMaxErrLines {
			errors = append(errors, fmt.Sprintf("有 %d 名学员已在当前批次中，已跳过", already))
		}
	}
	return total, success, failed, inserted, errors, nil
}

func batchMemberImportParseHeader(header []string) (map[string]int, error) {
	idx := make(map[string]int)
	for i, h := range header {
		key := strings.TrimSpace(h)
		if key == "" {
			continue
		}
		if canon, ok := batchMemberImportHeaderKeys[strings.ToLower(key)]; ok {
			idx[canon] = i
			continue
		}
		if canon, ok := batchMemberImportHeaderKeys[key]; ok {
			idx[canon] = i
		}
	}
	if _, ok := idx["username"]; !ok {
		if _, okID := idx["member_id"]; !okID {
			return nil, gerror.NewCode(consts.CodeInvalidParams)
		}
	}
	return idx, nil
}

func batchMemberImportRowEmpty(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

func batchMemberImportCell(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[col])
}
