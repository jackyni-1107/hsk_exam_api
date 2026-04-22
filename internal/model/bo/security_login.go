package bo

type LoginInput struct {
	UserType          int
	Username          string
	EncryptedPassword string
	CaptchaId         string
	CaptchaAnswer     string
	IP                string
	UserAgent         string
	TraceId           string
}

type LoginUserInfo struct {
	Id       int64
	Username string
	Nickname string
	Avatar   string
}

type LoginResult struct {
	Token    string
	UserInfo LoginUserInfo
}
