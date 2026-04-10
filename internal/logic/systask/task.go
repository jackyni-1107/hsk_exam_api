package systask

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/tasks"
)

func validateTaskTypeAndSchedule(typeVal int, cronExpr string, delaySec int) error {
	switch typeVal {
	case consts.TaskTypeCron:
		if strings.TrimSpace(cronExpr) == "" {
			return gerror.NewCode(consts.CodeCronExprRequired)
		}
	case consts.TaskTypeDelay:
		if delaySec <= 0 {
			return gerror.NewCode(consts.CodeDelaySecondsRequired)
		}
	default:
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	return nil
}

func (s *sSysTask) TaskList(ctx context.Context, page, size int, name, code, handler string, typ int, status *int) ([]sysentity.SysTask, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	m := dao.SysTask.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	if code != "" {
		m = m.WhereLike("code", "%"+code+"%")
	}
	if handler != "" {
		m = m.WhereLike("handler", "%"+handler+"%")
	}
	if typ > 0 {
		m = m.Where("type", typ)
	}
	if status != nil {
		m = m.Where("status", *status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysTask
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysTask) TaskCreate(ctx context.Context, name, code, cronExpr, handler, params, alertReceivers, remark, creator string, typ, delaySeconds, retryTimes, retryInterval, concurrency, alertOnFail, status int) (int64, error) {
	if err := validateTaskTypeAndSchedule(typ, cronExpr, delaySeconds); err != nil {
		return 0, err
	}
	n, err := dao.SysTask.Ctx(ctx).Where("code", code).Where("delete_flag", consts.DeleteFlagNotDeleted).Count()
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if n > 0 {
		return 0, gerror.NewCode(consts.CodeTaskCodeExists)
	}
	now := gtime.Now()
	r, err := dao.SysTask.Ctx(ctx).Insert(sysdo.SysTask{
		Name: name, Code: code, Type: typ, CronExpr: cronExpr, DelaySeconds: delaySeconds,
		Handler: handler, Params: params, RetryTimes: retryTimes, RetryInterval: retryInterval,
		Concurrency: concurrency, AlertOnFail: alertOnFail, AlertReceivers: alertReceivers,
		Status: status, Remark: remark, Creator: creator, CreateTime: now, UpdateTime: now,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	return id, nil
}

func (s *sSysTask) TaskUpdate(ctx context.Context, id int64, name, code, cronExpr, handler, params, alertReceivers, remark string, typ, delaySeconds, retryTimes, retryInterval, concurrency, alertOnFail, status int) error {
	if err := validateTaskTypeAndSchedule(typ, cronExpr, delaySeconds); err != nil {
		return err
	}
	var existing sysentity.SysTask
	err := dao.SysTask.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&existing)
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if existing.Id == 0 {
		return gerror.NewCode(consts.CodeTaskNotFound)
	}
	n, err := dao.SysTask.Ctx(ctx).Where("code", code).Where("delete_flag", consts.DeleteFlagNotDeleted).WhereNot("id", id).Count()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if n > 0 {
		return gerror.NewCode(consts.CodeTaskCodeExists)
	}
	_, err = dao.SysTask.Ctx(ctx).Where("id", id).Data(sysdo.SysTask{
		Name: name, Code: code, Type: typ, CronExpr: cronExpr, DelaySeconds: delaySeconds,
		Handler: handler, Params: params, RetryTimes: retryTimes, RetryInterval: retryInterval,
		Concurrency: concurrency, AlertOnFail: alertOnFail, AlertReceivers: alertReceivers,
		Status: status, Remark: remark, UpdateTime: gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysTask) TaskDelete(ctx context.Context, id int64) error {
	_, err := dao.SysTask.Ctx(ctx).Where("id", id).Data(sysdo.SysTask{
		DeleteFlag: consts.DeleteFlagDeleted, UpdateTime: gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysTask) TaskRun(ctx context.Context, id int64) (string, error) {
	var t sysentity.SysTask
	err := dao.SysTask.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&t)
	if err != nil {
		return "", gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if t.Id == 0 {
		return "", gerror.NewCode(consts.CodeTaskNotFound)
	}
	if t.Status != consts.TaskStatusEnabled {
		return "", gerror.NewCode(consts.CodeTaskDisabled)
	}
	runID := gconv.String(gtime.TimestampNano())
	go func() {
		_ = tasks.Execute(context.Background(), tasks.ExecRequest{
			TaskID: t.Id, RunID: runID, TriggerType: consts.TriggerTypeManual, RetryCount: 0,
		})
	}()
	return runID, nil
}

func (s *sSysTask) TaskLogList(ctx context.Context, page, size int, taskId int64, runId string, status *int) ([]sysentity.SysTaskLog, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	m := dao.SysTaskLog.Ctx(ctx)
	if taskId > 0 {
		m = m.Where("task_id", taskId)
	}
	if runId != "" {
		m = m.Where("run_id", runId)
	}
	if status != nil {
		m = m.Where("status", *status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysTaskLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}
