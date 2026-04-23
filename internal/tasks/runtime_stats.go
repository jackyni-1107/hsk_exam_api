package tasks

import (
	"context"

	"exam/internal/consts"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type RuntimeStatsSnapshot struct {
	DelayQueueSize         int
	DelayDueCount          int
	DelayScannerActive     bool
	DelayScannerTTLMillis  int64
	DelayOldestDueAtMillis int64
}

func RuntimeStats(ctx context.Context) (*RuntimeStatsSnapshot, error) {
	queueSizeVal, err := g.Redis().Do(ctx, "ZCARD", consts.TaskDelayQueueKey)
	if err != nil {
		return nil, err
	}
	nowMs := gtime.Now().TimestampMilli()
	dueCountVal, err := g.Redis().Do(ctx, "ZCOUNT", consts.TaskDelayQueueKey, "-inf", nowMs)
	if err != nil {
		return nil, err
	}
	lockTTLVal, err := g.Redis().Do(ctx, "PTTL", consts.TaskDelayScannerLockKey)
	if err != nil {
		return nil, err
	}
	oldestDueAtMs, err := loadDelayQueueOldestDueAt(ctx)
	if err != nil {
		return nil, err
	}
	lockTTL := gconv.Int64(lockTTLVal.Val())
	if lockTTL < 0 {
		lockTTL = 0
	}
	return &RuntimeStatsSnapshot{
		DelayQueueSize:         gconv.Int(queueSizeVal.Val()),
		DelayDueCount:          gconv.Int(dueCountVal.Val()),
		DelayScannerActive:     lockTTL > 0,
		DelayScannerTTLMillis:  lockTTL,
		DelayOldestDueAtMillis: oldestDueAtMs,
	}, nil
}

func loadDelayQueueOldestDueAt(ctx context.Context) (int64, error) {
	v, err := g.Redis().Do(ctx, "ZRANGE", consts.TaskDelayQueueKey, 0, 0, "WITHSCORES")
	if err != nil {
		return 0, err
	}
	items := gvar.New(v.Val()).Interfaces()
	if len(items) < 2 {
		return 0, nil
	}
	return gconv.Int64(items[1]), nil
}
