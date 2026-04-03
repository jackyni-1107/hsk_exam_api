// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysNotificationChannelConfig is the golang structure of table sys_notification_channel_config for DAO operations like Where/Data.
type SysNotificationChannelConfig struct {
	g.Meta     `orm:"table:sys_notification_channel_config, do:true"`
	Id         any         // 主键ID
	Channel    any         // 渠道类型：email-邮件，sms-短信
	Provider   any         // 提供商：email用smtp；sms用aliyun/tencent
	Name       any         // 配置名称
	IsActive   any         // 是否启用：0-否，1-是（同渠道仅一个可启用）
	ConfigJson any         // JSON配置，不同provider结构不同
	Creator    any         // 创建者
	CreateTime *gtime.Time // 创建时间
	Updater    any         // 更新者
	UpdateTime *gtime.Time // 更新时间
	DeleteFlag any         // 逻辑删除
}
