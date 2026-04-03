package v1

import "github.com/gogf/gf/v2/frame/g"

type ConfigListReq struct {
	g.Meta `path:"/config/list" method:"get" tags:"系统配置" summary:"配置列表"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Group  string `json:"group"`
	Key    string `json:"key"`
}

type ConfigListRes struct {
	List  []*ConfigItem `json:"list"`
	Total int           `json:"total"`
}

type ConfigItem struct {
	Id          int64  `json:"id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	GroupName   string `json:"group_name"`
	Remark      string `json:"remark"`
	CreateTime  string `json:"create_time"`
}

type ConfigCreateReq struct {
	g.Meta      `path:"/config" method:"post" tags:"系统配置" summary:"新增配置"`
	ConfigKey   string `json:"config_key" v:"required#err.invalid_params"`
	ConfigValue string `json:"config_value" v:"required#err.invalid_params"`
	ConfigType  string `json:"config_type"`
	GroupName   string `json:"group_name"`
	Remark      string `json:"remark"`
}

type ConfigCreateRes struct {
	Id int64 `json:"id"`
}

type ConfigUpdateReq struct {
	g.Meta      `path:"/config/{id}" method:"put" tags:"系统配置" summary:"更新配置"`
	Id          int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	ConfigValue string `json:"config_value"`
	Remark      string `json:"remark"`
}

type ConfigUpdateRes struct{}

type ConfigDeleteReq struct {
	g.Meta `path:"/config/{id}" method:"delete" tags:"系统配置" summary:"删除配置"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type ConfigDeleteRes struct{}

type ConfigGetReq struct {
	g.Meta `path:"/config/get" method:"get" tags:"系统配置" summary:"按键取值"`
	Key    string `json:"key" v:"required#err.invalid_params"`
}

type ConfigGetRes struct {
	Value string `json:"value"`
}

type DictTypeListReq struct {
	g.Meta   `path:"/dict/type/list" method:"get" tags:"字典" summary:"字典类型列表"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	DictType string `json:"dict_type"`
}

type DictTypeListRes struct {
	List  []*DictTypeItem `json:"list"`
	Total int             `json:"total"`
}

type DictTypeItem struct {
	Id         int64  `json:"id"`
	DictName   string `json:"dict_name"`
	DictType   string `json:"dict_type"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

type DictTypeCreateReq struct {
	g.Meta   `path:"/dict/type" method:"post" tags:"字典" summary:"新增字典类型"`
	DictName string `json:"dict_name" v:"required#err.invalid_params"`
	DictType string `json:"dict_type" v:"required#err.invalid_params"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

type DictTypeCreateRes struct {
	Id int64 `json:"id"`
}

type DictTypeUpdateReq struct {
	g.Meta   `path:"/dict/type/{id}" method:"put" tags:"字典" summary:"更新字典类型"`
	Id       int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	DictName string `json:"dict_name"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

type DictTypeUpdateRes struct{}

type DictTypeDeleteReq struct {
	g.Meta `path:"/dict/type/{id}" method:"delete" tags:"字典" summary:"删除字典类型"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type DictTypeDeleteRes struct{}

type DictDataListReq struct {
	g.Meta   `path:"/dict/data/list" method:"get" tags:"字典" summary:"字典数据列表"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	DictType string `json:"dict_type" v:"required#err.invalid_params"`
}

type DictDataListRes struct {
	List  []*DictDataItem `json:"list"`
	Total int             `json:"total"`
}

type DictDataItem struct {
	Id         int64  `json:"id"`
	DictType   string `json:"dict_type"`
	DictLabel  string `json:"dict_label"`
	DictValue  string `json:"dict_value"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

type DictDataCreateReq struct {
	g.Meta    `path:"/dict/data" method:"post" tags:"字典" summary:"新增字典数据"`
	DictType  string `json:"dict_type" v:"required#err.invalid_params"`
	DictLabel string `json:"dict_label" v:"required#err.invalid_params"`
	DictValue string `json:"dict_value" v:"required#err.invalid_params"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
}

type DictDataCreateRes struct {
	Id int64 `json:"id"`
}

type DictDataUpdateReq struct {
	g.Meta    `path:"/dict/data/{id}" method:"put" tags:"字典" summary:"更新字典数据"`
	Id        int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	DictLabel string `json:"dict_label"`
	DictValue string `json:"dict_value"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
}

type DictDataUpdateRes struct{}

type DictDataDeleteReq struct {
	g.Meta `path:"/dict/data/{id}" method:"delete" tags:"字典" summary:"删除字典数据"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type DictDataDeleteRes struct{}
