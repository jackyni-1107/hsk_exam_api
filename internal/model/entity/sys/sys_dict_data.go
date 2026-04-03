// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictData is the golang structure for table sys_dict_data.
type SysDictData struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`         // 主键ID
	DictType   string      `json:"dict_type"   orm:"dict_type"   description:"字典类型"`         // 字典类型
	DictLabel  string      `json:"dict_label"  orm:"dict_label"  description:"字典标签"`         // 字典标签
	DictValue  string      `json:"dict_value"  orm:"dict_value"  description:"字典值"`          // 字典值
	Sort       int         `json:"sort"        orm:"sort"        description:"排序"`           // 排序
	Status     int         `json:"status"      orm:"status"      description:"状态：0-正常，1-停用"` // 状态：0-正常，1-停用
	Remark     string      `json:"remark"      orm:"remark"      description:"备注"`           // 备注
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`          // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`         // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`          // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`         // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除"`         // 逻辑删除
}
