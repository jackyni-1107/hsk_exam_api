package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta        `path:"/auth/login" method:"post" tags:"客户端认证" summary:"客户端登录"`
	Username      string `json:"username" v:"required#err.invalid_params" dc:"用户名"`
	Password      string `json:"password" v:"required#err.invalid_params" dc:"密码"`
	CaptchaId     string `json:"captcha_id" dc:"验证码 ID（风控开启且失败次数达阈值时必填）"`
	CaptchaAnswer string `json:"captcha_answer" dc:"验证码答案"`
}

type LoginRes struct {
	Token    string     `json:"token" dc:"访问令牌"`
	UserInfo *LoginUser `json:"user_info" dc:"用户信息"`
}

type LoginUser struct {
	Id       int64    `json:"id" dc:"用户ID"`
	Username string   `json:"username" dc:"用户名"`
	Nickname string   `json:"nickname" dc:"昵称"`
	Avatar   string   `json:"avatar" dc:"头像"`
	Roles    []string `json:"roles,omitempty" dc:"角色列表"`
}
