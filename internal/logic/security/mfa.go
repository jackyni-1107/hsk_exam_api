package security

import (
	"context"
)

// MFANotImplemented 预留：等保/高安全场景可接入 TOTP、WebAuthn 等。
// 配置 security.mfa.enabled 为 true 时，业务侧应在此处实现校验并在登录流程中串联。
func MFANotImplemented(ctx context.Context) bool {
	return !LoadMFACfg(ctx).Enabled
}
