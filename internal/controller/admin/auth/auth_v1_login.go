package auth

import (
	"context"
	"encoding/json"
	auditsvc "exam/internal/service/audit"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"

	v1 "exam/api/admin/auth/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	secsvc "exam/internal/service/security"
	usersvc "exam/internal/service/sysuser"
	"exam/internal/utility"
)

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	ip, userAgent := "", ""
	if r != nil {
		ip = r.GetClientIp()
		userAgent = r.Header.Get("User-Agent")
	}
	traceId := middleware.GetTraceId(ctx)
	if secsvc.Security().CheckIPLoginRateLimit(ctx, ip) {
		auditsvc.Audit().RecordSecurityEvent(ctx, consts.EventTypeSuspiciousIP, 0, ip, userAgent, "login rate limit exceeded", traceId)
		return nil, gerror.NewCode(consts.CodeTooManyRequests)
	}
	name := secsvc.Security().NormalizeLoginName(req.Username)
	if secsvc.Security().ShouldRequireCaptcha(ctx, consts.UserTypeAdmin, name) {
		if req.CaptchaId == "" || !secsvc.Security().VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaAnswer) {
			return nil, gerror.NewCode(consts.CodeCaptchaRequired)
		}
	}
	if secsvc.Security().IsAccountLocked(ctx, consts.UserTypeAdmin, name) {
		return nil, gerror.NewCode(consts.CodeAccountLocked)
	}

	u, _ := usersvc.SysUser().FindByUsername(ctx, name)

	if u == nil {
		auditsvc.Audit().RecordLoginFailure(ctx, 0, req.Username, consts.UserTypeAdmin, ip, userAgent, "user not found", traceId)
		secsvc.Security().RecordLoginFailure(ctx, consts.UserTypeAdmin, req.Username, ip, userAgent, traceId)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}

	if u.Status == consts.StatusDisabled {
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeAdmin, ip, userAgent, "user disabled", traceId)
		return nil, gerror.NewCode(consts.CodeUserDisabled)
	}
	plainPassword, err := resolveLoginPassword(ctx, req.Password)
	if err != nil {
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeAdmin, ip, userAgent, "invalid encrypted password", traceId)
		secsvc.Security().RecordLoginFailure(ctx, consts.UserTypeAdmin, req.Username, ip, userAgent, traceId)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}
	if !utility.CheckPassword(u.Password, plainPassword) {
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeAdmin, ip, userAgent, "invalid password", traceId)
		secsvc.Security().RecordLoginFailure(ctx, consts.UserTypeAdmin, req.Username, ip, userAgent, traceId)
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
	key := consts.TokenRedisKeyPrefix + consts.UserTypeTagAdmin + ":" + token
	if err := g.Redis().SetEX(ctx, key, string(payload), ttl); err != nil {
		g.Log().Errorf(ctx, "redis set token failed: %v", err)
		auditsvc.Audit().RecordLoginFailure(ctx, u.Id, req.Username, consts.UserTypeAdmin, ip, userAgent, "redis set token failed", traceId)
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}
	_ = secsvc.Security().RegisterSession(ctx, consts.UserTypeAdmin, u.Id, token, ttl)

	secsvc.Security().ClearLoginFailure(ctx, consts.UserTypeAdmin, req.Username)
	auditsvc.Audit().RecordLoginSuccess(ctx, u.Id, u.Username, consts.UserTypeAdmin, ip, userAgent, traceId)

	var perms []string
	if p, perr := middleware.GetUserPermissions(ctx, u.Id); perr == nil {
		perms = p
	}

	return &v1.LoginRes{
		Token: token,
		UserInfo: &v1.LoginUser{
			Id: u.Id, Username: u.Username, Nickname: u.Nickname, Avatar: u.Avatar,
			Permissions: perms,
		},
	}, nil
}

func bearerToken(r *ghttp.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
}

func resolveLoginPassword(ctx context.Context, encryptedPassword string) (string, error) {
	if strings.TrimSpace(encryptedPassword) == "" {
		return "", gerror.NewCode(consts.CodeInvalidParams)
	}
	return secsvc.Security().DecryptLoginPassword(ctx, encryptedPassword)
}
