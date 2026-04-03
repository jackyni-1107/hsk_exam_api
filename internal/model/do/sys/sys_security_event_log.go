// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysSecurityEventLog is the golang structure of table sys_security_event_log for DAO operations like Where/Data.
type SysSecurityEventLog struct {
	g.Meta     `orm:"table:sys_security_event_log, do:true"`
	Id         any         // 主键ID
	EventType  any         // 事件类型：token_invalid/permission_denied/brute_force/suspicious_ip
	UserId     any         // 用户ID
	Ip         any         // 客户端IP
	UserAgent  any         // User-Agent
	Detail     any         // 详情
	TraceId    any         // 链路追踪ID
	CreateTime *gtime.Time // 创建时间
}
