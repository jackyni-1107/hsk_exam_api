package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
	"exam/internal/middleware"
	rolesvc "exam/internal/service/role"
)

func (c *ControllerV1) RoleCreate(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := rolesvc.Role().RoleCreate(ctx, req.Name, req.Code, req.Remark, creator, req.Status, req.Sort, req.Type)
	if err != nil {
		return nil, err
	}
	return &v1.RoleCreateRes{Id: id}, nil
}
