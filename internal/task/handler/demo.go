package handler

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type demoHandler struct{}

func (demoHandler) Execute(ctx context.Context, taskID int64, params string) error {
	g.Log().Infof(ctx, "[DemoHandler] task_id=%d params=%s", taskID, params)
	return nil
}

type examScoreFinalizeHandler struct{}

func (examScoreFinalizeHandler) Execute(ctx context.Context, taskID int64, params string) error {
	g.Log().Infof(ctx, "[ExamScoreFinalizeHandler] stub task_id=%d (implement batch finalize if needed)", taskID)
	return nil
}
