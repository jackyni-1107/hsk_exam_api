package tasks

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func logClusterTimeZoneHint(ctx context.Context) {
	g.Log().Info(ctx, "[Task] scheduler started (cluster lock: Redis SET NX)")
}
