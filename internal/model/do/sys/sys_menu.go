// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMenu is the golang structure of table sys_menu for DAO operations like Where/Data.
type SysMenu struct {
	g.Meta        `orm:"table:sys_menu, do:true"`
	Id            any         // 菜单ID
	Name          any         // 菜单名称
	Permission    any         // 权限标识
	Type          any         // 菜单类型
	Sort          any         // 显示顺序
	ParentId      any         // 父菜单ID
	Path          any         // 路由地址
	Icon          any         // 菜单图标
	Component     any         // 组件路径
	ComponentName any         // 组件名
	Status        any         // 菜单状态
	Visible       any         // 是否可见
	KeepAlive     any         // 是否缓存
	AlwaysShow    any         // 是否总是显示
	Creator       any         // 创建者
	CreateTime    *gtime.Time // 创建时间
	Updater       any         // 更新者
	UpdateTime    *gtime.Time // 更新时间
	DeleteFlag    any         // 逻辑删除标识：0-未删除，1-已删除
}
