package menu

import (
	"context"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/middleware"
	menusvc "exam/internal/service/menu"
)

func (c *ControllerV1) MenuDelete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = menusvc.Menu().MenuDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.MenuDeleteRes{}, nil
}
