package auth

import (
	"context"

	"exam/api/client/auth/v1"
)

type IAuth interface {
	ILogin
	ILogout
	ICaptcha
	IPublicKey
	IForgetPassword
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

type IPublicKey interface {
	PublicKey(ctx context.Context, req *v1.PublicKeyReq) (res *v1.PublicKeyRes, err error)
}

type IForgetPassword interface {
	ForgetPassword(ctx context.Context, req *v1.ForgetPasswordReq) (res *v1.ForgetPasswordRes, err error)
}
