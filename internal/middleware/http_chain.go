package middleware

import (
	"exam/internal/consts"

	"github.com/gogf/gf/v2/net/ghttp"
)

func ClientPublicChain() []ghttp.HandlerFunc {
	return []ghttp.HandlerFunc{
		Trace,
		Response,
		HandlerResponseI18n,
	}
}

func ClientProtectedChain() []ghttp.HandlerFunc {
	return append(ClientPublicChain(), Auth(consts.UserTypeClient))
}

func ClientMediaChain() []ghttp.HandlerFunc {
	return []ghttp.HandlerFunc{
		Trace,
		MiddlewareCORS,
		Response,
	}
}

func AdminPublicChain() []ghttp.HandlerFunc {
	return []ghttp.HandlerFunc{
		Trace,
		Response,
		HandlerResponseI18n,
	}
}

func AdminProtectedChain() []ghttp.HandlerFunc {
	return append(AdminPublicChain(),
		Auth(consts.UserTypeAdmin),
		RBACFromPath,
		Audit,
	)
}

func AdminDownloadChain() []ghttp.HandlerFunc {
	return []ghttp.HandlerFunc{
		Trace,
		Response,
		Auth(consts.UserTypeAdmin),
		RBACFromPath,
		Audit,
	}
}
