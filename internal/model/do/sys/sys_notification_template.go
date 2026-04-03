// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotificationTemplate is the golang structure of table sys_notification_template for DAO operations like Where/Data.
type SysNotificationTemplate struct {
	g.Meta     `orm:"table:sys_notification_template, do:true"`
	Id         any         // 主键ID
	Code       any         // 模板编码
	Name       any         // 模板名称
	Channel    any         // 渠道：sms/email/template
	Content    any         // 模板内容，支持变量 {{var}}
	Variables  any         // 变量列表，逗号分隔
	Status     any         // 状态：0-启用，1-停用
	Remark     any         // 备注
	Creator    any         // 创建者
	CreateTime *gtime.Time // 创建时间
	Updater    any         // 更新者
	UpdateTime *gtime.Time // 更新时间
	DeleteFlag any         // 逻辑删除
}
