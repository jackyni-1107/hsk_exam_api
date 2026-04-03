// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure of table sys_user for DAO operations like Where/Data.
type SysUser struct {
	g.Meta             `orm:"table:sys_user, do:true"`
	Id                 any         // 用户ID
	Username           any         // 用户账号
	Password           any         // 密码
	PasswordChangedAt  *gtime.Time // 密码最后修改时间
	MustChangePassword any         // 0否 1须改密
	Nickname           any         // 用户昵称
	Remark             any         // 备注
	Email              any         // 用户邮箱
	Mobile             any         // 手机号码
	Sex                any         // 用户性别
	Avatar             any         // 头像地址
	Status             any         // 帐号状态（0正常 1停用）
	LoginIp            any         // 最后登录IP
	LoginTime          *gtime.Time // 最后登录时间
	Creator            any         // 创建者
	CreateTime         *gtime.Time // 创建时间
	Updater            any         // 更新者
	UpdateTime         *gtime.Time // 更新时间
	DeleteFlag         any         // 逻辑删除标识：0-未删除，1-已删除
}
