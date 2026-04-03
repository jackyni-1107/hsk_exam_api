package user

import (
	"context"

	"exam/api/admin/user/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) UserRoleAssign(ctx context.Context, req *v1.UserRoleAssignReq) (res *v1.UserRoleAssignRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	_, err = dao.SystemUserRole.Ctx(ctx).Where("user_id", req.Id).Delete()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	for _, rid := range req.RoleIds {
		_, err = dao.SystemUserRole.Ctx(ctx).Insert(dosys.SysUserRole{
			UserId: req.Id, RoleId: rid, Creator: creator, Updater: creator,
			DeleteFlag: consts.DeleteFlagNotDeleted,
		})
		if err != nil {
			return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
		}
	}
	return &v1.UserRoleAssignRes{}, nil
}
