// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysExceptionLog is the golang structure for table sys_exception_log.
type SysExceptionLog struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键ID"`   // 主键ID
	TraceId    string      `json:"trace_id"    orm:"trace_id"    description:"链路追踪ID"` // 链路追踪ID
	Path       string      `json:"path"        orm:"path"        description:"请求路径"`   // 请求路径
	Method     string      `json:"method"      orm:"method"      description:"HTTP方法"` // HTTP方法
	ErrorMsg   string      `json:"error_msg"   orm:"error_msg"   description:"错误信息"`   // 错误信息
	Stack      string      `json:"stack"       orm:"stack"       description:"堆栈"`     // 堆栈
	UserId     int64       `json:"user_id"     orm:"user_id"     description:"用户ID"`   // 用户ID
	Ip         string      `json:"ip"          orm:"ip"          description:"客户端IP"`  // 客户端IP
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`   // 创建时间
}
