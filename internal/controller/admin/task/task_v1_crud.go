package task

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "exam/api/admin/task/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	taskpkg "exam/internal/task"
)

func (c *ControllerV1) TaskCreate(ctx context.Context, req *v1.TaskCreateReq) (res *v1.TaskCreateRes, err error) {
	if err := validateTaskTypeAndSchedule(req.Type, req.CronExpr, req.DelaySeconds); err != nil {
		return nil, err
	}
	var exist sysentity.SysTask
	_ = dao.SysTask.Ctx(ctx).Where("code", req.Code).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeTaskCodeExists, "")
	}
	id, err := dao.SysTask.Ctx(ctx).InsertAndGetId(sysdo.SysTask{
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		CronExpr:       req.CronExpr,
		DelaySeconds:   req.DelaySeconds,
		Handler:        req.Handler,
		Params:         req.Params,
		RetryTimes:     req.RetryTimes,
		RetryInterval:  req.RetryInterval,
		Concurrency:    req.Concurrency,
		AlertOnFail:    req.AlertOnFail,
		AlertReceivers: req.AlertReceivers,
		Status:         req.Status,
		Remark:         req.Remark,
		Creator:        "admin",
		CreateTime:     gtime.Now(),
		Updater:        "admin",
		UpdateTime:     gtime.Now(),
		DeleteFlag:     consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCreateRes{Id: id}, nil
}

func (c *ControllerV1) TaskUpdate(ctx context.Context, req *v1.TaskUpdateReq) (res *v1.TaskUpdateRes, err error) {
	if err := validateTaskTypeAndSchedule(req.Type, req.CronExpr, req.DelaySeconds); err != nil {
		return nil, err
	}
	var cur sysentity.SysTask
	if err := dao.SysTask.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&cur); err != nil {
		return nil, err
	}
	if cur.Id == 0 {
		return nil, gerror.NewCode(consts.CodeTaskNotFound, "")
	}
	var dup sysentity.SysTask
	_ = dao.SysTask.Ctx(ctx).
		Where("code", req.Code).
		WhereNot("id", req.Id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&dup)
	if dup.Id > 0 {
		return nil, gerror.NewCode(consts.CodeTaskCodeExists, "")
	}
	_, err = dao.SysTask.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysTask{
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		CronExpr:       req.CronExpr,
		DelaySeconds:   req.DelaySeconds,
		Handler:        req.Handler,
		Params:         req.Params,
		RetryTimes:     req.RetryTimes,
		RetryInterval:  req.RetryInterval,
		Concurrency:    req.Concurrency,
		AlertOnFail:    req.AlertOnFail,
		AlertReceivers: req.AlertReceivers,
		Status:         req.Status,
		Remark:         req.Remark,
		Updater:        "admin",
		UpdateTime:     gtime.Now(),
	}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.TaskUpdateRes{}, nil
}

func (c *ControllerV1) TaskDelete(ctx context.Context, req *v1.TaskDeleteReq) (res *v1.TaskDeleteRes, err error) {
	_, err = dao.SysTask.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(sysdo.SysTask{DeleteFlag: consts.DeleteFlagDeleted, Updater: "admin", UpdateTime: gtime.Now()}).Update()
	if err != nil {
		return nil, err
	}
	return &v1.TaskDeleteRes{}, nil
}

func (c *ControllerV1) TaskRun(ctx context.Context, req *v1.TaskRunReq) (res *v1.TaskRunRes, err error) {
	var t sysentity.SysTask
	if err := dao.SysTask.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&t); err != nil {
		return nil, err
	}
	if t.Id == 0 {
		return nil, gerror.NewCode(consts.CodeTaskNotFound, "")
	}
	if t.Status != consts.TaskStatusEnabled {
		return nil, gerror.NewCode(consts.CodeTaskDisabled, "")
	}
	runID := gconv.String(gtime.Now().TimestampMilli())
	go func() {
		_ = taskpkg.Execute(context.Background(), taskpkg.ExecRequest{
			TaskID:      t.Id,
			RunID:       runID,
			TriggerType: consts.TriggerTypeManual,
			RetryCount:  0,
		})
	}()
	return &v1.TaskRunRes{RunId: runID}, nil
}

func validateTaskTypeAndSchedule(typeVal int, cronExpr string, delaySec int) error {
	switch typeVal {
	case consts.TaskTypeCron:
		if strings.TrimSpace(cronExpr) == "" {
			return gerror.NewCode(consts.CodeCronExprRequired, "")
		}
	case consts.TaskTypeDelay:
		if delaySec <= 0 {
			return gerror.NewCode(consts.CodeDelaySecondsRequired, "")
		}
	default:
		return gerror.NewCode(consts.CodeInvalidParams, "")
	}
	return nil
}
