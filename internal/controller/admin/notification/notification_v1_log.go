package notification

import (
	"context"

	"exam/api/admin/notification/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) LogList(ctx context.Context, req *v1.LogListReq) (res *v1.LogListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SysNotificationLog.Ctx(ctx)
	if req.Channel != "" {
		model = model.Where("channel", req.Channel)
	}
	if req.Recipient != "" {
		model = model.WhereLike("recipient", "%"+req.Recipient+"%")
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []entity.SysNotificationLog
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.LogItem, 0, len(list))
	for _, e := range list {
		item := &v1.LogItem{
			Id: int64(e.Id), TemplateCode: e.TemplateCode, Channel: e.Channel,
			Recipient: e.Recipient, Status: e.Status, ErrorMsg: e.ErrorMsg,
		}
		if e.CreateTime != nil {
			item.CreateTime = e.CreateTime.Format("Y-m-d H:i:s")
		}
		items = append(items, item)
	}
	return &v1.LogListRes{List: items, Total: total}, nil
}
