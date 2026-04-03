// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysTask is the golang structure for table sys_task.
type SysTask struct {
	Id             int64       `json:"id"              orm:"id"              description:"主键"`                  // 主键
	Name           string      `json:"name"            orm:"name"            description:"任务名称"`                // 任务名称
	Code           string      `json:"code"            orm:"code"            description:"任务编码（唯一）"`            // 任务编码（唯一）
	Type           int         `json:"type"            orm:"type"            description:"类型：1-定时(cron)，2-延迟"`  // 类型：1-定时(cron)，2-延迟
	CronExpr       string      `json:"cron_expr"       orm:"cron_expr"       description:"cron 表达式（type=1 时用）"` // cron 表达式（type=1 时用）
	DelaySeconds   int         `json:"delay_seconds"   orm:"delay_seconds"   description:"延迟秒数（type=2 时用）"`     // 延迟秒数（type=2 时用）
	Handler        string      `json:"handler"         orm:"handler"         description:"处理器（如 DemoHandler）"`  // 处理器（如 DemoHandler）
	Params         string      `json:"params"          orm:"params"          description:"参数 JSON"`             // 参数 JSON
	RetryTimes     int         `json:"retry_times"     orm:"retry_times"     description:"重试次数"`                // 重试次数
	RetryInterval  int         `json:"retry_interval"  orm:"retry_interval"  description:"重试间隔（秒）"`             // 重试间隔（秒）
	Concurrency    int         `json:"concurrency"     orm:"concurrency"     description:"并发度（0=不限制）"`          // 并发度（0=不限制）
	AlertOnFail    int         `json:"alert_on_fail"   orm:"alert_on_fail"   description:"失败是否告警：0-否，1-是"`      // 失败是否告警：0-否，1-是
	AlertReceivers string      `json:"alert_receivers" orm:"alert_receivers" description:"告警接收人（手机/邮箱，逗号分隔）"`   // 告警接收人（手机/邮箱，逗号分隔）
	Status         int         `json:"status"          orm:"status"          description:"状态：0-启用，1-停用"`        // 状态：0-启用，1-停用
	Remark         string      `json:"remark"          orm:"remark"          description:"备注"`                  // 备注
	Creator        string      `json:"creator"         orm:"creator"         description:"创建者"`                 // 创建者
	CreateTime     *gtime.Time `json:"create_time"     orm:"create_time"     description:"创建时间"`                // 创建时间
	Updater        string      `json:"updater"         orm:"updater"         description:"更新者"`                 // 更新者
	UpdateTime     *gtime.Time `json:"update_time"     orm:"update_time"     description:"更新时间"`                // 更新时间
	DeleteFlag     int         `json:"delete_flag"     orm:"delete_flag"     description:"逻辑删除"`                // 逻辑删除
}
