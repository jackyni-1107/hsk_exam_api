package tasks

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/tasks/handler"
	"exam/internal/utility/notification"
)

// ExecRequest 单次任务执行入参。
type ExecRequest struct {
	TaskID      int64
	RunID       string
	TriggerType int
	RetryCount  int
	Params      string // 可选覆盖库中 params
}

// Execute 执行任务：并发控制、日志、Handler、重试与告警。
func Execute(ctx context.Context, req ExecRequest) error {
	var t sysentity.SysTask
	err := sysdao.SysTask.Ctx(ctx).Where("id", req.TaskID).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&t)
	if err != nil || t.Id == 0 {
		return err
	}
	params := t.Params
	if req.Params != "" {
		params = req.Params
	}

	// 1. 并发控制：concurrency=0 用集群锁；>0 用信号量
	if t.Concurrency <= 0 {
		ok, err := TryClusterExecLock(ctx, req.TaskID)
		if err != nil {
			return err
		}
		if !ok {
			g.Log().Infof(ctx, "[Task] skip, cluster exec lock held task_id=%d run_id=%s (another node or instance is running)", req.TaskID, req.RunID)
			return nil
		}
		defer func() { ReleaseClusterExecLock(ctx, req.TaskID) }()
	} else {
		acquired, err := TryAcquireSem(ctx, req.TaskID, t.Concurrency)
		if err != nil {
			return err
		}
		if !acquired {
			g.Log().Infof(ctx, "[Task] skip, concurrency limit task_id=%d", req.TaskID)
			return nil
		}
		defer func() { ReleaseSem(ctx, req.TaskID) }()
	}

	// 2. 写入运行中日志
	node, _ := os.Hostname()
	logData := sysdo.SysTaskLog{
		TaskId: req.TaskID, RunId: req.RunID, TriggerType: req.TriggerType,
		Status: consts.TaskRunStatusRunning, RetryCount: req.RetryCount, Node: node,
	}
	now := gtime.Now()
	logData.StartTime = now
	_, err = sysdao.SysTaskLog.Ctx(ctx).Insert(logData)
	if err != nil {
		g.Log().Error(ctx, "task log insert failed", err)
		return err
	}

	// 3. 执行 handler，recover panic
	var execErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				execErr = gerror.Newf("panic: %v", r)
			}
		}()
		execErr = handler.Execute(ctx, t.Handler, req.TaskID, params)
	}()

	// 4. 更新结束时间与状态
	endTime := gtime.Now()
	durationMs := int(endTime.TimestampMilli() - now.TimestampMilli())
	status := consts.TaskRunStatusSuccess
	errorMsg := ""
	result := "ok"
	if execErr != nil {
		status = consts.TaskRunStatusFailed
		errorMsg = execErr.Error()
		if len(errorMsg) > 1024 {
			errorMsg = errorMsg[:1024]
		}
		result = ""
	}
	_, _ = sysdao.SysTaskLog.Ctx(ctx).Where("task_id", req.TaskID).Where("run_id", req.RunID).
		Data(sysdo.SysTaskLog{
			Status: status, EndTime: endTime, DurationMs: durationMs,
			ErrorMsg: errorMsg, Result: result,
		}).Update()

	if execErr != nil {
		// 5. 重试或告警
		if req.RetryCount < t.RetryTimes && t.RetryInterval > 0 {
			EnqueueRetry(ctx, req.TaskID, req.RunID, req.TriggerType, req.RetryCount+1, t.RetryInterval)
		} else if t.AlertOnFail == 1 && t.AlertReceivers != "" {
			sendFailAlert(ctx, &t, req.RunID, errorMsg)
		}
		return execErr
	}
	return nil
}

func sendFailAlert(ctx context.Context, t *sysentity.SysTask, runID, errMsg string) {
	receivers := strings.Split(t.AlertReceivers, ",")
	vars := map[string]string{
		"task_name": t.Name,
		"task_code": t.Code,
		"run_id":    runID,
		"error_msg": errMsg,
	}
	for _, r := range receivers {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		content := notification.RenderTemplate("任务 {{task_name}} 执行失败\nRunID: {{run_id}}\n错误: {{error_msg}}", vars)
		if strings.Contains(r, "@") {
			_ = (&notification.EmailSender{}).Send(ctx, r, content)
		} else {
			_ = (&notification.SMSSender{}).Send(ctx, r, content)
		}
	}
}
