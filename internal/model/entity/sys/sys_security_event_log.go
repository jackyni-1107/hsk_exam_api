// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysSecurityEventLog is the golang structure for table sys_security_event_log.
type SysSecurityEventLog struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`                                                           // 主键ID
	EventType  string      `json:"event_type"  orm:"event_type"  description:"事件类型：token_invalid/permission_denied/brute_force/suspicious_ip"` // 事件类型：token_invalid/permission_denied/brute_force/suspicious_ip
	UserId     int64       `json:"user_id"     orm:"user_id"     description:"用户ID"`                                                           // 用户ID
	Ip         string      `json:"ip"          orm:"ip"          description:"客户端IP"`                                                          // 客户端IP
	UserAgent  string      `json:"user_agent"  orm:"user_agent"  description:"User-Agent"`                                                     // User-Agent
	Detail     string      `json:"detail"      orm:"detail"      description:"详情"`                                                             // 详情
	TraceId    string      `json:"trace_id"    orm:"trace_id"    description:"链路追踪ID"`                                                         // 链路追踪ID
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`                                                           // 创建时间
}
