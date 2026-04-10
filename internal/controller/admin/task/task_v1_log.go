package task

import (
	"context"

	v1 "exam/api/admin/task/v1"
	systasksvc "exam/internal/service/systask"
	"exam/internal/utility"
)

func (c *ControllerV1) TaskLogList(ctx context.Context, req *v1.TaskLogListReq) (res *v1.TaskLogListRes, err error) {
	list, total, err := systasksvc.SysTask().TaskLogList(ctx, req.Page, req.Size, req.TaskId, req.RunId, req.Status)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.TaskLogItem, 0, len(list))
	for _, e := range list {
		item := &v1.TaskLogItem{
			Id: e.Id, TaskId: e.TaskId, RunId: e.RunId, TriggerType: e.TriggerType,
			Status: e.Status, DurationMs: e.DurationMs, RetryCount: e.RetryCount,
			ErrorMsg: e.ErrorMsg, Result: e.Result, Node: e.Node,
		}
		if e.StartTime != nil {
			item.StartTime = utility.ToRFC3339UTC(e.StartTime)
		}
		if e.EndTime != nil {
			item.EndTime = utility.ToRFC3339UTC(e.EndTime)
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.TaskLogListRes{List: items, Total: total}, nil
}
