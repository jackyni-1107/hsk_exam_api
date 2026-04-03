// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package security

import (
	"context"
	"exam/internal/model/bo"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	ISecurity interface {
		// CreateCaptcha 生成验证码并存 Redis
		CreateCaptcha(ctx context.Context) (*bo.CaptchaChallenge, error)
		// VerifyCaptcha 校验答案（一次性）
		VerifyCaptcha(ctx context.Context, captchaId string, answer string) bool
		LoadLoginCfg(ctx context.Context) bo.LoginCfg
		LoadPasswordCfg(ctx context.Context) bo.PasswordCfg
		LoadSessionCfg(ctx context.Context) bo.SessionCfg
		LoadMFACfg(ctx context.Context) bo.MFACfg
		// TokenTTLSeconds 会话 Token 有效期（秒）
		TokenTTLSeconds(ctx context.Context) int64
		// NormalizeLoginName 登录名规范化（用于 Redis 键）
		NormalizeLoginName(name string) string
		// CheckIPLoginRateLimit 单 IP 每分钟尝试次数，超限返回 true
		CheckIPLoginRateLimit(ctx context.Context, ip string) (blocked bool)
		// IsAccountLocked 账号是否处于锁定窗口
		IsAccountLocked(ctx context.Context, userType int, username string) bool
		// ShouldRequireCaptcha 是否必须提交验证码（失败次数达到阈值）
		ShouldRequireCaptcha(ctx context.Context, userType int, username string) bool
		// RecordLoginFailure 记录失败：计数、可能加锁并记录 brute_force 安全事件
		RecordLoginFailure(ctx context.Context, userType int, username string, ip string, userAgent string, traceId string)
		// ClearLoginFailure 登录成功时清除失败计数与锁定（若存在可配置为手动解锁，此处仅清失败计数）
		ClearLoginFailure(ctx context.Context, userType int, username string)
		// UnlockAccount 管理员或到期后解锁（删除锁定键）
		UnlockAccount(ctx context.Context, userType int, username string)
		// MFANotImplemented 预留：等保/高安全场景可接入 TOTP、WebAuthn 等。
		// 配置 security.mfa.enabled 为 true 时，业务侧应在此处实现校验并在登录流程中串联。
		MFANotImplemented(ctx context.Context) bool
		// IsPasswordExpired 是否超过口令最长使用期限（maxAgeDays<=0 表示不启用）
		IsPasswordExpired(ctx context.Context, passwordChangedAt *gtime.Time) bool
		// ValidatePasswordNotInHistory 新口令不能与近期历史相同
		ValidatePasswordNotInHistory(ctx context.Context, userType int, userId int64, plainPassword string) error
		// SavePasswordHistory 在更新密码前将旧哈希写入历史表，并裁剪超出条数
		SavePasswordHistory(ctx context.Context, userType int, userId int64, oldHash string) error
		// ValidatePasswordPolicy 按配置校验口令复杂度
		ValidatePasswordPolicy(ctx context.Context, password string) error
		// RegisterSession 登记新 Token，超出并发数时剔除最旧会话
		RegisterSession(ctx context.Context, userType int, userId int64, token string, ttlSeconds int64) error
		// RemoveSessionToken 登出时从会话列表移除
		RemoveSessionToken(ctx context.Context, userType int, userId int64, token string)
		// RevokeAllUserSessions 撤销某用户全部 Token（强制下线）
		RevokeAllUserSessions(ctx context.Context, userType int, userId int64) error
	}
)

var (
	localSecurity ISecurity
)

func Security() ISecurity {
	if localSecurity == nil {
		panic("implement not found for interface ISecurity, forgot register?")
	}
	return localSecurity
}

func RegisterSecurity(i ISecurity) {
	localSecurity = i
}
