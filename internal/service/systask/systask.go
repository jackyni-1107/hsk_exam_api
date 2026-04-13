// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package systask

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysTask interface {
		TaskList(ctx context.Context, page int, size int, name string, code string, handler string, typ int, status *int) ([]sysentity.SysTask, int, error)
		TaskCreate(ctx context.Context, name string, code string, cronExpr string, handler string, params string, alertReceivers string, remark string, creator string, typ int, delaySeconds int, retryTimes int, retryInterval int, concurrency int, alertOnFail int, status int) (int64, error)
		TaskUpdate(ctx context.Context, id int64, name string, code string, cronExpr string, handler string, params string, alertReceivers string, remark string, typ int, delaySeconds int, retryTimes int, retryInterval int, concurrency int, alertOnFail int, status int) error
		TaskDelete(ctx context.Context, id int64) error
		TaskRun(ctx context.Context, id int64) (string, error)
		TaskLogList(ctx context.Context, page int, size int, taskId int64, runId string, status *int) ([]sysentity.SysTaskLog, int, error)
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
