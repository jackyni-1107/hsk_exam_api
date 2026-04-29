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

// ---------- 忘记密码 Redis 键 ----------

const (
	// ForgetPasswordCooldownUserKeyFmt 单账号冷却键（%s=标准化后的账号），用于 60 秒内防重复提交。
	ForgetPasswordCooldownUserKeyFmt = "hskexam:security:forget_password:cooldown:user:%s"

	// ForgetPasswordCooldownTTLSeconds 找回密码成功触发后的冷却时长（秒）。
	ForgetPasswordCooldownTTLSeconds = 60
)

const (
	// ForgetPasswordDailyIPCounterKeyFmt 每日 IP 访问计数键（第1个%s=IP，第2个%s=日期yyyyMMdd）。
	ForgetPasswordDailyIPCounterKeyFmt = "hskexam:security:forget_password:daily_counter:ip:%s:%s"
	// ForgetPasswordDailyUserCounterKeyFmt 每日账号访问计数键（第1个%s=账号，第2个%s=日期yyyyMMdd）。
	ForgetPasswordDailyUserCounterKeyFmt = "hskexam:security:forget_password:daily_counter:user:%s:%s"
	// ForgetPasswordBlockedIPKeyFmt IP 封禁键（%s=IP），命中后直接拒绝访问找回密码接口。
	ForgetPasswordBlockedIPKeyFmt = "hskexam:security:forget_password:blocked:ip:%s"
	// ForgetPasswordBlockedUserKeyFmt 账号封禁键（%s=账号），命中后直接拒绝访问找回密码接口。
	ForgetPasswordBlockedUserKeyFmt = "hskexam:security:forget_password:blocked:user:%s"

	// ForgetPasswordDailyAccessLimit 每个账号或 IP 每日最大访问次数（超过该值即触发 24h 封禁）。
	ForgetPasswordDailyAccessLimit = 10
	// ForgetPasswordBlockedTTLSeconds 超限后的封禁时长（秒）。当前=24小时。
	ForgetPasswordBlockedTTLSeconds = 24 * 3600
	// ForgetPasswordDailyCounterMaxTTL 每日计数键保留时长（秒），用于覆盖跨日并防止键长期堆积。
	ForgetPasswordDailyCounterMaxTTL = 48 * 3600
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
