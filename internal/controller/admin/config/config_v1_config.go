package config

import (
	"context"

	v1 "exam/api/admin/config/v1"
	"exam/internal/middleware"
	sysconfigsvc "exam/internal/service/sysconfig"
	"exam/internal/utility"
)

func (c *ControllerV1) ConfigList(ctx context.Context, req *v1.ConfigListReq) (res *v1.ConfigListRes, err error) {
	list, total, err := sysconfigsvc.SysConfig().ConfigList(ctx, req.Page, req.Size, req.Group, req.Key)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ConfigItem, 0, len(list))
	for _, e := range list {
		item := &v1.ConfigItem{
			Id: int64(e.Id), ConfigKey: e.ConfigKey, ConfigValue: e.ConfigValue,
			ConfigType: e.ConfigType, GroupName: e.GroupName, Remark: e.Remark,
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.ConfigListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) ConfigCreate(ctx context.Context, req *v1.ConfigCreateReq) (res *v1.ConfigCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := sysconfigsvc.SysConfig().ConfigCreate(ctx, req.ConfigKey, req.ConfigValue, req.ConfigType, req.GroupName, req.Remark, creator)
	if err != nil {
		return nil, err
	}
	return &v1.ConfigCreateRes{Id: id}, nil
}

func (c *ControllerV1) ConfigUpdate(ctx context.Context, req *v1.ConfigUpdateReq) (res *v1.ConfigUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysconfigsvc.SysConfig().ConfigUpdate(ctx, req.Id, req.ConfigValue, req.Remark, updater)
	if err != nil {
		return nil, err
	}
	return &v1.ConfigUpdateRes{}, nil
}

func (c *ControllerV1) ConfigDelete(ctx context.Context, req *v1.ConfigDeleteReq) (res *v1.ConfigDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysconfigsvc.SysConfig().ConfigDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.ConfigDeleteRes{}, nil
}

func (c *ControllerV1) ConfigGet(ctx context.Context, req *v1.ConfigGetReq) (res *v1.ConfigGetRes, err error) {
	val, err := sysconfigsvc.SysConfig().ConfigGet(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &v1.ConfigGetRes{Value: val}, nil
}
