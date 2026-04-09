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
	"exam/internal/dao"
	"exam/internal/middleware"
	sysentity "exam/internal/model/entity/sys"
	secsvc "exam/internal/service/security"
	"exam/internal/utility"
)

func bearerTokenClient(r *ghttp.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}
	ip := r.GetClientIp()
	if secsvc.Security().CheckIPLoginRateLimit(ctx, ip) {
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

	var u sysentity.SysMember
	// 与风控键一致：规范化后查询；LOWER 避免库中用户名大小写与输入不一致
	_ = dao.SysMember.Ctx(ctx).
		Wheref("username = ?", name).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&u)
	if u.Id == 0 || !utility.CheckPassword(u.Password, req.Password) {
		secsvc.Security().RecordLoginFailure(ctx, consts.UserTypeClient, name, ip, r.Header.Get("User-Agent"), middleware.GetTraceId(ctx))
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}
	if u.Status == consts.StatusDisabled {
		return nil, gerror.NewCode(consts.CodeUserDisabled)
	}
	secsvc.Security().ClearLoginFailure(ctx, consts.UserTypeClient, name)

	token := guid.S()
	ttl := secsvc.Security().TokenTTLSeconds(ctx)
	if ttl <= 0 {
		ttl = 86400
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"user_id": u.Id, "username": u.Username,
	})
	key := consts.TokenRedisKeyPrefix + "client:" + token
	if err := g.Redis().SetEX(ctx, key, string(payload), ttl); err != nil {
		return nil, gerror.Wrap(err, "redis")
	}
	_ = secsvc.Security().RegisterSession(ctx, consts.UserTypeClient, u.Id, token, ttl)

	return &v1.LoginRes{
		Token: token,
		UserInfo: &v1.LoginUser{
			Id: u.Id, Username: u.Username, Nickname: u.Nickname,
		},
	}, nil
}
