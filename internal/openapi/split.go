package openapi

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
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
func registerSplitOpenAPIHandlers(s *ghttp.Server) {
	s.BindHandler(pathOpenAPIAdmin, func(r *ghttp.Request) {
		doc := filterOpenAPIByPathPrefix(s.GetOpenApi(), "/api/admin", "（管理端）", "仅包含 `/api/admin` 前缀的接口。")
		writeOpenAPIJSON(r, doc)
	})
	s.BindHandler(pathOpenAPIClient, func(r *ghttp.Request) {
		doc := filterOpenAPIByPathPrefix(s.GetOpenApi(), "/api/client", "（客户端）", "仅包含 `/api/client` 前缀的接口。")
		writeOpenAPIJSON(r, doc)
	})
}
