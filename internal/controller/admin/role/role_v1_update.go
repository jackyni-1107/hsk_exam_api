package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
	"exam/internal/middleware"
	rolesvc "exam/internal/service/sysrole"
)

func (c *ControllerV1) RoleUpdate(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = rolesvc.SysRole().RoleUpdate(ctx, req.Id, req.Name, req.Code, req.Remark, updater, req.Status, req.Sort, req.Type)
	if err != nil {
		return nil, err
	}
	return &v1.RoleUpdateRes{}, nil
}
