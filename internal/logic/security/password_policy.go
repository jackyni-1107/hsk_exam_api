package security

import (
	"context"
	"unicode"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
)

// ValidatePasswordPolicy 按配置校验口令复杂度
func (s *sSecurity) ValidatePasswordPolicy(ctx context.Context, password string) error {
	cfg := s.LoadPasswordCfg(ctx)
	if len(password) < cfg.MinLength {
		return gerror.NewCode(consts.CodePasswordWeak)
	}
	var upper, lower, digit, special int
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			upper++
		case unicode.IsLower(r):
			lower++
		case unicode.IsDigit(r):
			digit++
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			special++
		}
	}
	if cfg.RequireUpper && upper == 0 {
		return gerror.NewCode(consts.CodePasswordWeak)
	}
	if cfg.RequireLower && lower == 0 {
		return gerror.NewCode(consts.CodePasswordWeak)
	}
	if cfg.RequireDigit && digit == 0 {
		return gerror.NewCode(consts.CodePasswordWeak)
	}
	if cfg.RequireSpecial && special == 0 {
		return gerror.NewCode(consts.CodePasswordWeak)
	}
	return nil
}
