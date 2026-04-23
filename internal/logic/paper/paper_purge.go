package paper

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
)

// PaperPurgePhysical 从数据库永久删除 exam_paper 及其题目树（option/question/block/section）。
// 调用方须已通过身份校验；删除前需校验 confirmText == "DELETE:{exam_paper_id}"。
// 若试卷仍被考试批次引用或存在未逻辑删除的答题会话，则拒绝删除。
func (s *sPaper) PaperPurgePhysical(ctx context.Context, examPaperId int64, confirmText string) error {
	if examPaperId <= 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	expected := fmt.Sprintf("DELETE:%d", examPaperId)
	if strings.TrimSpace(confirmText) != expected {
		return gerror.NewCode(consts.CodeExamPaperPurgeConfirmInvalid)
	}

	var paper examentity.ExamPaper
	if err := dao.ExamPaper.Ctx(ctx).
		Where(dao.ExamPaper.Columns().Id, examPaperId).
		Where(dao.ExamPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Scan(&paper); err != nil {
		return err
	}
	if paper.Id == 0 {
		return gerror.NewCode(consts.CodeExamPaperNotFound)
	}

	nBatchPaper, err := dao.ExamBatchPaper.Ctx(ctx).
		Where(dao.ExamBatchPaper.Columns().ExamPaperId, examPaperId).
		Count()
	if err != nil {
		return err
	}
	nBatchMember, err := dao.ExamBatchMember.Ctx(ctx).
		Where(dao.ExamBatchMember.Columns().ExamPaperId, examPaperId).
		Count()
	if err != nil {
		return err
	}
	if nBatchPaper > 0 || nBatchMember > 0 {
		return gerror.NewCode(consts.CodeExamPaperPurgeHasBatchBinding)
	}

	nAttempt, err := dao.ExamAttempt.Ctx(ctx).
		Where(dao.ExamAttempt.Columns().ExamPaperId, examPaperId).
		Where(dao.ExamAttempt.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return err
	}
	if nAttempt > 0 {
		return gerror.NewCode(consts.CodeExamPaperPurgeHasAttempts)
	}

	mockID := paper.MockExaminationPaperId
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return deletePaperTreeTx(ctx, tx, examPaperId)
	})
	if err != nil {
		return err
	}

	auditutil.RecordEntityDiff(ctx, dao.ExamPaper.Table(), examPaperId, &paper, nil)

	invalidatePaperCaches(ctx, examPaperId, mockID)
	return nil
}
