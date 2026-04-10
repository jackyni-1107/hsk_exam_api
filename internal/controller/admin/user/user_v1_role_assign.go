package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
	"exam/internal/middleware"
	usersvc "exam/internal/service/user"
)

func (c *ControllerV1) UserRoleAssign(ctx context.Context, req *v1.UserRoleAssignReq) (res *v1.UserRoleAssignRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	err = usersvc.User().UserRoleAssign(ctx, req.Id, req.RoleIds, creator)
	if err != nil {
		return nil, err
	}
	return &v1.UserRoleAssignRes{}, nil
}
