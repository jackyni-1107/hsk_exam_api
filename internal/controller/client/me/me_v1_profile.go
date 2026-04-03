package me

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/api/client/me/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	entitysys "exam/internal/model/entity/sys"
)

func gtimeRFC3339(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Time.Format("2006-01-02T15:04:05Z07:00")
}

func (c *ControllerV1) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	var u entitysys.SysMember
	err = dao.SysMember.Ctx(ctx).
		Where("id", d.UserId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&u)
	if err != nil || u.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.user_not_found")
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
		PasswordChangedAt:  gtimeRFC3339(u.PasswordChangedAt),
		LoginIp:            u.LoginIp,
		LoginTime:          gtimeRFC3339(u.LoginTime),
		CreateTime:         gtimeRFC3339(u.CreateTime),
	}, nil
}
