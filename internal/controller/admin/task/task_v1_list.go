package task

import (
	"context"

	v1 "exam/api/admin/task/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) TaskList(ctx context.Context, req *v1.TaskListReq) (res *v1.TaskListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SysTask.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Name != "" {
		model = model.WhereLike("name", "%"+req.Name+"%")
	}
	if req.Code != "" {
		model = model.WhereLike("code", "%"+req.Code+"%")
	}
	if req.Type > 0 {
		model = model.Where("type", req.Type)
	}
	if req.Status != nil {
		model = model.Where("status", *req.Status)
	}
	if req.Handler != "" {
		model = model.Where("handler", req.Handler)
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysTask
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
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
			item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.TaskListRes{List: items, Total: total}, nil
}
