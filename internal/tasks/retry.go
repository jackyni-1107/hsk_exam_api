package tasks

import (
	"context"
	"time"
)

// EnqueueRetry 在间隔后异步重试执行同一任务（简化实现：goroutine + sleep）。
func EnqueueRetry(ctx context.Context, taskID int64, runID string, triggerType int, retryCount int, delaySec int) {
	if delaySec <= 0 {
		return
	}
	go func() {
		time.Sleep(time.Duration(delaySec) * time.Second)
		bg := context.Background()
		_ = Execute(bg, ExecRequest{
			TaskID:      taskID,
			RunID:       runID,
			TriggerType: triggerType,
			RetryCount:  retryCount,
		})
	}()
}
