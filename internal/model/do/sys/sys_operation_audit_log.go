// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysOperationAuditLog is the golang structure of table sys_operation_audit_log for DAO operations like Where/Data.
type SysOperationAuditLog struct {
	g.Meta       `orm:"table:sys_operation_audit_log, do:true"`
	Id           any         // 主键ID
	UserId       any         // 用户ID
	Username     any         // 用户名
	UserType     any         // 用户类型：1-后台用户，2-前台用户
	Module       any         // 模块
	Action       any         // 操作类型：create/update/delete/query
	LogType      any         // 日志类型：operation-操作, api_access-API访问
	Method       any         // HTTP方法
	Path         any         // 请求路径
	RequestData  any         // 请求数据
	ResponseData any         // 响应数据（仅create/update/delete记录）
	Ip           any         // 客户端IP
	UserAgent    any         // User-Agent
	TraceId      any         // 链路追踪ID
	DeviceInfo   any         // 设备信息JSON：device_type, os, browser
	DurationMs   any         // 耗时(毫秒)
	CreateTime   *gtime.Time // 创建时间
}
