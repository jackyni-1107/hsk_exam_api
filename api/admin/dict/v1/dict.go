package v1

import "github.com/gogf/gf/v2/frame/g"

type DictTypeListReq struct {
	g.Meta   `path:"/dict/type/list" method:"get" tags:"字典" summary:"字典类型列表" permission:"dict:type:list"`
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
	Status     int    `json:"status" dc:"状态"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type DictTypeCreateReq struct {
	g.Meta   `path:"/dict/type" method:"post" tags:"字典" summary:"新增字典类型" permission:"dict:type:create"`
	DictName string `json:"dict_name" v:"required#err.invalid_params" dc:"字典名称"`
	DictType string `json:"dict_type" v:"required#err.invalid_params" dc:"字典类型"`
	Status   int    `json:"status" dc:"状态"`
	Remark   string `json:"remark" dc:"备注"`
}

type DictTypeCreateRes struct {
	Id int64 `json:"id" dc:"字典类型ID"`
}

type DictTypeUpdateReq struct {
	g.Meta   `path:"/dict/type/{id}" method:"put" tags:"字典" summary:"更新字典类型" permission:"dict:type:update"`
	Id       int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典类型ID"`
	DictName string `json:"dict_name" dc:"字典名称"`
	Status   int    `json:"status" dc:"状态"`
	Remark   string `json:"remark" dc:"备注"`
}

type DictTypeUpdateRes struct{}

type DictTypeDeleteReq struct {
	g.Meta `path:"/dict/type/{id}" method:"delete" tags:"字典" summary:"删除字典类型" permission:"dict:type:delete"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典类型ID"`
}

type DictTypeDeleteRes struct{}

type DictDataListReq struct {
	g.Meta   `path:"/dict/data/list" method:"get" tags:"字典" summary:"字典数据列表" permission:"dict:data:list"`
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
	Status     int    `json:"status" dc:"状态"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type DictDataCreateReq struct {
	g.Meta    `path:"/dict/data" method:"post" tags:"字典" summary:"新增字典数据" permission:"dict:data:create"`
	DictType  string `json:"dict_type" v:"required#err.invalid_params" dc:"字典类型"`
	DictLabel string `json:"dict_label" v:"required#err.invalid_params" dc:"字典标签"`
	DictValue string `json:"dict_value" v:"required#err.invalid_params" dc:"字典值"`
	Sort      int    `json:"sort" dc:"排序值"`
	Status    int    `json:"status" dc:"状态"`
	Remark    string `json:"remark" dc:"备注"`
}

type DictDataCreateRes struct {
	Id int64 `json:"id" dc:"字典数据ID"`
}

type DictDataUpdateReq struct {
	g.Meta    `path:"/dict/data/{id}" method:"put" tags:"字典" summary:"更新字典数据" permission:"dict:data:update"`
	Id        int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典数据ID"`
	DictLabel string `json:"dict_label" dc:"字典标签"`
	DictValue string `json:"dict_value" dc:"字典值"`
	Sort      int    `json:"sort" dc:"排序值"`
	Status    int    `json:"status" dc:"状态"`
	Remark    string `json:"remark" dc:"备注"`
}

type DictDataUpdateRes struct{}

type DictDataDeleteReq struct {
	g.Meta `path:"/dict/data/{id}" method:"delete" tags:"字典" summary:"删除字典数据" permission:"dict:data:delete"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"字典数据ID"`
}

type DictDataDeleteRes struct{}
