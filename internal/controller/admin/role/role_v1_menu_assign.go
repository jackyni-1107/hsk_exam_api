package role

import (
	"context"

	"exam/api/admin/role/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) RoleMenuAssign(ctx context.Context, req *v1.RoleMenuAssignReq) (res *v1.RoleMenuAssignRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	_, err = dao.SystemRoleMenu.Ctx(ctx).Where("role_id", req.Id).Delete()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	for _, mid := range req.MenuIds {
		_, err = dao.SystemRoleMenu.Ctx(ctx).Insert(dosys.SysRoleMenu{
			RoleId: req.Id, MenuId: mid, Creator: creator, Updater: creator,
			DeleteFlag: consts.DeleteFlagNotDeleted,
		})
		if err != nil {
			return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
		}
	}
	return &v1.RoleMenuAssignRes{}, nil
}
