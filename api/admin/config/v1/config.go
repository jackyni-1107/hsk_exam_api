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
	g.Meta `path:"/config/get" method:"get" tags:"系统配置" summary:"按键取值"`
	Key    string `json:"key" v:"required#err.invalid_params" dc:"配置键"`
}

type ConfigGetRes struct {
	Value string `json:"value" dc:"配置值"`
}

type DictTypeListReq struct {
	g.Meta   `path:"/dict/type/list" method:"get" tags:"字典" summary:"字典类型列表"`
	Page     int    `json:"page" dc:"页码"`
	Size     int    `json:"size" dc:"每页条数"`
	DictType string `json:"dict_type" dc:"字典类型"`
}

type DictTypeListRes struct {
	List  []*DictTypeItem `json:"list" dc:"列表"`
	Total int             `json:"total" dc:"总数"`
}

type DictTypeItem struct {
	Id         int64  `json:"id" dc:"字典类型ID"`
	DictName   string `json:"dict_name" dc:"字典名称"`
	DictType   string `json:"dict_type" dc:"字典类型"`
	Status     int    `json:"status" dc:"状态：0正常 1停用"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type DictTypeCreateReq struct {
	g.Meta   `path:"/dict/type" method:"post" tags:"字典" summary:"新增字典类型"`
	DictName string `json:"dict_name" v:"required#err.invalid_params" dc:"字典名称"`
	DictType string `json:"dict_type" v:"required#err.invalid_params" dc:"字典类型"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
	Remark   string `json:"remark" dc:"备注"`
}

type DictTypeCreateRes struct {
	Id int64 `json:"id" dc:"字典类型ID"`
}

type DictTypeUpdateReq struct {
	g.Meta   `path:"/dict/type/{id}" method:"put" tags:"字典" summary:"更新字典类型"`
	Id       int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典类型ID"`
	DictName string `json:"dict_name" dc:"字典名称"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
	Remark   string `json:"remark" dc:"备注"`
}

type DictTypeUpdateRes struct{}

type DictTypeDeleteReq struct {
	g.Meta `path:"/dict/type/{id}" method:"delete" tags:"字典" summary:"删除字典类型"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典类型ID"`
}

type DictTypeDeleteRes struct{}

type DictDataListReq struct {
	g.Meta   `path:"/dict/data/list" method:"get" tags:"字典" summary:"字典数据列表"`
	Page     int    `json:"page" dc:"页码"`
	Size     int    `json:"size" dc:"每页条数"`
	DictType string `json:"dict_type" v:"required#err.invalid_params" dc:"字典类型"`
}

type DictDataListRes struct {
	List  []*DictDataItem `json:"list" dc:"列表"`
	Total int             `json:"total" dc:"总数"`
}

type DictDataItem struct {
	Id         int64  `json:"id" dc:"字典数据ID"`
	DictType   string `json:"dict_type" dc:"字典类型"`
	DictLabel  string `json:"dict_label" dc:"字典标签"`
	DictValue  string `json:"dict_value" dc:"字典值"`
	Sort       int    `json:"sort" dc:"排序值"`
	Status     int    `json:"status" dc:"状态：0正常 1停用"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type DictDataCreateReq struct {
	g.Meta    `path:"/dict/data" method:"post" tags:"字典" summary:"新增字典数据"`
	DictType  string `json:"dict_type" v:"required#err.invalid_params" dc:"字典类型"`
	DictLabel string `json:"dict_label" v:"required#err.invalid_params" dc:"字典标签"`
	DictValue string `json:"dict_value" v:"required#err.invalid_params" dc:"字典值"`
	Sort      int    `json:"sort" dc:"排序值"`
	Status    int    `json:"status" dc:"状态：0正常 1停用"`
	Remark    string `json:"remark" dc:"备注"`
}

type DictDataCreateRes struct {
	Id int64 `json:"id" dc:"字典数据ID"`
}

type DictDataUpdateReq struct {
	g.Meta    `path:"/dict/data/{id}" method:"put" tags:"字典" summary:"更新字典数据"`
	Id        int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典数据ID"`
	DictLabel string `json:"dict_label" dc:"字典标签"`
	DictValue string `json:"dict_value" dc:"字典值"`
	Sort      int    `json:"sort" dc:"排序值"`
	Status    int    `json:"status" dc:"状态：0正常 1停用"`
	Remark    string `json:"remark" dc:"备注"`
}

type DictDataUpdateRes struct{}

type DictDataDeleteReq struct {
	g.Meta `path:"/dict/data/{id}" method:"delete" tags:"字典" summary:"删除字典数据"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典数据ID"`
}

type DictDataDeleteRes struct{}
