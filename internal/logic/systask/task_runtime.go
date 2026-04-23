package systask

import (
	"context"

	"exam/internal/consts"
	"exam/internal/model/bo"
	"exam/internal/tasks"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysTask) TaskRuntimeStats(ctx context.Context) (*bo.TaskRuntimeStats, error) {
	stats, err := tasks.RuntimeStats(ctx)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &bo.TaskRuntimeStats{
		DelayQueueSize:         stats.DelayQueueSize,
		DelayDueCount:          stats.DelayDueCount,
		DelayScannerActive:     stats.DelayScannerActive,
		DelayScannerTTLMillis:  stats.DelayScannerTTLMillis,
		DelayOldestDueAtMillis: stats.DelayOldestDueAtMillis,
	}, nil
}
