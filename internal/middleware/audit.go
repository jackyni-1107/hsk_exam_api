package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// Audit 操作审计占位：可在此写入审计日志。
func Audit(r *ghttp.Request) {
	r.Middleware.Next()
}
