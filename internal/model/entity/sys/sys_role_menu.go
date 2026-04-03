// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRoleMenu is the golang structure for table sys_role_menu.
type SysRoleMenu struct {
	Id         int64       `json:"id"          orm:"id"          description:"自增编号"`               // 自增编号
	RoleId     int64       `json:"role_id"     orm:"role_id"     description:"角色ID"`               // 角色ID
	MenuId     int64       `json:"menu_id"     orm:"menu_id"     description:"菜单ID"`               // 菜单ID
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`                // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`               // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`                // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`               // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除标识：0-未删除，1-已删除"` // 逻辑删除标识：0-未删除，1-已删除
}
