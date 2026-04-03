package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
)

type INotification interface {
	LogList(ctx context.Context, req *v1.LogListReq) (res *v1.LogListRes, err error)
	ChannelConfigList(ctx context.Context, req *v1.ChannelConfigListReq) (res *v1.ChannelConfigListRes, err error)
	ChannelConfigCreate(ctx context.Context, req *v1.ChannelConfigCreateReq) (res *v1.ChannelConfigCreateRes, err error)
	ChannelConfigUpdate(ctx context.Context, req *v1.ChannelConfigUpdateReq) (res *v1.ChannelConfigUpdateRes, err error)
	ChannelConfigDelete(ctx context.Context, req *v1.ChannelConfigDeleteReq) (res *v1.ChannelConfigDeleteRes, err error)
	ChannelConfigSetActive(ctx context.Context, req *v1.ChannelConfigSetActiveReq) (res *v1.ChannelConfigSetActiveRes, err error)
	TemplateList(ctx context.Context, req *v1.TemplateListReq) (res *v1.TemplateListRes, err error)
	TemplateCreate(ctx context.Context, req *v1.TemplateCreateReq) (res *v1.TemplateCreateRes, err error)
	TemplateUpdate(ctx context.Context, req *v1.TemplateUpdateReq) (res *v1.TemplateUpdateRes, err error)
	TemplateDelete(ctx context.Context, req *v1.TemplateDeleteReq) (res *v1.TemplateDeleteRes, err error)
	Send(ctx context.Context, req *v1.NotificationSendReq) (res *v1.NotificationSendRes, err error)
}
