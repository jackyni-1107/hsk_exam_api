package attempt

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
)

// examBatchFlags 与答题会话相关的批次策略（batchID<=0 时按「历史非批次」保守默认）。
type examBatchFlags struct {
	AllowMultipleAttempts bool
	MaxAttemptsPerMember  int
	SkipScoring           bool
	AutoSubmitOnDeadline  bool
}

func loadExamBatchFlags(ctx context.Context, batchID int64) (examBatchFlags, error) {
	if batchID <= 0 {
		return examBatchFlags{
			AllowMultipleAttempts: false,
			MaxAttemptsPerMember:  0,
			SkipScoring:           false,
			AutoSubmitOnDeadline:  true,
		}, nil
	}
	var b examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).
		Fields(
			dao.ExamBatch.Columns().AllowMultipleAttempts,
			dao.ExamBatch.Columns().MaxAttemptsPerMember,
			dao.ExamBatch.Columns().SkipScoring,
			dao.ExamBatch.Columns().AutoSubmitOnDeadline,
		).
		Where("id", batchID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&b); err != nil {
		return examBatchFlags{}, err
	}
	return examBatchFlags{
		AllowMultipleAttempts: b.AllowMultipleAttempts != 0,
		MaxAttemptsPerMember:  b.MaxAttemptsPerMember,
		SkipScoring:           b.SkipScoring != 0,
		AutoSubmitOnDeadline:  b.AutoSubmitOnDeadline != 0,
	}, nil
}
