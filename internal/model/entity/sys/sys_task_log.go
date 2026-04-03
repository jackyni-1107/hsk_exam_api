// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysTaskLog is the golang structure for table sys_task_log.
type SysTaskLog struct {
	Id          int64       `json:"id"           orm:"id"           description:"主键"`              // 主键
	TaskId      int64       `json:"task_id"      orm:"task_id"      description:"任务ID"`            // 任务ID
	RunId       string      `json:"run_id"       orm:"run_id"       description:"执行批次ID"`          // 执行批次ID
	TriggerType int         `json:"trigger_type" orm:"trigger_type" description:"1-定时，2-延迟，3-手动"`  // 1-定时，2-延迟，3-手动
	Status      int         `json:"status"       orm:"status"       description:"0-执行中，1-成功，2-失败"` // 0-执行中，1-成功，2-失败
	StartTime   *gtime.Time `json:"start_time"   orm:"start_time"   description:"开始时间"`            // 开始时间
	EndTime     *gtime.Time `json:"end_time"     orm:"end_time"     description:"结束时间"`            // 结束时间
	DurationMs  int         `json:"duration_ms"  orm:"duration_ms"  description:"耗时(ms)"`          // 耗时(ms)
	RetryCount  int         `json:"retry_count"  orm:"retry_count"  description:"已重试次数"`           // 已重试次数
	ErrorMsg    string      `json:"error_msg"    orm:"error_msg"    description:"错误信息"`            // 错误信息
	Result      string      `json:"result"       orm:"result"       description:"执行结果"`            // 执行结果
	Node        string      `json:"node"         orm:"node"         description:"执行节点（hostname）"`  // 执行节点（hostname）
	CreateTime  *gtime.Time `json:"create_time"  orm:"create_time"  description:"创建时间"`            // 创建时间
}
