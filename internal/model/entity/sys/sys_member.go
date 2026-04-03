// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMember is the golang structure for table sys_member.
type SysMember struct {
	Id                 int64       `json:"id"                   orm:"id"                   description:"用户ID"`               // 用户ID
	Username           string      `json:"username"             orm:"username"             description:"用户账号"`               // 用户账号
	Password           string      `json:"password"             orm:"password"             description:"密码"`                 // 密码
	PasswordChangedAt  *gtime.Time `json:"password_changed_at"  orm:"password_changed_at"  description:"密码最后修改时间"`           // 密码最后修改时间
	MustChangePassword int         `json:"must_change_password" orm:"must_change_password" description:"0否 1须改密"`            // 0否 1须改密
	Nickname           string      `json:"nickname"             orm:"nickname"             description:"用户昵称"`               // 用户昵称
	Avatar             string      `json:"avatar"               orm:"avatar"               description:"头像地址"`               // 头像地址
	Email              string      `json:"email"                orm:"email"                description:"用户邮箱"`               // 用户邮箱
	Mobile             string      `json:"mobile"               orm:"mobile"               description:"手机号码"`               // 手机号码
	Status             int         `json:"status"               orm:"status"               description:"帐号状态（0正常 1停用）"`      // 帐号状态（0正常 1停用）
	LoginIp            string      `json:"login_ip"             orm:"login_ip"             description:"最后登录IP"`             // 最后登录IP
	LoginTime          *gtime.Time `json:"login_time"           orm:"login_time"           description:"最后登录时间"`             // 最后登录时间
	Creator            string      `json:"creator"              orm:"creator"              description:"创建者"`                // 创建者
	CreateTime         *gtime.Time `json:"create_time"          orm:"create_time"          description:"创建时间"`               // 创建时间
	Updater            string      `json:"updater"              orm:"updater"              description:"更新者"`                // 更新者
	UpdateTime         *gtime.Time `json:"update_time"          orm:"update_time"          description:"更新时间"`               // 更新时间
	DeleteFlag         int         `json:"delete_flag"          orm:"delete_flag"          description:"逻辑删除标识：0-未删除，1-已删除"` // 逻辑删除标识：0-未删除，1-已删除
}
