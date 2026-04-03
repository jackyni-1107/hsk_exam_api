// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysPasswordHistory is the golang structure of table sys_password_history for DAO operations like Where/Data.
type SysPasswordHistory struct {
	g.Meta       `orm:"table:sys_password_history, do:true"`
	Id           any         //
	UserType     any         // 1=admin 2=client
	UserId       any         //
	PasswordHash any         //
	CreatedAt    *gtime.Time //
}
