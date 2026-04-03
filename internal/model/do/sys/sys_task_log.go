// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysTaskLog is the golang structure of table sys_task_log for DAO operations like Where/Data.
type SysTaskLog struct {
	g.Meta      `orm:"table:sys_task_log, do:true"`
	Id          any         // 主键
	TaskId      any         // 任务ID
	RunId       any         // 执行批次ID
	TriggerType any         // 1-定时，2-延迟，3-手动
	Status      any         // 0-执行中，1-成功，2-失败
	StartTime   *gtime.Time // 开始时间
	EndTime     *gtime.Time // 结束时间
	DurationMs  any         // 耗时(ms)
	RetryCount  any         // 已重试次数
	ErrorMsg    any         // 错误信息
	Result      any         // 执行结果
	Node        any         // 执行节点（hostname）
	CreateTime  *gtime.Time // 创建时间
}
