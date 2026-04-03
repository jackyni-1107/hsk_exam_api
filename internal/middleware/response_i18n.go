package middleware

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// HandlerResponseI18n 带国际化的响应中间件，替代 ghttp.MiddlewareHandlerResponse
// 根据 Language-Accept 或 Accept-Language 翻译错误 message
func HandlerResponseI18n(r *ghttp.Request) {
	r.Middleware.Next()

	var (
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code int
		msg  string
	)

	if err != nil {
		// GoFrame 在 middleware.TryCatch 恢复 panic 时会先 WriteStatus(500, exception)，
		// 把「exception recovered: …」写入 buffer；若不清理会与下方 WriteJson 拼接，客户端看到非纯 JSON。
		if r.Response.BufferLength() > 0 {
			r.Response.ClearBuffer()
			r.Response.WriteHeader(http.StatusOK)
		}
		ec := gerror.Code(err)
		code = ec.Code()
		i18nKey := ec.Message()
		if i18nKey == "" {
			i18nKey = err.Error()
		}
		// 仅对 err.xxx 格式的 key 做 i18n 翻译
		if strings.HasPrefix(i18nKey, "err.") {
			lang := parseLanguage(r)
			ctx := gi18n.WithLanguage(r.GetCtx(), lang)
			translated := g.I18n().T(ctx, i18nKey)
			if translated != "" && translated != i18nKey {
				msg = translated
			} else {
				msg = i18nKey
			}
		} else {
			msg = i18nKey
		}
	} else {
		code = 0
		msg = "OK"
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code,
		Message: msg,
		Data:    res,
	})
}

// parseLanguage 从请求头解析语言：优先 Language-Accept，其次 Accept-Language 第一个值，默认 zh-CN
func parseLanguage(r *ghttp.Request) string {
	if lang := r.Header.Get("Language-Accept"); lang != "" {
		return normalizeLang(lang)
	}
	if lang := r.Header.Get("Accept-Language"); lang != "" {
		// 取第一个语言标签，如 "zh-CN,zh;q=0.9,en;q=0.8" -> "zh-CN"
		first := strings.SplitN(strings.TrimSpace(lang), ",", 2)[0]
		first = strings.SplitN(strings.TrimSpace(first), ";", 2)[0]
		return normalizeLang(first)
	}
	return "zh-CN"
}

func normalizeLang(lang string) string {
	lang = strings.TrimSpace(lang)
	switch {
	case strings.HasPrefix(strings.ToLower(lang), "zh"):
		return "zh-CN"
	case strings.HasPrefix(strings.ToLower(lang), "en"):
		return "en"
	default:
		return lang
	}
}
