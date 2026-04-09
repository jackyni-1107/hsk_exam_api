package openapi

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"

	appcfg "exam/internal/config"
)

const (
	pathOpenAPIAdmin  = "/openapi/admin.json"
	pathOpenAPIClient = "/openapi/client.json"
)

// filterOpenAPIByPathPrefix 从完整规范中按路径前缀筛出子集（浅拷贝 Paths，不修改服务端原始对象）。
func filterOpenAPIByPathPrefix(full *goai.OpenApiV3, prefix string, titleSuffix, descLine string) *goai.OpenApiV3 {
	if full == nil {
		return goai.New()
	}
	out := *full
	out.Paths = goai.Paths{}
	for p, item := range full.Paths {
		if strings.HasPrefix(p, prefix) {
			out.Paths[p] = item
		}
	}
	out.Info = full.Info
	out.Info.Title = strings.TrimSpace(full.Info.Title + titleSuffix)
	if descLine != "" {
		if full.Info.Description != "" {
			out.Info.Description = descLine + "\n\n" + full.Info.Description
		} else {
			out.Info.Description = descLine
		}
	}
	return &out
}

func writeOpenAPIJSON(r *ghttp.Request, doc *goai.OpenApiV3) {
	r.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.Response.WriteJson(doc)
}

// registerSplitOpenAPIHandlers 注册管理端 / 客户端两份 OpenAPI JSON（依赖 GoFrame 在 Run 阶段生成的完整规范）。
func registerSplitOpenAPIHandlers(s *ghttp.Server, apiPrefix string) {
	adminP := appcfg.JoinHTTPPath(apiPrefix, "admin")
	clientP := appcfg.JoinHTTPPath(apiPrefix, "client")
	s.BindHandler(pathOpenAPIAdmin, func(r *ghttp.Request) {
		desc := fmt.Sprintf("仅包含 `%s` 前缀的接口。", adminP)
		doc := filterOpenAPIByPathPrefix(s.GetOpenApi(), adminP, "（管理端）", desc)
		writeOpenAPIJSON(r, doc)
	})
	s.BindHandler(pathOpenAPIClient, func(r *ghttp.Request) {
		desc := fmt.Sprintf("仅包含 `%s` 前缀的接口。", clientP)
		doc := filterOpenAPIByPathPrefix(s.GetOpenApi(), clientP, "（客户端）", desc)
		writeOpenAPIJSON(r, doc)
	})
}
