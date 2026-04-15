package auth

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"

	v1 "exam/api/client/auth/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	auditsvc "exam/internal/service/audit"
	membersvc "exam/internal/service/member"
	secsvc "exam/internal/service/security"
	"exam/internal/utility"
)

func bearerTokenClient(r *ghttp.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	httpReq := ghttp.RequestFromCtx(ctx)
	ip, userAgent := "", ""
	if httpReq != nil {
		ip = httpReq.GetClientIp()
		userAgent = httpReq.Header.Get("User-Agent")
	}
	traceId := middleware.GetTraceId(ctx)
	if secsvc.Security().CheckIPLoginRateLimit(ctx, ip) {
		auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeSuspiciousIP, 0, ip, userAgent, "login rate limit exceeded", traceId)
		return nil, gerror.NewCode(consts.CodeTooManyRequests)
	}
	name := secsvc.Security().NormalizeLoginName(req.Username)
	if secsvc.Security().ShouldRequireCaptcha(ctx, consts.UserTypeClient, name) {
		if req.CaptchaId == "" || !secsvc.Security().VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaAnswer) {
			return nil, gerror.NewCode(consts.CodeCaptchaRequired)
		}
	}
	if secsvc.Security().IsAccountLocked(ctx, consts.UserTypeClient, name) {
		return nil, gerror.NewCode(consts.CodeAccountLocked)
	}

	u, _ := membersvc.Member().FindByUsername(ctx, name)
	if u == nil {
		auditsvc.Audit().RecordLoginFailure(ctx, 0, req.Username, consts.UserTypeClient, ip, userAgent, "user not found", traceId)
		secsvc.Security().RecordLoginFailure(ctx, consts.UserTypeClient, req.Username, ip, userAgent, traceId)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}

	if u.Status == consts.StatusDisabled {
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeClient, ip, userAgent, "user disabled", traceId)
		return nil, gerror.NewCode(consts.CodeUserDisabled)
	}
	if !utility.CheckPassword(u.Password, req.Password) {
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeClient, ip, userAgent, "invalid password", traceId)
		secsvc.Security().RecordLoginFailure(ctx, consts.UserTypeClient, req.Username, ip, userAgent, traceId)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}

	token := guid.S()
	ttl := secsvc.Security().TokenTTLSeconds(ctx)
	if ttl <= 0 {
		ttl = consts.DefaultTokenTTLFallbackSeconds
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"user_id": u.Id, "username": u.Username,
	})
	key := consts.TokenRedisKeyPrefix + consts.UserTypeTagClient + ":" + token
	if err := g.Redis().SetEX(ctx, key, string(payload), ttl); err != nil {
		g.Log().Errorf(ctx, "redis set token failed: %v", err)
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeClient, ip, userAgent, "redis set token failed", traceId)
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}
	_ = secsvc.Security().RegisterSession(ctx, consts.UserTypeClient, u.Id, token, ttl)

	secsvc.Security().ClearLoginFailure(ctx, consts.UserTypeClient, req.Username)
	auditsvc.Audit().RecordLoginSuccess(ctx, u.Id, u.Username, consts.UserTypeClient, ip, userAgent, traceId)

	return &v1.LoginRes{
		Token: token,
		UserInfo: &v1.LoginUser{
			Id: u.Id, Username: u.Username, Nickname: u.Nickname,
		},
	}, nil
}
