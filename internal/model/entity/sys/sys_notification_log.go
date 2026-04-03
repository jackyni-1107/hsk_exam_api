// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotificationLog is the golang structure for table sys_notification_log.
type SysNotificationLog struct {
	Id           int64       `json:"id"            orm:"id"            description:"主键ID"`                  // 主键ID
	TemplateCode string      `json:"template_code" orm:"template_code" description:"模板编码"`                  // 模板编码
	Channel      string      `json:"channel"       orm:"channel"       description:"渠道：sms/email/template"` // 渠道：sms/email/template
	Recipient    string      `json:"recipient"     orm:"recipient"     description:"接收者（手机号/邮箱/用户ID）"`      // 接收者（手机号/邮箱/用户ID）
	Content      string      `json:"content"       orm:"content"       description:"发送内容"`                  // 发送内容
	Status       int         `json:"status"        orm:"status"        description:"状态：0-待发送，1-成功，2-失败"`    // 状态：0-待发送，1-成功，2-失败
	ErrorMsg     string      `json:"error_msg"     orm:"error_msg"     description:"失败原因"`                  // 失败原因
	CreateTime   *gtime.Time `json:"create_time"   orm:"create_time"   description:"创建时间"`                  // 创建时间
}
