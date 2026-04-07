package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	v1 "exam/api/admin/user/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) UserCreate(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	var exist sysentity.SysUser
	_ = dao.SystemUser.Ctx(ctx).Where("username", req.Username).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.user_exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	status := req.Status
	if status != consts.StatusNormal && status != consts.StatusDisabled {
		status = consts.StatusNormal
	}
	id, err := dao.SystemUser.Ctx(ctx).InsertAndGetId(sysdo.SysUser{
		Username: req.Username, Password: string(hash), Nickname: req.Nickname,
		Email: req.Email, Mobile: req.Mobile, Status: status,
		Creator: creator, Updater: creator, DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	for _, rid := range req.RoleIds {
		_, _ = dao.SystemUserRole.Ctx(ctx).Insert(sysdo.SysUserRole{
			UserId: id, RoleId: rid, Creator: creator, Updater: creator,
			DeleteFlag: consts.DeleteFlagNotDeleted,
		})
	}
	return &v1.UserCreateRes{Id: id}, nil
}
