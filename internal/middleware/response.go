package middleware

import (
	"runtime/debug"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	auditsvc "exam/internal/service/audit"
)

// Response 统一 REST 响应包装（与 GoFrame 默认 HandlerResponse 行为一致）。
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	if err := r.GetError(); err != nil {
		g.Log().Errorf(r.GetCtx(), "request error: path=%s method=%s err=%s",
			r.URL.Path, r.Method, err.Error())
		ctxData := GetCtxData(r.GetCtx())
		userId := int64(0)
		if ctxData != nil {
			userId = ctxData.UserId
		}
		stack := string(debug.Stack())
		auditsvc.Audit().RecordException(r.GetCtx(), r.URL.Path, r.Method, err.Error(), stack, userId, r.GetClientIp(), GetTraceId(r.GetCtx()))
	}
}
