package systask

import (
	"context"
	"strings"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/tasks"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type taskMutationInput struct {
	Name           string
	Code           string
	CronExpr       string
	Handler        string
	Params         string
	AlertReceivers string
	Remark         string
	Type           int
	DelaySeconds   int
	RetryTimes     int
	RetryInterval  int
	Concurrency    int
	AlertOnFail    int
	Status         int
}

const defaultTaskOperator = "system"

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
	if name = strings.TrimSpace(name); name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	if code = strings.TrimSpace(code); code != "" {
		m = m.WhereLike("code", "%"+code+"%")
	}
	if handler = strings.TrimSpace(handler); handler != "" {
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
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysTask) TaskCreate(ctx context.Context, input bo.TaskCreateInput) (int64, error) {
	input.Handler = strings.TrimSpace(input.Handler)
	mutation := taskMutationInput{
		Name:           input.Name,
		Code:           input.Code,
		CronExpr:       input.CronExpr,
		Handler:        input.Handler,
		Params:         input.Params,
		AlertReceivers: input.AlertReceivers,
		Remark:         input.Remark,
		Type:           input.Type,
		DelaySeconds:   input.DelaySeconds,
		RetryTimes:     input.RetryTimes,
		RetryInterval:  input.RetryInterval,
		Concurrency:    input.Concurrency,
		AlertOnFail:    input.AlertOnFail,
		Status:         input.Status,
	}
	data, err := buildTaskMutationData(mutation)
	if err != nil {
		return 0, err
	}
	if data.Name == "" || data.Code == "" {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	if err := tasks.ValidateHandler(input.Handler); err != nil {
		return 0, err
	}

	n, err := dao.SysTask.Ctx(ctx).
		Where("code", data.Code).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if n > 0 {
		return 0, gerror.NewCode(consts.CodeTaskCodeExists)
	}

	now := gtime.Now()
	data.Creator = normalizeTaskOperator(input.Creator)
	data.Updater = data.Creator
	data.CreateTime = now
	data.UpdateTime = now
	data.DeleteFlag = consts.DeleteFlagNotDeleted

	r, err := dao.SysTask.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	var after sysentity.SysTask
	if err := dao.SysTask.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SysTask.Table(), id, nil, &after)
	}
	tasks.ReloadSchedules(ctx)
	return id, nil
}

func (s *sSysTask) TaskUpdate(ctx context.Context, input bo.TaskUpdateInput) error {
	input.Handler = strings.TrimSpace(input.Handler)
	mutation := taskMutationInput{
		Name:           input.Name,
		Code:           input.Code,
		CronExpr:       input.CronExpr,
		Handler:        input.Handler,
		Params:         input.Params,
		AlertReceivers: input.AlertReceivers,
		Remark:         input.Remark,
		Type:           input.Type,
		DelaySeconds:   input.DelaySeconds,
		RetryTimes:     input.RetryTimes,
		RetryInterval:  input.RetryInterval,
		Concurrency:    input.Concurrency,
		AlertOnFail:    input.AlertOnFail,
		Status:         input.Status,
	}
	data, err := buildTaskMutationData(mutation)
	if err != nil {
		return err
	}
	if data.Name == "" || data.Code == "" {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	if err := tasks.ValidateHandler(input.Handler); err != nil {
		return err
	}

	existing, err := loadTaskByID(ctx, input.Id)
	if err != nil {
		return err
	}
	n, err := dao.SysTask.Ctx(ctx).
		Where("code", data.Code).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		WhereNot("id", input.Id).
		Count()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if n > 0 {
		return gerror.NewCode(consts.CodeTaskCodeExists)
	}

	data.Updater = normalizeTaskOperator(input.Updater)
	data.UpdateTime = gtime.Now()
	if _, err := dao.SysTask.Ctx(ctx).Where("id", input.Id).Data(data).Update(); err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var after sysentity.SysTask
	if err := dao.SysTask.Ctx(ctx).Where("id", input.Id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysTask.Table(), input.Id, existing, &after)
	}
	tasks.ReloadSchedules(ctx)
	return nil
}

func (s *sSysTask) TaskDelete(ctx context.Context, id int64) error {
	before, err := loadTaskByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = dao.SysTask.Ctx(ctx).Where("id", id).Data(sysdo.SysTask{
		DeleteFlag: consts.DeleteFlagDeleted,
		UpdateTime: gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var after sysentity.SysTask
	if err := dao.SysTask.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysTask.Table(), id, before, &after)
	}
	tasks.ReloadSchedules(ctx)
	return nil
}

func (s *sSysTask) TaskRun(ctx context.Context, id int64) (string, error) {
	runID, err := tasks.RunNow(ctx, id)
	if err != nil {
		return "", err
	}
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
	if runId = strings.TrimSpace(runId); runId != "" {
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
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func loadTaskByID(ctx context.Context, taskID int64) (*sysentity.SysTask, error) {
	var task sysentity.SysTask
	if err := dao.SysTask.Ctx(ctx).
		Where("id", taskID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&task); err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if task.Id == 0 {
		return nil, gerror.NewCode(consts.CodeTaskNotFound)
	}
	return &task, nil
}

func buildTaskMutationData(input taskMutationInput) (sysdo.SysTask, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Code = strings.TrimSpace(input.Code)
	input.CronExpr = strings.TrimSpace(input.CronExpr)
	input.Handler = strings.TrimSpace(input.Handler)
	input.Params = strings.TrimSpace(input.Params)
	input.AlertReceivers = strings.TrimSpace(input.AlertReceivers)
	input.Remark = strings.TrimSpace(input.Remark)
	input.Status = normalizeTaskStatus(input.Status)
	input.AlertOnFail = normalizeBinaryFlag(input.AlertOnFail)
	input.RetryTimes = normalizeNonNegative(input.RetryTimes)
	input.RetryInterval = normalizeNonNegative(input.RetryInterval)
	input.Concurrency = normalizeNonNegative(input.Concurrency)
	if err := validateTaskTypeAndSchedule(input.Type, input.CronExpr, input.DelaySeconds); err != nil {
		return sysdo.SysTask{}, err
	}

	data := sysdo.SysTask{
		Name:           input.Name,
		Code:           input.Code,
		Type:           input.Type,
		CronExpr:       input.CronExpr,
		DelaySeconds:   input.DelaySeconds,
		Handler:        input.Handler,
		Params:         input.Params,
		RetryTimes:     input.RetryTimes,
		RetryInterval:  input.RetryInterval,
		Concurrency:    input.Concurrency,
		AlertOnFail:    input.AlertOnFail,
		AlertReceivers: input.AlertReceivers,
		Status:         input.Status,
		Remark:         input.Remark,
	}
	if input.Type == consts.TaskTypeCron {
		data.DelaySeconds = 0
	}
	if input.Type == consts.TaskTypeDelay {
		data.CronExpr = ""
	}
	return data, nil
}

func normalizeTaskStatus(status int) int {
	if status == consts.TaskStatusEnabled || status == consts.TaskStatusDisabled {
		return status
	}
	return consts.TaskStatusEnabled
}

func normalizeBinaryFlag(flag int) int {
	if flag == 1 {
		return 1
	}
	return 0
}

func normalizeNonNegative(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

func normalizeTaskOperator(operator string) string {
	operator = strings.TrimSpace(operator)
	if operator == "" {
		return defaultTaskOperator
	}
	return operator
}
