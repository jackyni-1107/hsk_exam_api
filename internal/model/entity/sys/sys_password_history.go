// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysPasswordHistory is the golang structure for table sys_password_history.
type SysPasswordHistory struct {
	Id           int64       `json:"id"            orm:"id"            description:""`                 //
	UserType     int         `json:"user_type"     orm:"user_type"     description:"1=admin 2=client"` // 1=admin 2=client
	UserId       int64       `json:"user_id"       orm:"user_id"       description:""`                 //
	PasswordHash string      `json:"password_hash" orm:"password_hash" description:""`                 //
	CreatedAt    *gtime.Time `json:"created_at"    orm:"created_at"    description:""`                 //
}
