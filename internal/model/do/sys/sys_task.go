// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysTask is the golang structure of table sys_task for DAO operations like Where/Data.
type SysTask struct {
	g.Meta         `orm:"table:sys_task, do:true"`
	Id             any         // 主键
	Name           any         // 任务名称
	Code           any         // 任务编码（唯一）
	Type           any         // 类型：1-定时(cron)，2-延迟
	CronExpr       any         // cron 表达式（type=1 时用）
	DelaySeconds   any         // 延迟秒数（type=2 时用）
	Handler        any         // 处理器（如 DemoHandler）
	Params         any         // 参数 JSON
	RetryTimes     any         // 重试次数
	RetryInterval  any         // 重试间隔（秒）
	Concurrency    any         // 并发度（0=不限制）
	AlertOnFail    any         // 失败是否告警：0-否，1-是
	AlertReceivers any         // 告警接收人（手机/邮箱，逗号分隔）
	Status         any         // 状态：0-启用，1-停用
	Remark         any         // 备注
	Creator        any         // 创建者
	CreateTime     *gtime.Time // 创建时间
	Updater        any         // 更新者
	UpdateTime     *gtime.Time // 更新时间
	DeleteFlag     any         // 逻辑删除
}
