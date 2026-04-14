package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
	"exam/internal/middleware"
	usersvc "exam/internal/service/sysuser"
)

func (c *ControllerV1) UserDelete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = usersvc.SysUser().UserDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.UserDeleteRes{}, nil
}
