package auth

import (
	"context"

	"exam/api/client/auth/v1"
)

type IAuth interface {
	ILogin
	ILogout
	ICaptcha
}

type ILogin interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
}

type ILogout interface {
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
}

type ICaptcha interface {
	Captcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error)
}
