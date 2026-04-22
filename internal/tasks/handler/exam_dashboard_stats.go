package handler

import (
	"context"

	attemptsvc "exam/internal/service/attempt"
)

type examDashboardStatsHandler struct{}

func (examDashboardStatsHandler) Execute(ctx context.Context, taskID int64, params string) error {
	return attemptsvc.Attempt().RefreshAttemptDashboardSnapshot(ctx)
}
