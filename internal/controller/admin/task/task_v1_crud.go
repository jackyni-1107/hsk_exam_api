package task

import (
	"context"

	v1 "exam/api/admin/task/v1"
	systasksvc "exam/internal/service/systask"
)

func (c *ControllerV1) TaskCreate(ctx context.Context, req *v1.TaskCreateReq) (res *v1.TaskCreateRes, err error) {
	id, err := systasksvc.SysTask().TaskCreate(ctx,
		req.Name, req.Code, req.CronExpr, req.Handler, req.Params, req.AlertReceivers, req.Remark,
		"admin", req.Type, req.DelaySeconds, req.RetryTimes, req.RetryInterval,
		req.Concurrency, req.AlertOnFail, req.Status,
	)
	if err != nil {
		return nil, err
	}
	return &v1.TaskCreateRes{Id: id}, nil
}

func (c *ControllerV1) TaskUpdate(ctx context.Context, req *v1.TaskUpdateReq) (res *v1.TaskUpdateRes, err error) {
	err = systasksvc.SysTask().TaskUpdate(ctx,
		req.Id, req.Name, req.Code, req.CronExpr, req.Handler, req.Params, req.AlertReceivers, req.Remark,
		req.Type, req.DelaySeconds, req.RetryTimes, req.RetryInterval,
		req.Concurrency, req.AlertOnFail, req.Status,
	)
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
