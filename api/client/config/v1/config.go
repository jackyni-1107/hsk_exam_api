package v1

import "github.com/gogf/gf/v2/frame/g"

type MultiGetReq struct {
	g.Meta `path:"/config" method:"get" tags:"客户端-系统配置" summary:"获取配置"`
	Keys   []string `json:"keys" in:"query" v:"required|length:1,50#err.invalid_params|err.invalid_params" dc:"配置键列表"`
}

type MultiGetRes struct {
	Items []*ConfigKV `json:"items" dc:"配置项列表"`
}

type ConfigKV struct {
	Key   string `json:"key" dc:"配置键"`
	Value string `json:"value" dc:"配置值"`
}
