// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictData is the golang structure of table sys_dict_data for DAO operations like Where/Data.
type SysDictData struct {
	g.Meta     `orm:"table:sys_dict_data, do:true"`
	Id         any         // 主键ID
	DictType   any         // 字典类型
	DictLabel  any         // 字典标签
	DictValue  any         // 字典值
	Sort       any         // 排序
	Status     any         // 状态：0-正常，1-停用
	Remark     any         // 备注
	Creator    any         // 创建者
	CreateTime *gtime.Time // 创建时间
	Updater    any         // 更新者
	UpdateTime *gtime.Time // 更新时间
	DeleteFlag any         // 逻辑删除
}
