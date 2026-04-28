package consts

// 登录审计日志类型（sys_login_log.log_type 等）
const (
	AuditLogTypeLoginSuccess = "login_success"
	AuditLogTypeLoginFail    = "login_fail"
	AuditLogTypeLogout       = "logout"
)

const (
	EventTypeTokenInvalid     = "token_invalid"
	EventTypePermissionDenied = "permission_denied"
	EventTypeBruteForce       = "brute_force"
	EventTypeSuspiciousIP     = "suspicious_ip"
)
