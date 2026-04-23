package task

import (
	"context"

	"exam/api/admin/task/v1"
)

type ITask interface {
	TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error)
	TaskRuntimeStats(ctx context.Context, req *v1.TaskRuntimeStatsReq) (res *v1.TaskRuntimeStatsRes, err error)
	TaskCreate(ctx context.Context, req *v1.TaskCreateReq) (res *v1.TaskCreateRes, err error)
	TaskUpdate(ctx context.Context, req *v1.TaskUpdateReq) (res *v1.TaskUpdateRes, err error)
	TaskDelete(ctx context.Context, req *v1.TaskDeleteReq) (res *v1.TaskDeleteRes, err error)
	TaskRun(ctx context.Context, req *v1.TaskRunReq) (res *v1.TaskRunRes, err error)
	TaskLogList(ctx context.Context, req *v1.TaskLogListReq) (res *v1.TaskLogListRes, err error)
}
