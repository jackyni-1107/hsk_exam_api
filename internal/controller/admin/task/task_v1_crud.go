package task

import (
	"context"

	v1 "exam/api/admin/task/v1"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	systasksvc "exam/internal/service/systask"
)

func (c *ControllerV1) TaskCreate(ctx context.Context, req *v1.TaskCreateReq) (res *v1.TaskCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil && d.Username != "" {
		creator = d.Username
	}
	id, err := systasksvc.SysTask().TaskCreate(ctx, bo.TaskCreateInput{
		Name:           req.Name,
		Code:           req.Code,
		CronExpr:       req.CronExpr,
		Handler:        req.Handler,
		Params:         req.Params,
		AlertReceivers: req.AlertReceivers,
		Remark:         req.Remark,
		Creator:        creator,
		Type:           req.Type,
		DelaySeconds:   req.DelaySeconds,
		RetryTimes:     req.RetryTimes,
		RetryInterval:  req.RetryInterval,
		Concurrency:    req.Concurrency,
		AlertOnFail:    req.AlertOnFail,
		Status:         req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskCreateRes{Id: id}, nil
}

func (c *ControllerV1) TaskUpdate(ctx context.Context, req *v1.TaskUpdateReq) (res *v1.TaskUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil && d.Username != "" {
		updater = d.Username
	}
	err = systasksvc.SysTask().TaskUpdate(ctx, bo.TaskUpdateInput{
		Id:             req.Id,
		Name:           req.Name,
		Code:           req.Code,
		CronExpr:       req.CronExpr,
		Handler:        req.Handler,
		Params:         req.Params,
		AlertReceivers: req.AlertReceivers,
		Remark:         req.Remark,
		Updater:        updater,
		Type:           req.Type,
		DelaySeconds:   req.DelaySeconds,
		RetryTimes:     req.RetryTimes,
		RetryInterval:  req.RetryInterval,
		Concurrency:    req.Concurrency,
		AlertOnFail:    req.AlertOnFail,
		Status:         req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TaskUpdateRes{}, nil
}

func (c *ControllerV1) TaskDelete(ctx context.Context, req *v1.TaskDeleteReq) (res *v1.TaskDeleteRes, err error) {
	if err = systasksvc.SysTask().TaskDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.TaskDeleteRes{}, nil
}

func (c *ControllerV1) TaskRun(ctx context.Context, req *v1.TaskRunReq) (res *v1.TaskRunRes, err error) {
	runID, err := systasksvc.SysTask().TaskRun(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.TaskRunRes{RunId: runID}, nil
}
