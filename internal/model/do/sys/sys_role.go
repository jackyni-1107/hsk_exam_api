// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRole is the golang structure of table sys_role for DAO operations like Where/Data.
type SysRole struct {
	g.Meta     `orm:"table:sys_role, do:true"`
	Id         any         // 角色ID
	Name       any         // 角色名称
	Code       any         // 角色权限字符串
	Sort       any         // 显示顺序
	Status     any         // 角色状态（0正常 1停用）
	Type       any         // 角色类型
	Remark     any         // 备注
	Creator    any         // 创建者
	CreateTime *gtime.Time // 创建时间
	Updater    any         // 更新者
	UpdateTime *gtime.Time // 更新时间
	DeleteFlag any         // 逻辑删除标识：0-未删除，1-已删除
}
