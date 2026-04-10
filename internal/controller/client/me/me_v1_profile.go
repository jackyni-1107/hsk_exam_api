package me

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/client/me/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	membersvc "exam/internal/service/member"
	"exam/internal/utility"
)

func (c *ControllerV1) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	u, err := membersvc.Member().MemberProfile(ctx, d.UserId)
	if err != nil {
		return nil, err
	}
	if u == nil {
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
		PasswordChangedAt:  utility.ToRFC3339UTC(u.PasswordChangedAt),
		LoginIp:            u.LoginIp,
		LoginTime:          utility.ToRFC3339UTC(u.LoginTime),
		CreateTime:         utility.ToRFC3339UTC(u.CreateTime),
	}, nil
}
