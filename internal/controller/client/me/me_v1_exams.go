package me

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/client/me/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	examsvc "exam/internal/service/exam"
	"exam/internal/util"
)

func (c *ControllerV1) MyExams(ctx context.Context, req *v1.ExamsReq) (res *v1.ExamsRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	rows, err := examsvc.Exam().MyExamBatches(ctx, d.UserId)
	if err != nil {
		return nil, err
	}
	list := make([]v1.ExamBatchItem, 0, len(rows))
	for _, r := range rows {
		pids := r.MockExaminationPaperIds
		if pids == nil {
			pids = []int64{}
		}
		list = append(list, v1.ExamBatchItem{
			BatchId:                 r.BatchId,
			Title:                   r.Title,
			MockExaminationPaperId:  r.MockExaminationPaperId,
			MockExaminationPaperIds: pids,
			PaperTitle:              r.PaperTitle,
			ExamStartAt:             util.ToRFC3339UTC(r.ExamStartAt),
			ExamEndAt:               util.ToRFC3339UTC(r.ExamEndAt),
			WindowStatus:            r.WindowStatus,
		})
	}
	return &v1.ExamsRes{List: list}, nil
}
