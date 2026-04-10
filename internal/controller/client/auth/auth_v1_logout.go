package auth

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "exam/api/client/auth/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	secsvc "exam/internal/service/security"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r != nil {
		tok := bearerTokenClient(r)
		if tok != "" {
			key := consts.TokenRedisKeyPrefix + consts.UserTypeTagClient + ":" + tok
			_, _ = g.Redis().Del(ctx, key)
		}
		if d := middleware.GetCtxData(ctx); d != nil && tok != "" {
			secsvc.Security().RemoveSessionToken(ctx, consts.UserTypeClient, d.UserId, tok)
		}
	}
	return &v1.LogoutRes{}, nil
}
