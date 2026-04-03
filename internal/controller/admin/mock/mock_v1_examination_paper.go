package mock

import (
	"context"
	"strings"

	v1 "exam/api/admin/mock/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	examdao "exam/internal/dao/exam"
	entitymock "exam/internal/model/entity/mock"

	"github.com/gogf/gf/v2/errors/gerror"
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

func (c *ControllerV1) ExaminationPaperList(ctx context.Context, req *v1.ExaminationPaperListReq) (res *v1.ExaminationPaperListRes, err error) {
	m := dao.MockExaminationPaper.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted).Where("status", 1)
	if req.LevelId > 0 {
		m = m.Where("level_id", req.LevelId)
	}
	f := examinationPaperImportFilter(req.ImportStatus)
	colMock := dao.MockExaminationPaper.Columns().Id
	switch f {
	case "imported":
		m = m.Where(`id IN (SELECT `+examdao.ExamPaper.Columns().MockExaminationPaperId+
			` FROM `+examdao.ExamPaper.Table()+` WHERE `+examdao.ExamPaper.Columns().DeleteFlag+`=? )`, consts.DeleteFlagNotDeleted)
	case "not_imported":
		m = m.Where(`id NOT IN (SELECT `+examdao.ExamPaper.Columns().MockExaminationPaperId+
			` FROM `+examdao.ExamPaper.Table()+` WHERE `+examdao.ExamPaper.Columns().DeleteFlag+`=? )`, consts.DeleteFlagNotDeleted)
	}
	var list []entitymock.MockExaminationPaper
	err = m.OrderAsc("seq").OrderAsc(colMock).Scan(&list)
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
	items := make([]*v1.ExaminationPaperItem, 0, len(list))
	for _, e := range list {
		_, imp := importedSet[e.Id]
		items = append(items, &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
			Imported: imp,
		})
	}
	return &v1.ExaminationPaperListRes{List: items}, nil
}

func (c *ControllerV1) ExaminationPaperDetail(ctx context.Context, req *v1.ExaminationPaperDetailReq) (res *v1.ExaminationPaperDetailRes, err error) {
	var e entitymock.MockExaminationPaper
	err = dao.MockExaminationPaper.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.not_found")
	}
	imported := false
	n, err := examdao.ExamPaper.Ctx(ctx).
		Where(examdao.ExamPaper.Columns().MockExaminationPaperId, e.Id).
		Where(examdao.ExamPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	imported = n > 0
	return &v1.ExaminationPaperDetailRes{
		Paper: &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
			Imported: imported,
		},
	}, nil
}
