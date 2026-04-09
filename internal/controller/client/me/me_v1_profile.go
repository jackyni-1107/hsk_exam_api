package me

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/client/me/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"
)

func (c *ControllerV1) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	var u sysentity.SysMember
	err = dao.SysMember.Ctx(ctx).
		Where("id", d.UserId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&u)
	if err != nil || u.Id == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	return &v1.ProfileRes{
		Id:                 u.Id,
		Username:           u.Username,
		Nickname:           u.Nickname,
		Avatar:             u.Avatar,
		Email:              u.Email,
		Mobile:             u.Mobile,
		Status:             u.Status,
		MustChangePassword: u.MustChangePassword,
		PasswordChangedAt:  util.ToRFC3339UTC(u.PasswordChangedAt),
		LoginIp:            u.LoginIp,
		LoginTime:          util.ToRFC3339UTC(u.LoginTime),
		CreateTime:         util.ToRFC3339UTC(u.CreateTime),
	}, nil
}
