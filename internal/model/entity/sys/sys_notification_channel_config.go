// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotificationChannelConfig is the golang structure for table sys_notification_channel_config.
type SysNotificationChannelConfig struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`                              // 主键ID
	Channel    string      `json:"channel"     orm:"channel"     description:"渠道类型：email-邮件，sms-短信"`              // 渠道类型：email-邮件，sms-短信
	Provider   string      `json:"provider"    orm:"provider"    description:"提供商：email用smtp；sms用aliyun/tencent"` // 提供商：email用smtp；sms用aliyun/tencent
	Name       string      `json:"name"        orm:"name"        description:"配置名称"`                              // 配置名称
	IsActive   int         `json:"is_active"   orm:"is_active"   description:"是否启用：0-否，1-是（同渠道仅一个可启用）"`           // 是否启用：0-否，1-是（同渠道仅一个可启用）
	ConfigJson string      `json:"config_json" orm:"config_json" description:"JSON配置，不同provider结构不同"`             // JSON配置，不同provider结构不同
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`                               // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`                              // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`                               // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`                              // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除"`                              // 逻辑删除
}
