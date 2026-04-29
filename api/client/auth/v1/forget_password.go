package v1

import "github.com/gogf/gf/v2/frame/g"

type ForgetPasswordReq struct {
	g.Meta   `path:"/auth/forget-password" method:"post" tags:"客户端认证" summary:"忘记密码"`
	Username string `json:"username" v:"required#err.invalid_params" dc:"用户名"`
}

type ForgetPasswordRes struct{}
