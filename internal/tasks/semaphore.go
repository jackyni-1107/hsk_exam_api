package tasks

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
)

func semRedisKey(taskID int64) string {
	return fmt.Sprintf(consts.TaskSemKeyFmt, taskID)
}

// TryAcquireSem 基于 Redis 计数的简易并发度限制（limit 为最大同时执行数）。
func TryAcquireSem(ctx context.Context, taskID int64, limit int) (bool, error) {
	if limit <= 0 {
		return true, nil
	}
	key := semRedisKey(taskID)
	n, err := g.Redis().Incr(ctx, key)
	if err != nil {
		return false, err
	}
	if n == 1 {
		_, _ = g.Redis().Expire(ctx, key, consts.TaskSemExpireSeconds)
	}
	if int(n) > limit {
		_, _ = g.Redis().Decr(ctx, key)
		return false, nil
	}
	return true, nil
}

// ReleaseSem 释放一次并发计数。
func ReleaseSem(ctx context.Context, taskID int64) {
	_, _ = g.Redis().Decr(ctx, semRedisKey(taskID))
}
