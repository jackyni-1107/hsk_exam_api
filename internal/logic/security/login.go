package security

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/service/audit"
)

const (
	loginRLPrefix   = "login:rl:"
	loginFailPrefix = "login:fail:"
	loginLockPrefix = "login:lock:"
)

func userTypeTag(ut int) string {
	if ut == 2 {
		return "client"
	}
	return "admin"
}

// NormalizeLoginName 登录名规范化（用于 Redis 键）
func NormalizeLoginName(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// CheckIPLoginRateLimit 单 IP 每分钟尝试次数，超限返回 true
func CheckIPLoginRateLimit(ctx context.Context, ip string) (blocked bool) {
	cfg := LoadLoginCfg(ctx)
	if cfg.RateLimitPerMinute <= 0 || ip == "" {
		return false
	}
	key := loginRLPrefix + ip
	n, err := g.Redis().Incr(ctx, key)
	if err != nil {
		return false
	}
	if n == 1 {
		_, _ = g.Redis().Expire(ctx, key, 60)
	}
	return int(n) > cfg.RateLimitPerMinute
}

// IsAccountLocked 账号是否处于锁定窗口
func IsAccountLocked(ctx context.Context, userType int, username string) bool {
	name := NormalizeLoginName(username)
	if name == "" {
		return false
	}
	key := loginLockPrefix + userTypeTag(userType) + ":" + name
	n, err := g.Redis().Get(ctx, key)
	if err != nil || n.IsEmpty() {
		return false
	}
	return true
}

// ShouldRequireCaptcha 是否必须提交验证码（失败次数达到阈值）
func ShouldRequireCaptcha(ctx context.Context, userType int, username string) bool {
	cfg := LoadLoginCfg(ctx)
	if !cfg.CaptchaEnabled || cfg.CaptchaAfterFailures <= 0 {
		return false
	}
	name := NormalizeLoginName(username)
	if name == "" {
		return false
	}
	key := loginFailPrefix + userTypeTag(userType) + ":" + name
	val, err := g.Redis().Get(ctx, key)
	if err != nil || val.IsEmpty() {
		return false
	}
	return val.Int() >= cfg.CaptchaAfterFailures
}

// RecordLoginFailure 记录失败：计数、可能加锁并记录 brute_force 安全事件
func RecordLoginFailure(ctx context.Context, userType int, username, ip, userAgent, traceId string) {
	cfg := LoadLoginCfg(ctx)
	name := NormalizeLoginName(username)
	if name == "" {
		return
	}
	failKey := loginFailPrefix + userTypeTag(userType) + ":" + name
	n, err := g.Redis().Incr(ctx, failKey)
	if err != nil {
		return
	}
	if n == 1 && cfg.FailureWindowSeconds > 0 {
		_, _ = g.Redis().Expire(ctx, failKey, int64(cfg.FailureWindowSeconds))
	}

	if cfg.MaxFailuresBeforeLock > 0 && int(n) >= cfg.MaxFailuresBeforeLock {
		lockKey := loginLockPrefix + userTypeTag(userType) + ":" + name
		if cfg.LockDurationSeconds > 0 {
			_ = g.Redis().SetEX(ctx, lockKey, "1", int64(cfg.LockDurationSeconds))
		}
		audit.Audit().RecordSecurityEvent(ctx, "brute_force", 0, ip, userAgent,
			fmt.Sprintf("account locked after %d failed logins: %s", n, name), traceId)
	}
}

// ClearLoginFailure 登录成功时清除失败计数与锁定（若存在可配置为手动解锁，此处仅清失败计数）
func ClearLoginFailure(ctx context.Context, userType int, username string) {
	name := NormalizeLoginName(username)
	if name == "" {
		return
	}
	failKey := loginFailPrefix + userTypeTag(userType) + ":" + name
	_, _ = g.Redis().Del(ctx, failKey)
}

// UnlockAccount 管理员或到期后解锁（删除锁定键）
func UnlockAccount(ctx context.Context, userType int, username string) {
	name := NormalizeLoginName(username)
	if name == "" {
		return
	}
	lockKey := loginLockPrefix + userTypeTag(userType) + ":" + name
	_, _ = g.Redis().Del(ctx, lockKey)
	failKey := loginFailPrefix + userTypeTag(userType) + ":" + name
	_, _ = g.Redis().Del(ctx, failKey)
}
