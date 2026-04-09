package v1

import "github.com/gogf/gf/v2/frame/g"

type CaptchaReq struct {
	g.Meta `path:"/auth/captcha" method:"get" tags:"客户端认证" summary:"获取登录验证码"`
}

type CaptchaRes struct {
	CaptchaId     string `json:"captcha_id" dc:"验证码 ID"`
	Question      string `json:"question" dc:"算术题（文本）"`
	QuestionImage string `json:"question_image" dc:"算术题 PNG 的 Base64，前端可用 data:image/png;base64,{question_image} 展示"`
}
