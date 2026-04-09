package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
)

// HandlerResponseI18n 带国际化的响应中间件，替代 ghttp.MiddlewareHandlerResponse
// 根据 Language-Accept 或 Accept-Language 翻译错误 message；支持多语言资源缺失时回退
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
		raw := errorMessageKey(err)
		lang := parseLanguage(r)
		ctx := r.GetCtx()
		msg = translateWithFallback(ctx, lang, raw)
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

// errorMessageKey 取用于展示/翻译的文案：优先链路上显式 gerror.Text（可与 Code 的 message 不同），否则 Code.Message，否则 Error()
func errorMessageKey(err error) string {
	if err == nil {
		return ""
	}
	for e := err; e != nil; e = errors.Unwrap(e) {
		var ge *gerror.Error
		if !errors.As(e, &ge) {
			continue
		}
		if t := strings.TrimSpace(ge.Text()); t != "" {
			return ge.TextWithArgs()
		}
	}
	ec := gerror.Code(err)
	if ec != nil && ec != gcode.CodeNil {
		if m := strings.TrimSpace(ec.Message()); m != "" {
			return m
		}
	}
	return err.Error()
}

// translateWithFallback 按请求语言翻译；若无词条或该语言未配置，依次回退 zh-CN、en
func translateWithFallback(ctx context.Context, lang, key string) string {
	if key == "" {
		return ""
	}
	seen := make(map[string]struct{})
	try := func(l string) (string, bool) {
		l = strings.TrimSpace(l)
		if l == "" {
			return "", false
		}
		if _, ok := seen[l]; ok {
			return "", false
		}
		seen[l] = struct{}{}
		cctx := gi18n.WithLanguage(ctx, l)
		out := g.I18n().T(cctx, key)
		if out != key {
			return out, true
		}
		if gc := g.I18n().GetContent(cctx, key); gc != "" {
			return gc, true
		}
		return "", false
	}
	if s, ok := try(lang); ok {
		return s
	}
	if s, ok := try("zh-CN"); ok {
		return s
	}
	if s, ok := try("en"); ok {
		return s
	}
	return key
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

// normalizeLang 将常见 BCP 47 标签映射到 manifest/i18n 下的语言名
func normalizeLang(lang string) string {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return "zh-CN"
	}
	low := strings.ToLower(lang)
	switch {
	case strings.HasPrefix(low, "zh-hant"), strings.HasPrefix(low, "zh-tw"), strings.HasPrefix(low, "zh-hk"):
		return "zh-TW"
	case strings.HasPrefix(low, "zh"):
		return "zh-CN"
	case strings.HasPrefix(low, "en"):
		return "en"
	case strings.HasPrefix(low, "ja"):
		return "ja"
	default:
		if i := strings.Index(low, "-"); i > 0 {
			return lang[:i] + strings.ToUpper(lang[i:])
		}
		return low
	}
}
