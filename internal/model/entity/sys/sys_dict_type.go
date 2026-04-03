// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictType is the golang structure for table sys_dict_type.
type SysDictType struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`         // 主键ID
	DictName   string      `json:"dict_name"   orm:"dict_name"   description:"字典名称"`         // 字典名称
	DictType   string      `json:"dict_type"   orm:"dict_type"   description:"字典类型（唯一标识）"`   // 字典类型（唯一标识）
	Status     int         `json:"status"      orm:"status"      description:"状态：0-正常，1-停用"` // 状态：0-正常，1-停用
	Remark     string      `json:"remark"      orm:"remark"      description:"备注"`           // 备注
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`          // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`         // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`          // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`         // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除"`         // 逻辑删除
}
