package task

import (
	"context"

	v1 "exam/api/admin/task/v1"
	systasksvc "exam/internal/service/systask"
	"exam/internal/utility"
)

func (c *ControllerV1) TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error) {
	list, total, err := systasksvc.SysTask().TaskList(ctx, req.Page, req.Size, req.Name, req.Code, req.Handler, req.Type, req.Status)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.TaskItem, 0, len(list))
	for _, e := range list {
		item := &v1.TaskItem{
			Id: e.Id, Name: e.Name, Code: e.Code, Type: e.Type, CronExpr: e.CronExpr,
			DelaySeconds: e.DelaySeconds, Handler: e.Handler, Params: e.Params,
			RetryTimes: e.RetryTimes, RetryInterval: e.RetryInterval, Concurrency: e.Concurrency,
			AlertOnFail: e.AlertOnFail, AlertReceivers: e.AlertReceivers, Status: e.Status, Remark: e.Remark,
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.TaskListRes{List: items, Total: total}, nil
}
