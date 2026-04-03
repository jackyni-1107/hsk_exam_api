package security

import (
	"context"

	appcfg "exam/internal/config"
	"exam/internal/model/bo"
)

func LoadLoginCfg(ctx context.Context) bo.LoginCfg {
	return appcfg.Config.Login
}

func LoadPasswordCfg(ctx context.Context) bo.PasswordCfg {
	return appcfg.Config.Password
}

func LoadSessionCfg(ctx context.Context) bo.SessionCfg {
	return appcfg.Config.Session
}

func LoadMFACfg(ctx context.Context) bo.MFACfg {
	return appcfg.Config.MFA
}

// TokenTTLSeconds 会话 Token 有效期（秒），供服务接口实现使用。
func TokenTTLSeconds(ctx context.Context) int64 {
	c := LoadSessionCfg(ctx)
	if c.TokenTTLSeconds <= 0 {
		return 86400
	}
	return c.TokenTTLSeconds
}
