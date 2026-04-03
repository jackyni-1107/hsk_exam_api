// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotificationLog is the golang structure of table sys_notification_log for DAO operations like Where/Data.
type SysNotificationLog struct {
	g.Meta       `orm:"table:sys_notification_log, do:true"`
	Id           any         // 主键ID
	TemplateCode any         // 模板编码
	Channel      any         // 渠道：sms/email/template
	Recipient    any         // 接收者（手机号/邮箱/用户ID）
	Content      any         // 发送内容
	Status       any         // 状态：0-待发送，1-成功，2-失败
	ErrorMsg     any         // 失败原因
	CreateTime   *gtime.Time // 创建时间
}
