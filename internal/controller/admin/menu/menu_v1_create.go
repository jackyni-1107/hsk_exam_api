package menu

import (
	"context"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/middleware"
	menusvc "exam/internal/service/menu"
)

func (c *ControllerV1) MenuCreate(ctx context.Context, req *v1.MenuCreateReq) (res *v1.MenuCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := menusvc.Menu().MenuCreate(ctx,
		req.Name, req.Permission, req.Path, req.Icon, req.Component, req.ComponentName, creator,
		req.Type, req.Sort, req.ParentId, req.Status,
		boolToInt(req.Visible), boolToInt(req.KeepAlive), boolToInt(req.AlwaysShow),
	)
	if err != nil {
		return nil, err
	}
	return &v1.MenuCreateRes{Id: id}, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
