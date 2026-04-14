package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
	notisvc "exam/internal/service/sysnotification"
	"exam/internal/utility"
)

func (c *ControllerV1) LogList(ctx context.Context, req *v1.LogListReq) (res *v1.LogListRes, err error) {
	list, total, err := notisvc.SysNotification().LogList(ctx, req.Page, req.Size, req.Channel, req.Recipient)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.LogItem, 0, len(list))
	for _, e := range list {
		item := &v1.LogItem{
			Id: int64(e.Id), TemplateCode: e.TemplateCode, Channel: e.Channel,
			Recipient: e.Recipient, Status: e.Status, ErrorMsg: e.ErrorMsg,
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.LogListRes{List: items, Total: total}, nil
}
