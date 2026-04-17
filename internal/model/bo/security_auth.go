package bo

type CaptchaChallenge struct {
	CaptchaId string `json:"captcha_id"`
	Question  string `json:"question"`
}

type LoginCfg struct {
	RateLimitPerMinute    int  `json:"RateLimitPerMinute"`
	MaxFailuresBeforeLock int  `json:"MaxFailuresBeforeLock"`
	FailureWindowSeconds  int  `json:"FailureWindowSeconds"`
	LockDurationSeconds   int  `json:"LockDurationSeconds"`
	CaptchaEnabled        bool `json:"CaptchaEnabled"`
	CaptchaAfterFailures  int  `json:"CaptchaAfterFailures"`
}

type PasswordCfg struct {
	MinLength      int  `json:"MinLength"`
	RequireUpper   bool `json:"RequireUpper"`
	RequireLower   bool `json:"RequireLower"`
	RequireDigit   bool `json:"RequireDigit"`
	RequireSpecial bool `json:"RequireSpecial"`
	MaxAgeDays     int  `json:"MaxAgeDays"`
	HistoryCount   int  `json:"HistoryCount"`
}

type SessionCfg struct {
	TokenTTLSeconds       int64 `json:"TokenTTLSeconds"`
	MaxConcurrentSessions int   `json:"MaxConcurrentSessions"`
}

type SM2Cfg struct {
	// PrivateKeyPem SM2 私钥 PEM（建议通过环境变量注入，不入库）
	PrivateKeyPem string `json:"privateKeyPem"`
	// PublicKeyHex 可选：显式配置前端加密公钥（十六进制）；未配置时由私钥推导
	PublicKeyHex string `json:"publicKeyHex"`
}
