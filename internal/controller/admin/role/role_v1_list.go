package role

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/admin/role/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"
)

func (c *ControllerV1) RoleList(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	model := dao.SystemRole.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Name != "" {
		model = model.WhereLike("name", "%"+req.Name+"%")
	}
	// status: -1=全部 0=正常 1=停用，仅当明确指定时过滤
	if req.Status == consts.StatusNormal || req.Status == consts.StatusDisabled {
		model = model.Where("status", req.Status)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}

	var roles []sysentity.SysRole
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&roles)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
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
			item.CreateTime = r.CreateTime.Format("Y-m-d H:i:s")
		}
		var roleMenus []sysentity.SysRoleMenu
		_ = dao.SystemRoleMenu.Ctx(ctx).Where("role_id", r.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&roleMenus)
		for _, rm := range roleMenus {
			item.MenuIds = append(item.MenuIds, rm.MenuId)
		}
		list = append(list, item)
	}

	return &v1.RoleListRes{List: list, Total: total}, nil
}
