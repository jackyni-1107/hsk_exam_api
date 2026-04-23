package tasks

import (
	"context"
	"strings"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/tasks/handler"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func StartRuntime(ctx context.Context) {
	logClusterTimeZoneHint(ctx)
	_ = triggerHandlerInit()
	ReloadSchedules(ctx)
	go refreshCronLoop(ctx)
	go delayQueueLoop(ctx)
}

func StartScheduler(ctx context.Context) {
	StartRuntime(ctx)
}

func ReloadSchedules(ctx context.Context) {
	loadCronTasks(ctx)
}

func ValidateHandler(handlerName string) error {
	handlerName = strings.TrimSpace(handlerName)
	if handlerName == "" {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	if handler.Get(handlerName) == nil {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	return nil
}

func RunNow(ctx context.Context, taskID int64) (string, error) {
	var task sysentity.SysTask
	err := sysdao.SysTask.Ctx(ctx).
		Where("id", taskID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&task)
	if err != nil {
		return "", err
	}
	if task.Id == 0 {
		return "", gerror.NewCode(consts.CodeTaskNotFound)
	}
	if task.Status != consts.TaskStatusEnabled {
		return "", gerror.NewCode(consts.CodeTaskDisabled)
	}

	runID := newRunID()
	if err := dispatchTask(ctx, &task, runID); err != nil {
		return "", err
	}
	return runID, nil
}

func executeAsync(req ExecRequest) {
	go func() {
		_ = Execute(context.Background(), req)
	}()
}

func scheduleAsync(ctx context.Context, req ExecRequest, delaySec int) error {
	if delaySec <= 0 {
		executeAsync(req)
		return nil
	}
	return enqueueDelayedExec(ctx, req, delaySec)
}

func dispatchTask(ctx context.Context, task *sysentity.SysTask, runID string) error {
	req, delaySec := buildDispatchRequest(task, runID)
	if delaySec > 0 {
		return scheduleAsync(ctx, req, delaySec)
	}
	executeAsync(req)
	return nil
}

func newRunID() string {
	return gconv.String(gtime.TimestampNano())
}

func buildDispatchRequest(task *sysentity.SysTask, runID string) (ExecRequest, int) {
	req := ExecRequest{
		TaskID:      task.Id,
		RunID:       runID,
		RetryCount:  0,
		TriggerType: consts.TriggerTypeManual,
	}
	if task.Type == consts.TaskTypeDelay && task.DelaySeconds > 0 {
		req.TriggerType = consts.TriggerTypeDelay
		return req, task.DelaySeconds
	}
	return req, 0
}
