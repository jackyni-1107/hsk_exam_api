// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysLoginLog is the golang structure for table sys_login_log.
type SysLoginLog struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`                               // 主键ID
	LogType    string      `json:"log_type"    orm:"log_type"    description:"类型：login_success/login_fail/logout"` // 类型：login_success/login_fail/logout
	UserId     int64       `json:"user_id"     orm:"user_id"     description:"用户ID"`                               // 用户ID
	Username   string      `json:"username"    orm:"username"    description:"用户名"`                                // 用户名
	UserType   int         `json:"user_type"   orm:"user_type"   description:"用户类型：1-后台，2-前台"`                     // 用户类型：1-后台，2-前台
	Ip         string      `json:"ip"          orm:"ip"          description:"客户端IP"`                              // 客户端IP
	UserAgent  string      `json:"user_agent"  orm:"user_agent"  description:"User-Agent"`                         // User-Agent
	DeviceInfo string      `json:"device_info" orm:"device_info" description:"设备信息JSON"`                           // 设备信息JSON
	TraceId    string      `json:"trace_id"    orm:"trace_id"    description:"链路追踪ID"`                             // 链路追踪ID
	FailReason string      `json:"fail_reason" orm:"fail_reason" description:"失败原因"`                               // 失败原因
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`                               // 创建时间
}
