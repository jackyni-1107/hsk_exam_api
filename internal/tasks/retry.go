package tasks

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// EnqueueRetry schedules a delayed retry for the same task execution.
func EnqueueRetry(ctx context.Context, taskID int64, runID string, triggerType int, retryCount int, delaySec int) {
	if delaySec <= 0 {
		return
	}
	req := ExecRequest{
		TaskID:      taskID,
		RunID:       runID,
		TriggerType: triggerType,
		RetryCount:  retryCount,
	}
	if err := scheduleAsync(ctx, req, delaySec); err != nil {
		g.Log().Errorf(ctx, "enqueue retry failed task_id=%d run_id=%s: %v", taskID, runID, err)
		scheduleAsyncFallback(ctx, req, delaySec)
	}
}
