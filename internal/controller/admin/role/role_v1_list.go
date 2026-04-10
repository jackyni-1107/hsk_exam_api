package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
	rolesvc "exam/internal/service/role"
	"exam/internal/utility"
)

func (c *ControllerV1) RoleList(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	roles, total, err := rolesvc.Role().RoleList(ctx, req.Page, req.Size, req.Name, req.Status)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.RoleItem, 0, len(roles))
	for _, r := range roles {
		item := &v1.RoleItem{
			Id:     r.Id,
			Name:   r.Name,
			Code:   r.Code,
			Status: r.Status,
			Sort:   r.Sort,
			Type:   r.Type,
			Remark: r.Remark,
		}
		if r.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(r.CreateTime)
		}
		menuIds, _ := rolesvc.Role().RoleMenuIds(ctx, r.Id)
		item.MenuIds = menuIds
		list = append(list, item)
	}

	return &v1.RoleListRes{List: list, Total: total}, nil
}
