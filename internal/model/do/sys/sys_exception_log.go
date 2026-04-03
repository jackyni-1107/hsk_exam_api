// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysExceptionLog is the golang structure of table sys_exception_log for DAO operations like Where/Data.
type SysExceptionLog struct {
	g.Meta     `orm:"table:sys_exception_log, do:true"`
	Id         any         // 主键ID
	TraceId    any         // 链路追踪ID
	Path       any         // 请求路径
	Method     any         // HTTP方法
	ErrorMsg   any         // 错误信息
	Stack      any         // 堆栈
	UserId     any         // 用户ID
	Ip         any         // 客户端IP
	CreateTime *gtime.Time // 创建时间
}
