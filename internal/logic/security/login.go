package security

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
	"exam/internal/service/audit"
)

func userTypeTag(ut int) string {
	if ut == consts.UserTypeClient {
		return consts.UserTypeTagClient
	}
	return consts.UserTypeTagAdmin
}

// NormalizeLoginName 登录名规范化（用于 Redis 键）
func (s *sSecurity) NormalizeLoginName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// CheckIPLoginRateLimit 单 IP 每分钟尝试次数，超限返回 true
func (s *sSecurity) CheckIPLoginRateLimit(ctx context.Context, ip string) (blocked bool) {
	cfg := s.LoadLoginCfg(ctx)
	if cfg.RateLimitPerMinute <= 0 || ip == "" {
		return false
	}
	key := consts.LoginRateLimitKeyPrefix + ip
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
func (s *sSecurity) IsAccountLocked(ctx context.Context, userType int, username string) bool {
	name := s.NormalizeLoginName(username)
	if name == "" {
		return false
	}
	key := consts.LoginLockKeyPrefix + userTypeTag(userType) + ":" + name
	n, err := g.Redis().Get(ctx, key)
	if err != nil || n.IsEmpty() {
		return false
	}
	return true
}

// ShouldRequireCaptcha 是否必须提交验证码（失败次数达到阈值）
func (s *sSecurity) ShouldRequireCaptcha(ctx context.Context, userType int, username string) bool {
	cfg := s.LoadLoginCfg(ctx)
	if !cfg.CaptchaEnabled || cfg.CaptchaAfterFailures <= 0 {
		return false
	}
	name := s.NormalizeLoginName(username)
	if name == "" {
		return false
	}
	key := consts.LoginFailCountKeyPrefix + userTypeTag(userType) + ":" + name
	val, err := g.Redis().Get(ctx, key)
	if err != nil || val.IsEmpty() {
		return false
	}
	return val.Int() >= cfg.CaptchaAfterFailures
}

// RecordLoginFailure 记录失败：计数、可能加锁并记录 brute_force 安全事件
func (s *sSecurity) RecordLoginFailure(ctx context.Context, userType int, username, ip, userAgent, traceId string) {
	cfg := s.LoadLoginCfg(ctx)
	name := s.NormalizeLoginName(username)
	if name == "" {
		return
	}
	failKey := consts.LoginFailCountKeyPrefix + userTypeTag(userType) + ":" + name
	n, err := g.Redis().Incr(ctx, failKey)
	if err != nil {
		return
	}
	if n == 1 && cfg.FailureWindowSeconds > 0 {
		_, _ = g.Redis().Expire(ctx, failKey, int64(cfg.FailureWindowSeconds))
	}

	if cfg.MaxFailuresBeforeLock > 0 && int(n) >= cfg.MaxFailuresBeforeLock {
		lockKey := consts.LoginLockKeyPrefix + userTypeTag(userType) + ":" + name
		if cfg.LockDurationSeconds > 0 {
			_ = g.Redis().SetEX(ctx, lockKey, "1", int64(cfg.LockDurationSeconds))
		}
		audit.Audit().RecordSecurityEvent(ctx, consts.SecurityEventBruteForce, 0, ip, userAgent,
			fmt.Sprintf("account locked after %d failed logins: %s", n, name), traceId)
	}
}

// ClearLoginFailure 登录成功时清除失败计数与锁定（若存在可配置为手动解锁，此处仅清失败计数）
func (s *sSecurity) ClearLoginFailure(ctx context.Context, userType int, username string) {
	name := s.NormalizeLoginName(username)
	if name == "" {
		return
	}
	failKey := consts.LoginFailCountKeyPrefix + userTypeTag(userType) + ":" + name
	_, _ = g.Redis().Del(ctx, failKey)
}

// UnlockAccount 管理员或到期后解锁（删除锁定键）
func (s *sSecurity) UnlockAccount(ctx context.Context, userType int, username string) {
	name := s.NormalizeLoginName(username)
	if name == "" {
		return
	}
	lockKey := consts.LoginLockKeyPrefix + userTypeTag(userType) + ":" + name
	_, _ = g.Redis().Del(ctx, lockKey)
	failKey := consts.LoginFailCountKeyPrefix + userTypeTag(userType) + ":" + name
	_, _ = g.Redis().Del(ctx, failKey)
}
