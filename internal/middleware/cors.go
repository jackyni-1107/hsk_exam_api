package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// MiddlewareCORS is a middleware handler for CORS with default options.
func MiddlewareCORS(r *ghttp.Request) {
	//r.Response.CORSDefault()
	r.Middleware.Next()
}
