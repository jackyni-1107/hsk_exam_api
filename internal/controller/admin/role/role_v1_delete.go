package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
	"exam/internal/middleware"
	rolesvc "exam/internal/service/sysrole"
)

func (c *ControllerV1) RoleDelete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = rolesvc.SysRole().RoleDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.RoleDeleteRes{}, nil
}
