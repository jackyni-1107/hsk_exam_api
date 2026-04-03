// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysUserDao is the data access object for the table sys_user.
type SysUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysUserColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysUserColumns defines and stores column names for the table sys_user.
type SysUserColumns struct {
	Id                 string // 用户ID
	Username           string // 用户账号
	Password           string // 密码
	PasswordChangedAt  string // 密码最后修改时间
	MustChangePassword string // 0否 1须改密
	Nickname           string // 用户昵称
	Remark             string // 备注
	Email              string // 用户邮箱
	Mobile             string // 手机号码
	Sex                string // 用户性别
	Avatar             string // 头像地址
	Status             string // 帐号状态（0正常 1停用）
	LoginIp            string // 最后登录IP
	LoginTime          string // 最后登录时间
	Creator            string // 创建者
	CreateTime         string // 创建时间
	Updater            string // 更新者
	UpdateTime         string // 更新时间
	DeleteFlag         string // 逻辑删除标识：0-未删除，1-已删除
}

// sysUserColumns holds the columns for the table sys_user.
var sysUserColumns = SysUserColumns{
	Id:                 "id",
	Username:           "username",
	Password:           "password",
	PasswordChangedAt:  "password_changed_at",
	MustChangePassword: "must_change_password",
	Nickname:           "nickname",
	Remark:             "remark",
	Email:              "email",
	Mobile:             "mobile",
	Sex:                "sex",
	Avatar:             "avatar",
	Status:             "status",
	LoginIp:            "login_ip",
	LoginTime:          "login_time",
	Creator:            "creator",
	CreateTime:         "create_time",
	Updater:            "updater",
	UpdateTime:         "update_time",
	DeleteFlag:         "delete_flag",
}

// NewSysUserDao creates and returns a new DAO object for table data access.
func NewSysUserDao(handlers ...gdb.ModelHandler) *SysUserDao {
	return &SysUserDao{
		group:    "default",
		table:    "sys_user",
		columns:  sysUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysUserDao) Columns() SysUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysUserDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
