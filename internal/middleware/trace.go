package middleware

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// Trace 透传 / 生成 Trace-Id（依赖 GoFrame 上下文）。
func Trace(r *ghttp.Request) {
	r.Middleware.Next()
}

// GetTraceId 从上下文取链路 ID。
func GetTraceId(ctx context.Context) string {
	return gctx.CtxId(ctx)
}
