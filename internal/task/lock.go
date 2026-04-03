package task

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

const clusterExecLockTTLSeconds = 300

func clusterExecLockKey(taskID int64) string {
	return fmt.Sprintf("task:cluster_exec:%d", taskID)
}

// TryClusterExecLock 集群内任务互斥执行锁（单节点 Redis）。
func TryClusterExecLock(ctx context.Context, taskID int64) (bool, error) {
	key := clusterExecLockKey(taskID)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", clusterExecLockTTLSeconds)
	if err != nil {
		return false, err
	}
	return v.String() == "OK", nil
}

// ReleaseClusterExecLock 释放集群执行锁。
func ReleaseClusterExecLock(ctx context.Context, taskID int64) {
	_, _ = g.Redis().Del(ctx, clusterExecLockKey(taskID))
}
