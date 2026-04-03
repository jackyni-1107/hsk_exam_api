// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysLoginLog is the golang structure of table sys_login_log for DAO operations like Where/Data.
type SysLoginLog struct {
	g.Meta     `orm:"table:sys_login_log, do:true"`
	Id         any         // 主键ID
	LogType    any         // 类型：login_success/login_fail/logout
	UserId     any         // 用户ID
	Username   any         // 用户名
	UserType   any         // 用户类型：1-后台，2-前台
	Ip         any         // 客户端IP
	UserAgent  any         // User-Agent
	DeviceInfo any         // 设备信息JSON
	TraceId    any         // 链路追踪ID
	FailReason any         // 失败原因
	CreateTime *gtime.Time // 创建时间
}
