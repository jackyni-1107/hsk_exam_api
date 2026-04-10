package tasks

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
)

func clusterExecLockKey(taskID int64) string {
	return fmt.Sprintf(consts.TaskClusterExecLockKeyFmt, taskID)
}

// TryClusterExecLock 集群内任务互斥执行锁（单节点 Redis）。
func TryClusterExecLock(ctx context.Context, taskID int64) (bool, error) {
	key := clusterExecLockKey(taskID)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", consts.TaskClusterExecLockTTLSeconds)
	if err != nil {
		return false, err
	}
	return v.String() == "OK", nil
}

// ReleaseClusterExecLock 释放集群执行锁。
func ReleaseClusterExecLock(ctx context.Context, taskID int64) {
	_, _ = g.Redis().Del(ctx, clusterExecLockKey(taskID))
}
