package exampaper

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	"exam/internal/dao"
	mockdao "exam/internal/dao/mock"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/model/entity/mock"
)

// ByMockID 按 mock_examination_paper.id 查询已导入且未删除的 exam_paper。
func ByMockID(ctx context.Context, mockExaminationPaperID int64) (examentity.ExamPaper, error) {
	var p examentity.ExamPaper
	if mockExaminationPaperID <= 0 {
		return p, gerror.NewCode(consts.CodeInvalidParams)
	}
	err := dao.ExamPaper.Ctx(ctx).
		Where(dao.ExamPaper.Columns().MockExaminationPaperId, mockExaminationPaperID).
		Where(dao.ExamPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Scan(&p)
	if err != nil {
		return p, err
	}
	if p.Id == 0 {
		var m mock.MockExaminationPaper
		_ = mockdao.MockExaminationPaper.Ctx(ctx).
			Where(mockdao.MockExaminationPaper.Columns().Id, mockExaminationPaperID).
			Where(mockdao.MockExaminationPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
			Scan(&m)
		if m.Id == 0 {
			return p, gerror.NewCode(consts.CodeMockExamPaperNotFound)
		}
		return p, gerror.NewCode(consts.CodeExamPaperNotImported)
	}
	return p, nil
}

// ExamPaperByMockLevelID 按 mock_levels.id 选取一张已导入 exam_paper 的卷：同等级多张 mock 卷时优先 id 较大且已完成导入者。
func ExamPaperByMockLevelID(ctx context.Context, mockLevelID int64) (examentity.ExamPaper, error) {
	var p examentity.ExamPaper
	if mockLevelID <= 0 {
		return p, gerror.NewCode(consts.CodeInvalidParams)
	}
	var rows []struct {
		Id int64 `json:"id"`
	}
	err := mockdao.MockExaminationPaper.Ctx(ctx).
		Fields(mockdao.MockExaminationPaper.Columns().Id).
		Where(mockdao.MockExaminationPaper.Columns().LevelId, mockLevelID).
		Where(mockdao.MockExaminationPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		OrderDesc(mockdao.MockExaminationPaper.Columns().Id).
		Scan(&rows)
	if err != nil {
		return p, err
	}
	if len(rows) == 0 {
		return p, gerror.NewCode(consts.CodeMockExamPaperNotFound)
	}
	var lastErr error
	for _, r := range rows {
		ep, err := ByMockID(ctx, r.Id)
		if err == nil && ep.Id != 0 {
			return ep, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return p, lastErr
	}
	return p, gerror.NewCode(consts.CodeExamPaperNotImported)
}

// EnsureMockExaminationPaper 校验 mock 卷存在且未删除（导入前调用）。
func EnsureMockExaminationPaper(ctx context.Context, mockExaminationPaperID int64) error {
	if mockExaminationPaperID <= 0 {
		return gerror.NewCode(consts.CodeMockExaminationPaperIdRequired)
	}
	var m mock.MockExaminationPaper
	_ = mockdao.MockExaminationPaper.Ctx(ctx).
		Where(mockdao.MockExaminationPaper.Columns().Id, mockExaminationPaperID).
		Where(mockdao.MockExaminationPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Scan(&m)
	if m.Id == 0 {
		return gerror.NewCode(consts.CodeMockExamPaperNotFound)
	}
	return nil
}
