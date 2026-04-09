package openapi

import "github.com/gogf/gf/v2/net/ghttp"

// RegisterSplitEndpoints 注册按端拆分的 OpenAPI JSON：
//   - GET /openapi/admin.json  管理端
//   - GET /openapi/client.json 学员端
//
// apiPrefix 与 server.apiPrefix 一致（如 "/api" 或 ""），用于按实际挂载路径筛选 Paths。
// 完整规范仍为配置项 server.openapiPath（默认 /api.json）。
func RegisterSplitEndpoints(s *ghttp.Server, apiPrefix string) {
	registerSplitOpenAPIHandlers(s, apiPrefix)
}
