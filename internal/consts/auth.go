package consts

// ---------- Token ----------

const (
	// TokenRedisKeyPrefix Redis key 前缀（hskexam:security:token:{userTypeTag}:{token}）
	TokenRedisKeyPrefix = "hskexam:security:token:"
	// TokenTTL 默认 7 天（秒）
	TokenTTL = 7 * 24 * 3600
	// DefaultTokenTTLFallbackSeconds 配置缺失时的保底会话时长（1 天）
	DefaultTokenTTLFallbackSeconds int64 = 86400
)

// ---------- 用户类型 Tag（用于 Redis 键拼接等） ----------

const (
	UserTypeTagAdmin  = "admin"
	UserTypeTagClient = "client"
)

// ---------- Security Redis 键前缀（hskexam:security:业务名:） ----------

const (
	LoginRateLimitKeyPrefix = "hskexam:security:login_rate_limit:"
	LoginFailCountKeyPrefix = "hskexam:security:login_fail_count:"
	LoginLockKeyPrefix      = "hskexam:security:login_lock:"
	SessionListKeyPrefix    = "hskexam:security:session_list:"
	CaptchaKeyPrefix        = "hskexam:security:captcha:"
	CaptchaTTLSeconds       = 300
)

// ---------- Security 事件类型 ----------

const (
	SecurityEventBruteForce       = "brute_force"
	SecurityEventPermissionDenied = "permission_denied"
)

// ---------- RBAC 权限缓存 ----------

const (
	PermCacheKeyPrefix  = "hskexam:rbac:user_perms:"
	PermCacheTTLSeconds = 300
)

// ---------- 表名 ----------

const (
	TableSysPasswordHistory = "sys_password_history"
)
