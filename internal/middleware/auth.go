package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
)

type tokenPayload struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
}

func userTypeTag(ut int) string {
	if ut == consts.UserTypeClient {
		return "client"
	}
	return "admin"
}

// Auth 校验 Bearer Token 并注入 CtxData。
func Auth(expectedUserType int) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		raw := strings.TrimSpace(r.Header.Get("Authorization"))
		token := strings.TrimPrefix(raw, "Bearer ")
		token = strings.TrimSpace(token)
		if token == "" {
			r.SetError(gerror.NewCode(consts.CodeTokenRequired))
			r.ExitAll()
			return
		}
		key := consts.TokenRedisKeyPrefix + userTypeTag(expectedUserType) + ":" + token
		val, err := g.Redis().Get(r.GetCtx(), key)
		if err != nil || val.IsEmpty() {
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid))
			r.ExitAll()
			return
		}
		var p tokenPayload
		if json.Unmarshal([]byte(val.String()), &p) != nil || p.UserId == 0 {
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid))
			r.ExitAll()
			return
		}
		ctx := SetCtxData(r.GetCtx(), &CtxData{
			UserId:   p.UserId,
			UserType: expectedUserType,
			Username: p.Username,
		})
		r.SetCtx(ctx)
		r.Middleware.Next()
	}
}
