// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package systask

import (
	"context"
	"exam/internal/model/bo"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysTask interface {
		TaskList(ctx context.Context, page int, size int, name string, code string, handler string, typ int, status *int) ([]sysentity.SysTask, int, error)
		TaskCreate(ctx context.Context, input bo.TaskCreateInput) (int64, error)
		TaskUpdate(ctx context.Context, input bo.TaskUpdateInput) error
		TaskDelete(ctx context.Context, id int64) error
		TaskRun(ctx context.Context, id int64) (string, error)
		TaskLogList(ctx context.Context, page int, size int, taskId int64, runId string, status *int) ([]sysentity.SysTaskLog, int, error)
		TaskRuntimeStats(ctx context.Context) (*bo.TaskRuntimeStats, error)
	}
)

var (
	localSysTask ISysTask
)

func SysTask() ISysTask {
	if localSysTask == nil {
		panic("implement not found for interface ISysTask, forgot register?")
	}
	return localSysTask
}

func RegisterSysTask(i ISysTask) {
	localSysTask = i
}
