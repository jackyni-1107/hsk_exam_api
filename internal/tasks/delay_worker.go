package tasks

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// delayWorkerLoop 延迟任务扫描占位：当前重试由 EnqueueRetry 直接调度，此处仅保活。
func delayWorkerLoop(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = g.Cfg()
		}
	}
}
