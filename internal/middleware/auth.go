package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
	auditsvc "exam/internal/service/audit"
)

type tokenPayload struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
}

func userTypeTag(ut int) string {
	if ut == consts.UserTypeClient {
		return consts.UserTypeTagClient
	}
	return consts.UserTypeTagAdmin
}

func Auth(expectedUserType int) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		var (
			ctx     = r.GetCtx()
			ip      = r.GetClientIp()
			ua      = r.Header.Get("User-Agent")
			traceId = GetTraceId(ctx)
			raw     = strings.TrimSpace(r.Header.Get("Authorization"))
			token   = strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
			uType   = userTypeTag(expectedUserType)
		)

		// 1. 检查 Token 是否为空
		if token == "" {
			g.Log().Warningf(ctx, "[Auth] Missing Token | IP: %s | URL: %s", ip, r.RequestURI) // 打印告警日志
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "missing authorization header", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenRequired, "请提供身份认证凭证"))
			r.ExitAll()
			return
		}

		// 2. 从 Redis 获取 Token 数据
		key := consts.TokenRedisKeyPrefix + uType + ":" + token
		val, err := g.Redis().Get(ctx, key)
		if err != nil {
			g.Log().Errorf(ctx, "[Auth] Redis Error: %v | Key: %s", err, key) // 打印错误日志
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "redis lookup failed", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid, "服务繁忙"))
			r.ExitAll()
			return
		}

		// 3. 检查 Token 是否过期或不存在
		if val.IsEmpty() {
			// 这里使用 Debug 或 Info，因为 Token 过期是正常业务现象，不是系统故障
			g.Log().Infof(ctx, "[Auth] Token Expired or Invalid | Type: %s | IP: %s", uType, ip)
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "token expired or not found", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid, "登录已过期"))
			r.ExitAll()
			return
		}

		// 4. 解析 Payload
		var p tokenPayload
		if err := json.Unmarshal(val.Bytes(), &p); err != nil || p.UserId == 0 {
			g.Log().Errorf(ctx, "[Auth] Payload Parse Failed | Content: %s | Error: %v", val.String(), err)
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "malformed token payload", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid))
			r.ExitAll()
			return
		}

		// 5. 鉴权通过日志（可选，通常设为 Debug 级别防止日志量过大）
		g.Log().Debugf(ctx, "[Auth] Success | UserID: %d | UserType: %s | IP: %s", p.UserId, uType, ip)

		// 注入上下文并继续
		newCtx := SetCtxData(ctx, &CtxData{
			UserId:   p.UserId,
			UserType: expectedUserType,
			Username: p.Username,
		})
		r.SetCtx(newCtx)
		r.Middleware.Next()
	}
}
