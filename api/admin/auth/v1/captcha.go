package v1

import "github.com/gogf/gf/v2/frame/g"

type CaptchaReq struct {
	g.Meta `path:"/auth/captcha" method:"get" tags:"管理端认证" summary:"获取登录验证码"`
}

type CaptchaRes struct {
	CaptchaId string `json:"captcha_id" dc:"验证码 ID"`
	Question  string `json:"question" dc:"算术题"`
}
