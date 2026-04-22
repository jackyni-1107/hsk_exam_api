package auth

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "exam/api/admin/auth/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	secsvc "exam/internal/service/security"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	ip, userAgent := "", ""
	if r != nil {
		ip = r.GetClientIp()
		userAgent = r.Header.Get("User-Agent")
	}
	loginRes, err := secsvc.Security().Login(ctx, bo.LoginInput{
		UserType:          consts.UserTypeAdmin,
		Username:          req.Username,
		EncryptedPassword: req.Password,
		CaptchaId:         req.CaptchaId,
		CaptchaAnswer:     req.CaptchaAnswer,
		IP:                ip,
		UserAgent:         userAgent,
		TraceId:           middleware.GetTraceId(ctx),
	})
	if err != nil {
		return nil, err
	}

	var perms []string
	if p, perr := middleware.GetUserPermissions(ctx, loginRes.UserInfo.Id); perr == nil {
		perms = p
	}

	return &v1.LoginRes{
		Token: loginRes.Token,
		UserInfo: &v1.LoginUser{
			Id:       loginRes.UserInfo.Id,
			Username: loginRes.UserInfo.Username,
			Nickname: loginRes.UserInfo.Nickname,
			Avatar:   loginRes.UserInfo.Avatar,
			Permissions: perms,
		},
	}, nil
}

func bearerToken(r *ghttp.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
}
