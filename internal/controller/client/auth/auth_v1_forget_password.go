package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "exam/api/client/auth/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysentity "exam/internal/model/entity/sys"
	auditsvc "exam/internal/service/audit"
	membersvc "exam/internal/service/member"
	secsvc "exam/internal/service/security"
	notisvc "exam/internal/service/sysnotification"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) ForgetPassword(ctx context.Context, req *v1.ForgetPasswordReq) (res *v1.ForgetPasswordRes, err error) {
	email := strings.TrimSpace(req.Email)
	if email == "" {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	httpReq := ghttp.RequestFromCtx(ctx)
	ip, userAgent := "", ""
	if httpReq != nil {
		ip = httpReq.GetClientIp()
		userAgent = httpReq.Header.Get("User-Agent")
	}
	traceId := middleware.GetTraceId(ctx)
	if blocked := secsvc.Security().CheckForgetPasswordAccess(ctx, ip, email); blocked {
		auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_rate_limited", 0, ip, userAgent, fmt.Sprintf("email=%s", email), traceId)
		return &v1.ForgetPasswordRes{}, nil
	}

	member, err := membersvc.Member().FindByUsername(ctx, email)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	// Avoid email enumeration: return success when account does not exist.
	if member == nil {
		auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_user_not_found", 0, ip, userAgent, fmt.Sprintf("email=%s", email), traceId)
		return &v1.ForgetPasswordRes{}, nil
	}
	if member.Status == consts.StatusDisabled {
		auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_user_disabled", member.Id, ip, userAgent, fmt.Sprintf("email=%s", email), traceId)
		return &v1.ForgetPasswordRes{}, nil
	}
	if coolingDown := secsvc.Security().CheckForgetPasswordCooldown(ctx, email); coolingDown {
		auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_cooldown", member.Id, ip, userAgent, fmt.Sprintf("email=%s", email), traceId)
		return &v1.ForgetPasswordRes{}, nil
	}

	password, err := generateResetPassword(ctx)
	if err != nil {
		return nil, err
	}
	if err = membersvc.Member().MemberUpdatePwd(ctx, member.Id, password); err != nil {
		return nil, err
	}
	auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_password_reset", member.Id, ip, userAgent, fmt.Sprintf("trigger=forget_password email=%s", email), traceId)

	template, err := loadActiveTemplateByCode(ctx, "forget_password")
	if err != nil {
		return nil, err
	}
	recipient := recipientByChannel(member, template.Channel)
	if recipient == "" {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	vars, _ := json.Marshal(map[string]string{
		"password": password,
	})
	if _, err = notisvc.SysNotification().Send(ctx, "forget_password", template.Channel, recipient, string(vars)); err != nil {
		g.Log().Errorf(ctx, "[forget_password] notification send failed email=%s template=forget_password channel=%s recipient=%s err=%v", email, template.Channel, recipient, err)
		auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_notify_failed", member.Id, ip, userAgent, fmt.Sprintf("channel=%s recipient=%s", template.Channel, recipient), traceId)
		return &v1.ForgetPasswordRes{}, nil
	}
	auditsvc.Audit().RecordSecurityEvent(ctx, "forget_password_notify_sent", member.Id, ip, userAgent, fmt.Sprintf("channel=%s recipient=%s", template.Channel, recipient), traceId)
	return &v1.ForgetPasswordRes{}, nil
}

func loadActiveTemplateByCode(ctx context.Context, code string) (*sysentity.SysNotificationTemplate, error) {
	var tpl sysentity.SysNotificationTemplate
	if err := dao.SysNotificationTemplate.Ctx(ctx).
		Where("code", code).
		Where("status", 0).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&tpl); err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if tpl.Id == 0 {
		return nil, gerror.NewCode(consts.CodeTemplateNotFound)
	}
	return &tpl, nil
}

func recipientByChannel(member *sysentity.SysMember, channel string) string {
	switch channel {
	case "email":
		return strings.TrimSpace(member.Email)
	case "sms":
		return strings.TrimSpace(member.Mobile)
	default:
		return ""
	}
}

func generateResetPassword(ctx context.Context) (string, error) {
	cfg := secsvc.Security().LoadPasswordCfg(ctx)
	password, err := utility.GeneratePasswordByPolicy(cfg)
	if err != nil {
		return "", err
	}
	if err = secsvc.Security().ValidatePasswordPolicy(ctx, password); err != nil {
		return "", gerror.NewCode(consts.CodePasswordWeak)
	}
	return password, nil
}
