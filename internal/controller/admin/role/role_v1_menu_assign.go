package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
	"exam/internal/middleware"
	rolesvc "exam/internal/service/sysrole"
)

func (c *ControllerV1) RoleMenuAssign(ctx context.Context, req *v1.RoleMenuAssignReq) (res *v1.RoleMenuAssignRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	err = rolesvc.SysRole().RoleMenuAssign(ctx, req.Id, req.MenuIds, creator)
	if err != nil {
		return nil, err
	}
	return &v1.RoleMenuAssignRes{}, nil
}
