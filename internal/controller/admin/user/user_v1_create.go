package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
	"exam/internal/middleware"
	usersvc "exam/internal/service/sysuser"
)

func (c *ControllerV1) UserCreate(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := usersvc.SysUser().UserCreate(ctx, req.Username, req.Password, req.Nickname, req.Email, req.Mobile, creator, req.Status, req.RoleIds)
	if err != nil {
		return nil, err
	}
	return &v1.UserCreateRes{Id: id}, nil
}
