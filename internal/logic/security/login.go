package security

import (
	"context"
	"strings"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	sysdo "exam/internal/model/do/sys"
	"exam/internal/service/audit"
	membersvc "exam/internal/service/member"
	rolesvc "exam/internal/service/sysrole"
	"exam/internal/service/sysuser"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type loginAccount struct {
	Id                int64
	Username          string
	PasswordHash      string
	Nickname          string
	Avatar            string
	Status            int
	PasswordChangedAt *gtime.Time
}

func (s *sSecurity) Login(ctx context.Context, input bo.LoginInput) (*bo.LoginResult, error) {
	if s.CheckIPLoginRateLimit(ctx, input.IP) {
		s.recordSuspiciousIP(ctx, input)
		return nil, gerror.NewCode(consts.CodeTooManyRequests)
	}

	loginName := s.NormalizeLoginName(input.Username)
	if s.ShouldRequireCaptcha(ctx, input.UserType, loginName) {
		if input.CaptchaId == "" || input.CaptchaAnswer == "" {
			return nil, gerror.NewCode(consts.CodeCaptchaRequired)
		}
		if !s.VerifyCaptcha(ctx, input.CaptchaId, input.CaptchaAnswer) {
			return nil, gerror.NewCode(consts.CodeCaptchaInvalid)
		}
	}
	if s.IsAccountLocked(ctx, input.UserType, loginName) {
		return nil, gerror.NewCode(consts.CodeAccountLocked)
	}

	account, err := s.loadLoginAccount(ctx, input.UserType, loginName)
	if err != nil {
		g.Log().Errorf(ctx, "load login account failed: userType=%d, username=%s, err=%v", input.UserType, loginName, err)
		s.recordLoginFailure(ctx, input, 0, "user lookup failed", false)
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}
	if account == nil {
		s.recordLoginFailure(ctx, input, 0, "user not found", true)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}
	if account.Status == consts.StatusDisabled {
		s.recordLoginFailure(ctx, input, account.Id, "user disabled", false)
		return nil, gerror.NewCode(consts.CodeUserDisabled)
	}
	if s.IsPasswordExpired(ctx, account.PasswordChangedAt) {
		s.recordLoginFailure(ctx, input, account.Id, "password expired", false)
		return nil, gerror.NewCode(consts.CodePasswordExpired)
	}

	plainPassword, err := s.resolveEncryptedPassword(ctx, input.EncryptedPassword)
	if err != nil {
		s.recordLoginFailure(ctx, input, account.Id, "invalid encrypted password", true)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}
	if !utility.CheckPassword(account.PasswordHash, plainPassword) {
		s.recordLoginFailure(ctx, input, account.Id, "invalid password", true)
		return nil, gerror.NewCode(consts.CodeInvalidCredentials)
	}

	permissions, err := s.loadLoginPermissions(ctx, input.UserType, account.Id)
	if err != nil {
		g.Log().Errorf(ctx, "load login permissions failed: userType=%d, userId=%d, err=%v", input.UserType, account.Id, err)
		s.recordLoginFailure(ctx, input, account.Id, "load permissions failed", false)
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}

	token, err := s.IssueToken(ctx, input.UserType, account.Id, account.Username)
	if err != nil {
		g.Log().Errorf(ctx, "issue login token failed: userType=%d, userId=%d, err=%v", input.UserType, account.Id, err)
		s.recordLoginFailure(ctx, input, account.Id, "issue token failed", false)
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}

	s.bestEffortUpdateLoginMeta(ctx, input.UserType, account.Id, input.IP)
	s.ClearLoginFailure(ctx, input.UserType, input.Username)
	audit.Audit().RecordLoginSuccess(ctx, account.Id, account.Username, input.UserType, input.IP, input.UserAgent, input.TraceId)

	return &bo.LoginResult{
		Token: token,
		UserInfo: bo.LoginUserInfo{
			Id:          account.Id,
			Username:    account.Username,
			Nickname:    account.Nickname,
			Avatar:      account.Avatar,
			Permissions: permissions,
		},
	}, nil
}

func (s *sSecurity) resolveEncryptedPassword(ctx context.Context, encryptedPassword string) (string, error) {
	if strings.TrimSpace(encryptedPassword) == "" {
		return "", gerror.NewCode(consts.CodeInvalidParams)
	}
	return s.DecryptLoginPassword(ctx, encryptedPassword)
}

func (s *sSecurity) loadLoginAccount(ctx context.Context, userType int, username string) (*loginAccount, error) {
	switch userType {
	case consts.UserTypeClient:
		member, err := membersvc.Member().FindByUsername(ctx, username)
		if err != nil || member == nil {
			return nil, err
		}
		return &loginAccount{
			Id:                member.Id,
			Username:          member.Username,
			PasswordHash:      member.Password,
			Nickname:          member.Nickname,
			Avatar:            member.Avatar,
			Status:            member.Status,
			PasswordChangedAt: member.PasswordChangedAt,
		}, nil
	default:
		user, err := sysuser.SysUser().FindByUsername(ctx, username)
		if err != nil || user == nil {
			return nil, err
		}
		return &loginAccount{
			Id:                user.Id,
			Username:          user.Username,
			PasswordHash:      user.Password,
			Nickname:          user.Nickname,
			Avatar:            user.Avatar,
			Status:            user.Status,
			PasswordChangedAt: user.PasswordChangedAt,
		}, nil
	}
}

func (s *sSecurity) recordSuspiciousIP(ctx context.Context, input bo.LoginInput) {
	audit.Audit().RecordSecurityEvent(ctx, consts.EventTypeSuspiciousIP, 0, input.IP, input.UserAgent, "login rate limit exceeded", input.TraceId)
}

func (s *sSecurity) recordLoginFailure(ctx context.Context, input bo.LoginInput, userId int64, reason string, increaseFailureCount bool) {
	audit.Audit().RecordLoginFailure(ctx, userId, input.Username, input.UserType, input.IP, input.UserAgent, reason, input.TraceId)
	if increaseFailureCount {
		s.RecordLoginFailure(ctx, input.UserType, input.Username, input.IP, input.UserAgent, input.TraceId)
	}
}

func (s *sSecurity) loadLoginPermissions(ctx context.Context, userType int, userId int64) ([]string, error) {
	if userType != consts.UserTypeAdmin || userId <= 0 {
		return nil, nil
	}
	return rolesvc.SysRole().PermissionCodesByUser(ctx, userId)
}

func (s *sSecurity) bestEffortUpdateLoginMeta(ctx context.Context, userType int, userId int64, ip string) {
	if userId <= 0 {
		return
	}
	now := gtime.Now()
	var err error
	switch userType {
	case consts.UserTypeClient:
		_, err = dao.SysMember.Ctx(ctx).Where("id", userId).Data(sysdo.SysMember{
			LoginIp:   ip,
			LoginTime: now,
		}).Update()
	default:
		_, err = dao.SystemUser.Ctx(ctx).Where("id", userId).Data(sysdo.SysUser{
			LoginIp:   ip,
			LoginTime: now,
		}).Update()
	}
	if err != nil {
		g.Log().Warningf(ctx, "update login meta failed: userType=%d, userId=%d, err=%v", userType, userId, err)
	}
}
