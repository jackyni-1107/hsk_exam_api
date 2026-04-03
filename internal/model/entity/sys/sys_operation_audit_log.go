// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysOperationAuditLog is the golang structure for table sys_operation_audit_log.
type SysOperationAuditLog struct {
	Id           int64       `json:"id"            orm:"id"            description:"主键ID"`                                // 主键ID
	UserId       int64       `json:"user_id"       orm:"user_id"       description:"用户ID"`                                // 用户ID
	Username     string      `json:"username"      orm:"username"      description:"用户名"`                                 // 用户名
	UserType     int         `json:"user_type"     orm:"user_type"     description:"用户类型：1-后台用户，2-前台用户"`                  // 用户类型：1-后台用户，2-前台用户
	Module       string      `json:"module"        orm:"module"        description:"模块"`                                  // 模块
	Action       string      `json:"action"        orm:"action"        description:"操作类型：create/update/delete/query"`     // 操作类型：create/update/delete/query
	LogType      string      `json:"log_type"      orm:"log_type"      description:"日志类型：operation-操作, api_access-API访问"` // 日志类型：operation-操作, api_access-API访问
	Method       string      `json:"method"        orm:"method"        description:"HTTP方法"`                              // HTTP方法
	Path         string      `json:"path"          orm:"path"          description:"请求路径"`                                // 请求路径
	RequestData  string      `json:"request_data"  orm:"request_data"  description:"请求数据"`                                // 请求数据
	ResponseData string      `json:"response_data" orm:"response_data" description:"响应数据（仅create/update/delete记录）"`       // 响应数据（仅create/update/delete记录）
	Ip           string      `json:"ip"            orm:"ip"            description:"客户端IP"`                               // 客户端IP
	UserAgent    string      `json:"user_agent"    orm:"user_agent"    description:"User-Agent"`                          // User-Agent
	TraceId      string      `json:"trace_id"      orm:"trace_id"      description:"链路追踪ID"`                              // 链路追踪ID
	DeviceInfo   string      `json:"device_info"   orm:"device_info"   description:"设备信息JSON：device_type, os, browser"`   // 设备信息JSON：device_type, os, browser
	DurationMs   int         `json:"duration_ms"   orm:"duration_ms"   description:"耗时(毫秒)"`                              // 耗时(毫秒)
	CreateTime   *gtime.Time `json:"create_time"   orm:"create_time"   description:"创建时间"`                                // 创建时间
}
