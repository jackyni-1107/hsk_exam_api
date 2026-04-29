package v1

import "github.com/gogf/gf/v2/frame/g"

type ForgetPasswordReq struct {
	g.Meta `path:"/auth/forget-password" method:"post" tags:"客户端认证" summary:"忘记密码"`
	Email  string `json:"email" v:"required#err.invalid_params" dc:"邮箱"`
}

type ForgetPasswordRes struct{}
