package auth

import (
	"context"
	auditsvc "exam/internal/service/audit"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "exam/api/client/auth/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	secsvc "exam/internal/service/security"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	ip, userAgent := "", ""
	if r != nil {
		ip = r.GetClientIp()
		userAgent = r.Header.Get("User-Agent")
	}
	if r != nil {
		tok := bearerTokenClient(r)
		if d := middleware.GetCtxData(ctx); d != nil && tok != "" {
			auditsvc.Audit().RecordLogout(ctx, d.UserId, d.Username, d.UserType, ip, userAgent, middleware.GetTraceId(ctx))
			if err := secsvc.Security().RevokeToken(ctx, consts.UserTypeClient, d.UserId, tok); err != nil {
				g.Log().Warningf(ctx, "revoke client token failed: %v", err)
			}
		}
	}
	return &v1.LogoutRes{}, nil
}
