package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
	"exam/internal/middleware"
	usersvc "exam/internal/service/user"
)

func (c *ControllerV1) UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = usersvc.User().UserUpdate(ctx, req.Id, req.Password, req.Nickname, req.Email, req.Mobile, updater, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.UserUpdateRes{}, nil
}
