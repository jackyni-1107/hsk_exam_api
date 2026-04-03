package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"exam/api/admin/user/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := dosys.SysUser{Updater: updater}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
		}
		data.Password = string(hash)
	}
	if req.Nickname != "" {
		data.Nickname = req.Nickname
	}
	if req.Email != "" {
		data.Email = req.Email
	}
	if req.Mobile != "" {
		data.Mobile = req.Mobile
	}
	data.Status = req.Status
	_, err = dao.SystemUser.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.UserUpdateRes{}, nil
}
