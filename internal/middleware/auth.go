package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
	auditsvc "exam/internal/service/audit"
	secsvc "exam/internal/service/security"
)

func authUserTypeTag(ut int) string {
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
			uType   = authUserTypeTag(expectedUserType)
		)

		if token == "" {
			g.Log().Warningf(ctx, "[Auth] Missing Token | IP: %s | URL: %s", ip, r.RequestURI)
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "missing authorization header", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenRequired, "authorization token required"))
			r.ExitAll()
			return
		}

		payload, err := secsvc.Security().LoadTokenPayload(ctx, expectedUserType, token)
		if err != nil {
			g.Log().Errorf(ctx, "[Auth] Token Lookup Failed | Type: %s | Error: %v", uType, err)
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "token lookup failed", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid, "service unavailable"))
			r.ExitAll()
			return
		}
		if payload == nil {
			g.Log().Infof(ctx, "[Auth] Token Expired or Invalid | Type: %s | IP: %s", uType, ip)
			auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeTokenInvalid, 0, ip, ua, "token expired or not found", traceId)
			r.SetError(gerror.NewCode(consts.CodeTokenInvalid, "token expired"))
			r.ExitAll()
			return
		}

		g.Log().Debugf(ctx, "[Auth] Success | UserID: %d | UserType: %s | IP: %s", payload.UserId, uType, ip)

		newCtx := SetCtxData(ctx, &CtxData{
			UserId:   payload.UserId,
			UserType: expectedUserType,
			Username: payload.Username,
		})
		r.SetCtx(newCtx)
		r.Middleware.Next()
	}
}
