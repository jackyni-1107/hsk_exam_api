package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
	"exam/internal/middleware"
	notisvc "exam/internal/service/SysNotification"
	"exam/internal/utility"
)

func (c *ControllerV1) ChannelConfigList(ctx context.Context, req *v1.ChannelConfigListReq) (res *v1.ChannelConfigListRes, err error) {
	list, err := notisvc.SysNotification().ChannelConfigList(ctx, req.Channel)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ChannelConfigItem, 0, len(list))
	for _, e := range list {
		item := &v1.ChannelConfigItem{
			Id:         int64(e.Id),
			Channel:    e.Channel,
			Provider:   e.Provider,
			Name:       e.Name,
			IsActive:   e.IsActive,
			ConfigJson: e.ConfigJson,
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.ChannelConfigListRes{List: items}, nil
}

func (c *ControllerV1) ChannelConfigCreate(ctx context.Context, req *v1.ChannelConfigCreateReq) (res *v1.ChannelConfigCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := notisvc.SysNotification().ChannelConfigCreate(ctx, req.Channel, req.Provider, req.Name, req.ConfigJson, creator)
	if err != nil {
		return nil, err
	}
	return &v1.ChannelConfigCreateRes{Id: id}, nil
}

func (c *ControllerV1) ChannelConfigUpdate(ctx context.Context, req *v1.ChannelConfigUpdateReq) (res *v1.ChannelConfigUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = notisvc.SysNotification().ChannelConfigUpdate(ctx, req.Id, req.Name, req.ConfigJson, updater)
	if err != nil {
		return nil, err
	}
	return &v1.ChannelConfigUpdateRes{}, nil
}

func (c *ControllerV1) ChannelConfigDelete(ctx context.Context, req *v1.ChannelConfigDeleteReq) (res *v1.ChannelConfigDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = notisvc.SysNotification().ChannelConfigDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.ChannelConfigDeleteRes{}, nil
}

func (c *ControllerV1) ChannelConfigSetActive(ctx context.Context, req *v1.ChannelConfigSetActiveReq) (res *v1.ChannelConfigSetActiveRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = notisvc.SysNotification().ChannelConfigSetActive(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.ChannelConfigSetActiveRes{}, nil
}
