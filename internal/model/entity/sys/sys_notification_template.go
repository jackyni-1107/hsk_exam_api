// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotificationTemplate is the golang structure for table sys_notification_template.
type SysNotificationTemplate struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`                  // 主键ID
	Code       string      `json:"code"        orm:"code"        description:"模板编码"`                  // 模板编码
	Name       string      `json:"name"        orm:"name"        description:"模板名称"`                  // 模板名称
	Channel    string      `json:"channel"     orm:"channel"     description:"渠道：sms/email/template"` // 渠道：sms/email/template
	Content    string      `json:"content"     orm:"content"     description:"模板内容，支持变量 {{var}}"`     // 模板内容，支持变量 {{var}}
	Variables  string      `json:"variables"   orm:"variables"   description:"变量列表，逗号分隔"`             // 变量列表，逗号分隔
	Status     int         `json:"status"      orm:"status"      description:"状态：0-启用，1-停用"`          // 状态：0-启用，1-停用
	Remark     string      `json:"remark"      orm:"remark"      description:"备注"`                    // 备注
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`                   // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`                  // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`                   // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`                  // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除"`                  // 逻辑删除
}
