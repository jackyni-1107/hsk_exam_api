package menu

import (
	"context"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/middleware"
	menusvc "exam/internal/service/sysmenu"
)

func (c *ControllerV1) MenuUpdate(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = menusvc.SysMenu().MenuUpdate(ctx,
		req.Id, req.Name, req.Permission, req.Path, req.Icon, req.Component, req.ComponentName, updater,
		req.Type, req.Sort, req.ParentId, req.Status,
		boolToInt(req.Visible), boolToInt(req.KeepAlive), boolToInt(req.AlwaysShow),
	)
	if err != nil {
		return nil, err
	}
	return &v1.MenuUpdateRes{}, nil
}
