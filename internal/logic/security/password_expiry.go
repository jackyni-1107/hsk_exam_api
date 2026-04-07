package security

import (
	"context"

	"github.com/gogf/gf/v2/os/gtime"
)

// IsPasswordExpired 是否超过口令最长使用期限（maxAgeDays<=0 表示不启用）
func (s *sSecurity) IsPasswordExpired(ctx context.Context, passwordChangedAt *gtime.Time) bool {
	cfg := s.LoadPasswordCfg(ctx)
	if cfg.MaxAgeDays <= 0 {
		return false
	}
	if passwordChangedAt == nil {
		return false
	}
	deadline := passwordChangedAt.AddDate(0, 0, cfg.MaxAgeDays)
	return gtime.Now().After(deadline)
}
