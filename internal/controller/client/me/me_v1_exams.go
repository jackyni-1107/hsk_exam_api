package me

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/client/me/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	batchsvc "exam/internal/service/batch"
	"exam/internal/utility"
)

func (c *ControllerV1) MyExams(ctx context.Context, req *v1.ExamsReq) (res *v1.ExamsRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	rows, err := batchsvc.Batch().MyExamBatches(ctx, d.UserId)
	if err != nil {
		return nil, err
	}
	list := make([]v1.ExamBatchItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, v1.ExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			ExamPaperId:            r.ExamPaperId,
			MockExaminationPaperId: r.MockExaminationPaperId,
			PaperTitle:             r.PaperTitle,
			ExamStartAt:            utility.ToRFC3339UTC(r.ExamStartAt),
			ExamEndAt:              utility.ToRFC3339UTC(r.ExamEndAt),
			AttemptId:              r.AttemptId,
			WindowStatus:           r.WindowStatus,
			BatchKind:              r.BatchKind,
			AllowMultipleAttempts:  r.AllowMultipleAttempts,
		})
	}
	return &v1.ExamsRes{List: list}, nil
}
