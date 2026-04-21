package cmd

import (
	"context"

	"exam/internal/tasks"
)

func startBackgroundRuntimes(ctx context.Context) {
	go SyncAnswer(ctx)
	tasks.StartScheduler(ctx)
}
