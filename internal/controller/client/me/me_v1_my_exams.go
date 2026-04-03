package me

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/client/me/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	svcExam "exam/internal/service/exam"
	"exam/internal/util"
)

func (c *ControllerV1) MyExams(ctx context.Context, req *v1.MyExamsReq) (res *v1.MyExamsRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	page := req.Page
	size := req.Size
	rows, total, err := svcExam.Exam().MyExamBatches(ctx, d.UserId, page, size)
	if err != nil {
		return nil, err
	}
	list := make([]v1.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {
		ids := r.MockLevelIds
		if ids == nil {
			ids = []int64{}
		}
		list = append(list, v1.MyExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			MockExaminationPaperId: r.MockExaminationPaperId,
			PaperTitle:             r.PaperTitle,
			ExamStartAt:            util.ToRFC3339UTC(r.ExamStartAt),
			ExamEndAt:              util.ToRFC3339UTC(r.ExamEndAt),
			MockLevelId:            r.MockLevelId,
			MockLevelIds:           ids,
			WindowStatus:           r.WindowStatus,
		})
	}
	return &v1.MyExamsRes{List: list, Total: total}, nil
}
