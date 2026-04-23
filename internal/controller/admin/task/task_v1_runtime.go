package task

import (
	"context"
	"time"

	v1 "exam/api/admin/task/v1"
	systasksvc "exam/internal/service/systask"
)

func (c *ControllerV1) TaskRuntimeStats(ctx context.Context, req *v1.TaskRuntimeStatsReq) (res *v1.TaskRuntimeStatsRes, err error) {
	stats, err := systasksvc.SysTask().TaskRuntimeStats(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.TaskRuntimeStatsRes{
		DelayQueueSize:        stats.DelayQueueSize,
		DelayDueCount:         stats.DelayDueCount,
		DelayScannerActive:    stats.DelayScannerActive,
		DelayScannerTTLMillis: stats.DelayScannerTTLMillis,
	}
	if stats.DelayOldestDueAtMillis > 0 {
		res.DelayOldestDueAt = time.UnixMilli(stats.DelayOldestDueAtMillis).UTC().Format(time.RFC3339)
	}
	return res, nil
}
