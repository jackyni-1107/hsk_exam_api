package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// Response 统一 REST 响应包装（与 GoFrame 默认 HandlerResponse 行为一致）。
func Response(r *ghttp.Request) {
	r.Middleware.Next()
}
