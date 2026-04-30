package bo

// MemberImportResult 客户批量导入统计
type MemberImportResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

// MemberImportOptions 客户批量导入附加配置
type MemberImportOptions struct {
	// UseRandomPassword 为 true 时，密码为空行使用随机密码。
	// 为 false 时，密码为空行按邮箱规则生成（第1/3/5位 + "@" + FixedPasswordSuffix）。
	UseRandomPassword bool `json:"use_random_password"`
	// EmailPasswordPickPositions 固定规则取位，1-based，逗号分隔，例如 "1,3,5"。
	EmailPasswordPickPositions string `json:"email_password_pick_positions"`
	// FixedPasswordSuffix 固定规则后缀，默认 "hskmock"。
	FixedPasswordSuffix string `json:"fixed_password_suffix"`
	// SendPasswordNotice 导入成功后是否发送账号密码通知（模板 forget_password，渠道 email）。
	SendPasswordNotice bool `json:"send_password_notice"`
}
