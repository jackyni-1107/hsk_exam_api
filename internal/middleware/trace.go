package middleware

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"
)

type traceCtxKey string

const traceIdKey traceCtxKey = "trace_id"

// Trace 链路追踪中间件：生成或透传 TraceId
func Trace(r *ghttp.Request) {
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		traceId = uuid.New().String()
	}
	ctx := context.WithValue(r.GetCtx(), traceIdKey, traceId)
	r.SetCtx(ctx)
	r.Middleware.Next()
	r.Response.Header().Set("X-Trace-Id", traceId)
}

// GetTraceId 从 context 获取 TraceId
func GetTraceId(ctx context.Context) string {
	if v := ctx.Value(traceIdKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
