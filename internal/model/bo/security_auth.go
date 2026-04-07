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

type MFACfg struct {
	Enabled bool `json:"enabled"`
}
