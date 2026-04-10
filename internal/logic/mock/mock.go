package mock

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	"exam/internal/dao"
	examdao "exam/internal/dao/exam"
	mockentity "exam/internal/model/entity/mock"
	mocksvc "exam/internal/service/mock"
)

func examinationPaperImportFilter(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "imported", "1", "yes", "true":
		return "imported"
	case "not_imported", "not-imported", "0", "no", "false":
		return "not_imported"
	default:
		return ""
	}
}

func loadImportedMockPaperIDSet(ctx context.Context, mockIDs []int64) (map[int64]struct{}, error) {
	out := make(map[int64]struct{})
	if len(mockIDs) == 0 {
		return out, nil
	}
	args := make([]interface{}, len(mockIDs))
	for i, id := range mockIDs {
		args[i] = id
	}
	col := examdao.ExamPaper.Columns().MockExaminationPaperId
	var rows []struct {
		Mid int64 `orm:"mock_examination_paper_id"`
	}
	err := examdao.ExamPaper.Ctx(ctx).
		Fields(col).
		WhereIn(col, args).
		Where(examdao.ExamPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.Mid] = struct{}{}
	}
	return out, nil
}

func (s *sMock) MockLevelsList(ctx context.Context) ([]mockentity.MockLevels, error) {
	var list []mockentity.MockLevels
	err := dao.MockLevels.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted).OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, nil
}

func (s *sMock) ExaminationPaperList(ctx context.Context, levelId int64, importStatus string) ([]*mocksvc.MockExaminationPaperWithImport, error) {
	m := dao.MockExaminationPaper.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted).Where("status", 1)
	if levelId > 0 {
		m = m.Where("level_id", levelId)
	}
	f := examinationPaperImportFilter(importStatus)
	switch f {
	case "imported":
		m = m.Where(`id IN (SELECT `+examdao.ExamPaper.Columns().MockExaminationPaperId+
			` FROM `+examdao.ExamPaper.Table()+` WHERE `+examdao.ExamPaper.Columns().DeleteFlag+`=? )`, consts.DeleteFlagNotDeleted)
	case "not_imported":
		m = m.Where(`id NOT IN (SELECT `+examdao.ExamPaper.Columns().MockExaminationPaperId+
			` FROM `+examdao.ExamPaper.Table()+` WHERE `+examdao.ExamPaper.Columns().DeleteFlag+`=? )`, consts.DeleteFlagNotDeleted)
	}
	var list []mockentity.MockExaminationPaper
	err := m.OrderAsc("seq").OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	mockIDs := make([]int64, len(list))
	for i, e := range list {
		mockIDs[i] = e.Id
	}
	importedSet, err := loadImportedMockPaperIDSet(ctx, mockIDs)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*mocksvc.MockExaminationPaperWithImport, 0, len(list))
	for _, e := range list {
		_, imp := importedSet[e.Id]
		items = append(items, &mocksvc.MockExaminationPaperWithImport{
			MockExaminationPaper: e,
			Imported:             imp,
		})
	}
	return items, nil
}

func (s *sMock) ExaminationPaperDetail(ctx context.Context, id int64) (*mocksvc.MockExaminationPaperWithImport, error) {
	var e mockentity.MockExaminationPaper
	err := dao.MockExaminationPaper.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeResourceNotFound)
	}
	n, err := examdao.ExamPaper.Ctx(ctx).
		Where(examdao.ExamPaper.Columns().MockExaminationPaperId, e.Id).
		Where(examdao.ExamPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &mocksvc.MockExaminationPaperWithImport{
		MockExaminationPaper: e,
		Imported:             n > 0,
	}, nil
}
