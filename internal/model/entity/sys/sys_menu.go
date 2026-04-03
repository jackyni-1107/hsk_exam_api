// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMenu is the golang structure for table sys_menu.
type SysMenu struct {
	Id            int64       `json:"id"             orm:"id"             description:"菜单ID"`               // 菜单ID
	Name          string      `json:"name"           orm:"name"           description:"菜单名称"`               // 菜单名称
	Permission    string      `json:"permission"     orm:"permission"     description:"权限标识"`               // 权限标识
	Type          int         `json:"type"           orm:"type"           description:"菜单类型"`               // 菜单类型
	Sort          int         `json:"sort"           orm:"sort"           description:"显示顺序"`               // 显示顺序
	ParentId      int64       `json:"parent_id"      orm:"parent_id"      description:"父菜单ID"`              // 父菜单ID
	Path          string      `json:"path"           orm:"path"           description:"路由地址"`               // 路由地址
	Icon          string      `json:"icon"           orm:"icon"           description:"菜单图标"`               // 菜单图标
	Component     string      `json:"component"      orm:"component"      description:"组件路径"`               // 组件路径
	ComponentName string      `json:"component_name" orm:"component_name" description:"组件名"`                // 组件名
	Status        int         `json:"status"         orm:"status"         description:"菜单状态"`               // 菜单状态
	Visible       bool        `json:"visible"        orm:"visible"        description:"是否可见"`               // 是否可见
	KeepAlive     bool        `json:"keep_alive"     orm:"keep_alive"     description:"是否缓存"`               // 是否缓存
	AlwaysShow    bool        `json:"always_show"    orm:"always_show"    description:"是否总是显示"`             // 是否总是显示
	Creator       string      `json:"creator"        orm:"creator"        description:"创建者"`                // 创建者
	CreateTime    *gtime.Time `json:"create_time"    orm:"create_time"    description:"创建时间"`               // 创建时间
	Updater       string      `json:"updater"        orm:"updater"        description:"更新者"`                // 更新者
	UpdateTime    *gtime.Time `json:"update_time"    orm:"update_time"    description:"更新时间"`               // 更新时间
	DeleteFlag    int         `json:"delete_flag"    orm:"delete_flag"    description:"逻辑删除标识：0-未删除，1-已删除"` // 逻辑删除标识：0-未删除，1-已删除
}
