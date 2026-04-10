package systask

import (
	"context"

	sysentity "exam/internal/model/entity/sys"
)

type ISysTask interface {
	TaskList(ctx context.Context, page, size int, name, code, handler string, typ int, status *int) ([]sysentity.SysTask, int, error)
	TaskCreate(ctx context.Context, name, code, cronExpr, handler, params, alertReceivers, remark, creator string, typ, delaySeconds, retryTimes, retryInterval, concurrency, alertOnFail, status int) (int64, error)
	TaskUpdate(ctx context.Context, id int64, name, code, cronExpr, handler, params, alertReceivers, remark string, typ, delaySeconds, retryTimes, retryInterval, concurrency, alertOnFail, status int) error
	TaskDelete(ctx context.Context, id int64) error
	TaskRun(ctx context.Context, id int64) (string, error)
	TaskLogList(ctx context.Context, page, size int, taskId int64, runId string, status *int) ([]sysentity.SysTaskLog, int, error)
}

var localSysTask ISysTask

func SysTask() ISysTask {
	if localSysTask == nil {
		panic("implement not found for interface ISysTask, forgot register?")
	}
	return localSysTask
}

func RegisterSysTask(i ISysTask) {
	localSysTask = i
}
