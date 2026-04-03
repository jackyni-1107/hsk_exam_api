package middleware

import (
	"context"
)

type ctxKey int

const ctxUserKey ctxKey = 1

// CtxData 认证后注入上下文的用户信息。
type CtxData struct {
	UserId   int64
	UserType int
	Username string
}

// SetCtxData 写入 context（供 Auth 中间件调用）。
func SetCtxData(ctx context.Context, d *CtxData) context.Context {
	return context.WithValue(ctx, ctxUserKey, d)
}

// GetCtxData 读取当前请求用户信息。
func GetCtxData(ctx context.Context) *CtxData {
	v := ctx.Value(ctxUserKey)
	if v == nil {
		return nil
	}
	d, _ := v.(*CtxData)
	return d
}
