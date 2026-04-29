package security

import (
	"context"
	"fmt"
	"strings"
	"time"

	"exam/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

func (s *sSecurity) CheckForgetPasswordAccess(ctx context.Context, ip string, username string) bool {
	day := time.Now().Format("20060102")
	normalizedUser := strings.ToLower(strings.TrimSpace(username))

	if ip != "" && s.isForgetPasswordBlocked(ctx, fmt.Sprintf(consts.ForgetPasswordBlockedIPKeyFmt, ip)) {
		return true
	}
	if normalizedUser != "" && s.isForgetPasswordBlocked(ctx, fmt.Sprintf(consts.ForgetPasswordBlockedUserKeyFmt, normalizedUser)) {
		return true
	}

	if ip != "" && s.increaseForgetPasswordDailyCounter(ctx,
		fmt.Sprintf(consts.ForgetPasswordDailyIPCounterKeyFmt, ip, day),
		fmt.Sprintf(consts.ForgetPasswordBlockedIPKeyFmt, ip),
	) {
		return true
	}
	if normalizedUser != "" && s.increaseForgetPasswordDailyCounter(ctx,
		fmt.Sprintf(consts.ForgetPasswordDailyUserCounterKeyFmt, normalizedUser, day),
		fmt.Sprintf(consts.ForgetPasswordBlockedUserKeyFmt, normalizedUser),
	) {
		return true
	}
	return false
}

func (s *sSecurity) CheckForgetPasswordCooldown(ctx context.Context, username string) bool {
	normalizedUser := strings.ToLower(strings.TrimSpace(username))
	if normalizedUser == "" {
		return false
	}
	key := fmt.Sprintf(consts.ForgetPasswordCooldownUserKeyFmt, normalizedUser)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", consts.ForgetPasswordCooldownTTLSeconds)
	if err != nil {
		return false
	}
	return v == nil || strings.ToUpper(v.String()) != "OK"
}

func (s *sSecurity) isForgetPasswordBlocked(ctx context.Context, key string) bool {
	v, err := g.Redis().Get(ctx, key)
	return err == nil && !v.IsEmpty()
}

func (s *sSecurity) increaseForgetPasswordDailyCounter(ctx context.Context, counterKey string, blockedKey string) bool {
	n, err := g.Redis().Incr(ctx, counterKey)
	if err != nil {
		return false
	}
	if n == 1 {
		_, _ = g.Redis().Expire(ctx, counterKey, consts.ForgetPasswordDailyCounterMaxTTL)
	}
	if int(n) > consts.ForgetPasswordDailyAccessLimit {
		_ = g.Redis().SetEX(ctx, blockedKey, "1", consts.ForgetPasswordBlockedTTLSeconds)
		return true
	}
	return false
}
