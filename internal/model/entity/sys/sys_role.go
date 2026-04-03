// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRole is the golang structure for table sys_role.
type SysRole struct {
	Id         int64       `json:"id"          orm:"id"          description:"角色ID"`               // 角色ID
	Name       string      `json:"name"        orm:"name"        description:"角色名称"`               // 角色名称
	Code       string      `json:"code"        orm:"code"        description:"角色权限字符串"`            // 角色权限字符串
	Sort       int         `json:"sort"        orm:"sort"        description:"显示顺序"`               // 显示顺序
	Status     int         `json:"status"      orm:"status"      description:"角色状态（0正常 1停用）"`      // 角色状态（0正常 1停用）
	Type       int         `json:"type"        orm:"type"        description:"角色类型"`               // 角色类型
	Remark     string      `json:"remark"      orm:"remark"      description:"备注"`                 // 备注
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`                // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`               // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`                // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`               // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除标识：0-未删除，1-已删除"` // 逻辑删除标识：0-未删除，1-已删除
}
