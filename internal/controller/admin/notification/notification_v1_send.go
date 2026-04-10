package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
	notisvc "exam/internal/service/sysnotification"
)

func (c *ControllerV1) Send(ctx context.Context, req *v1.NotificationSendReq) (res *v1.NotificationSendRes, err error) {
	ok, err := notisvc.Sysnotification().Send(ctx, req.TemplateCode, req.Channel, req.Recipient, req.Variables)
	if err != nil {
		return nil, err
	}
	return &v1.NotificationSendRes{Ok: ok}, nil
}
