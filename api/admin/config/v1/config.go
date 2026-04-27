package v1

import "github.com/gogf/gf/v2/frame/g"

type ConfigListReq struct {
	g.Meta `path:"/config/list" method:"get" tags:"系统配置" summary:"配置列表" permission:"config:list"`
	Page   int    `json:"page" dc:"页码"`
	Size   int    `json:"size" dc:"每页条数"`
	Group  string `json:"group" dc:"配置分组"`
	Key    string `json:"key" dc:"配置键"`
}

type ConfigListRes struct {
	List  []*ConfigItem `json:"list" dc:"列表"`
	Total int           `json:"total" dc:"总数"`
}

type ConfigItem struct {
	Id          int64  `json:"id" dc:"配置ID"`
	ConfigKey   string `json:"config_key" dc:"配置键"`
	ConfigValue string `json:"config_value" dc:"配置值"`
	ConfigType  string `json:"config_type" dc:"配置类型"`
	GroupName   string `json:"group_name" dc:"分组名称"`
	Remark      string `json:"remark" dc:"备注"`
	CreateTime  string `json:"create_time" dc:"创建时间"`
}

type ConfigCreateReq struct {
	g.Meta      `path:"/config" method:"post" tags:"系统配置" summary:"新增配置" permission:"config:create"`
	ConfigKey   string `json:"config_key" v:"required#err.invalid_params" dc:"配置键"`
	ConfigValue string `json:"config_value" v:"required#err.invalid_params" dc:"配置值"`
	ConfigType  string `json:"config_type" dc:"配置类型"`
	GroupName   string `json:"group_name" dc:"分组名称"`
	Remark      string `json:"remark" dc:"备注"`
}

type ConfigCreateRes struct {
	Id int64 `json:"id" dc:"配置ID"`
}

type ConfigUpdateReq struct {
	g.Meta      `path:"/config/{id}" method:"put" tags:"系统配置" summary:"更新配置" permission:"config:update"`
	Id          int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
	ConfigValue string `json:"config_value" dc:"配置值"`
	Remark      string `json:"remark" dc:"备注"`
}

type ConfigUpdateRes struct{}

type ConfigDeleteReq struct {
	g.Meta `path:"/config/{id}" method:"delete" tags:"系统配置" summary:"删除配置" permission:"config:delete"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
}

type ConfigDeleteRes struct{}

type ConfigGetReq struct {
	g.Meta `path:"/config/get" method:"get" tags:"系统配置" summary:"按键获取配置"`
	Key    string `json:"key" v:"required#err.invalid_params" dc:"配置键"`
}

type ConfigGetRes struct {
	Value string `json:"value" dc:"配置值"`
}
