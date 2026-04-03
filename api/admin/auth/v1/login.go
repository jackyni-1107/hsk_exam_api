package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta        `path:"/auth/login" method:"post" tags:"管理端认证" summary:"管理端登录"`
	Username      string `json:"username" v:"required#err.invalid_params"`
	Password      string `json:"password" v:"required#err.invalid_params"`
	CaptchaId     string `json:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer"`
}

type LoginRes struct {
	Token    string     `json:"token"`
	UserInfo *LoginUser `json:"user_info"`
}

type LoginUser struct {
	Id       int64    `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles,omitempty"`
}
