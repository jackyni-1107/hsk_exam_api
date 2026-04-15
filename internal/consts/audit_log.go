package consts

// 登录审计日志类型（sys_login_log.log_type 等）
const (
	AuditLogTypeLoginSuccess = 1
	AuditLogTypeLoginFail    = 2
	AuditLogTypeLogout       = 3
)

const (
	EventTypeTokenInvalid     = "token_invalid"
	EventTypePermissionDenied = "permission_denied"
	EventTypeBruteForce       = "brute_force"
	EventTypeSuspiciousIP     = "suspicious_ip"
)
