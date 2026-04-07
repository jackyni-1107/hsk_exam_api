package task

import (
	"context"

	v1 "exam/api/admin/task/v1"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) TaskLogList(ctx context.Context, req *v1.TaskLogListReq) (res *v1.TaskLogListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SysTaskLog.Ctx(ctx)
	if req.TaskId > 0 {
		model = model.Where("task_id", req.TaskId)
	}
	if req.RunId != "" {
		model = model.Where("run_id", req.RunId)
	}
	// 未传 status 时不能按 int 零值过滤，否则默认变成 WHERE status=0，成功/失败记录被误过滤
	if req.Status != nil {
		model = model.Where("status", *req.Status)
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count task log")
	}
	var list []sysentity.SysTaskLog
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.Wrap(err, "list task log")
	}
	items := make([]*v1.TaskLogItem, 0, len(list))
	for _, e := range list {
		item := &v1.TaskLogItem{
			Id: e.Id, TaskId: e.TaskId, RunId: e.RunId, TriggerType: e.TriggerType,
			Status: e.Status, DurationMs: e.DurationMs, RetryCount: e.RetryCount,
			ErrorMsg: e.ErrorMsg, Result: e.Result, Node: e.Node,
		}
		if e.StartTime != nil {
			item.StartTime = util.ToRFC3339UTC(e.StartTime)
		}
		if e.EndTime != nil {
			item.EndTime = util.ToRFC3339UTC(e.EndTime)
		}
		if e.CreateTime != nil {
			item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.TaskLogListRes{List: items, Total: total}, nil
}
