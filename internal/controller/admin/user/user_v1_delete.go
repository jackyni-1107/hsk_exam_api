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

func (c *ControllerV1) UserDelete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	if req.Id == consts.SuperAdminUserId {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.cannot_delete_super_admin")
	}
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SystemUser.Ctx(ctx).Where("id", req.Id).Data(dosys.SysUser{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.UserDeleteRes{}, nil
}
