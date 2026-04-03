package auth

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "exam/api/admin/auth/v1"
	"exam/internal/consts"
	"exam/internal/logic/security"
	"exam/internal/middleware"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r != nil {
		tok := bearerToken(r)
		if tok != "" {
			key := consts.TokenRedisKeyPrefix + "admin:" + tok
			_, _ = g.Redis().Del(ctx, key)
		}
		if d := middleware.GetCtxData(ctx); d != nil && tok != "" {
			security.RemoveSessionToken(ctx, consts.UserTypeAdmin, d.UserId, tok)
		}
	}
	return &v1.LogoutRes{}, nil
}
