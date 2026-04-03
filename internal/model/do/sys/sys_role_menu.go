// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRoleMenu is the golang structure of table sys_role_menu for DAO operations like Where/Data.
type SysRoleMenu struct {
	g.Meta     `orm:"table:sys_role_menu, do:true"`
	Id         any         // 自增编号
	RoleId     any         // 角色ID
	MenuId     any         // 菜单ID
	Creator    any         // 创建者
	CreateTime *gtime.Time // 创建时间
	Updater    any         // 更新者
	UpdateTime *gtime.Time // 更新时间
	DeleteFlag any         // 逻辑删除标识：0-未删除，1-已删除
}
