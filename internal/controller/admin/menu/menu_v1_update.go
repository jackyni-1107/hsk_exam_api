package menu

import (
	"context"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	menusvc "exam/internal/service/sysmenu"
)

func (c *ControllerV1) MenuUpdate(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = menusvc.SysMenu().MenuUpdate(ctx, req.Id, bo.MenuUpdateInput{
		Name:          req.Name,
		Permission:    req.Permission,
		Type:          req.Type,
		Sort:          req.Sort,
		ParentID:      req.ParentId,
		Path:          req.Path,
		Icon:          req.Icon,
		Component:     req.Component,
		ComponentName: req.ComponentName,
		Status:        req.Status,
		Visible:       req.Visible,
		KeepAlive:     req.KeepAlive,
		AlwaysShow:    req.AlwaysShow,
	}, updater)
	if err != nil {
		return nil, err
	}
	return &v1.MenuUpdateRes{}, nil
}
