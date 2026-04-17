package security

import (
	"context"

	appcfg "exam/internal/config"
	"exam/internal/consts"
	"exam/internal/model/bo"
)

func (s *sSecurity) LoadLoginCfg(ctx context.Context) bo.LoginCfg {
	return appcfg.Config.Login
}

func (s *sSecurity) LoadPasswordCfg(ctx context.Context) bo.PasswordCfg {
	return appcfg.Config.Password
}

func (s *sSecurity) LoadSessionCfg(ctx context.Context) bo.SessionCfg {
	return appcfg.Config.Session
}

func (s *sSecurity) LoadSM2Cfg(ctx context.Context) bo.SM2Cfg {
	return appcfg.Config.SM2
}

// TokenTTLSeconds 会话 Token 有效期（秒），供服务接口实现使用。
func (s *sSecurity) TokenTTLSeconds(ctx context.Context) int64 {
	c := s.LoadSessionCfg(ctx)
	if c.TokenTTLSeconds <= 0 {
		return consts.DefaultTokenTTLFallbackSeconds
	}
	return c.TokenTTLSeconds
}
