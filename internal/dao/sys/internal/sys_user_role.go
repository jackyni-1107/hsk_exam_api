// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysUserRoleDao is the data access object for the table sys_user_role.
type SysUserRoleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysUserRoleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysUserRoleColumns defines and stores column names for the table sys_user_role.
type SysUserRoleColumns struct {
	Id         string // 自增编号
	UserId     string // 用户ID
	RoleId     string // 角色ID
	Creator    string // 创建者
	CreateTime string // 创建时间
	Updater    string // 更新者
	UpdateTime string // 更新时间
	DeleteFlag string // 逻辑删除标识：0-未删除，1-已删除
}

// sysUserRoleColumns holds the columns for the table sys_user_role.
var sysUserRoleColumns = SysUserRoleColumns{
	Id:         "id",
	UserId:     "user_id",
	RoleId:     "role_id",
	Creator:    "creator",
	CreateTime: "create_time",
	Updater:    "updater",
	UpdateTime: "update_time",
	DeleteFlag: "delete_flag",
}

// NewSysUserRoleDao creates and returns a new DAO object for table data access.
func NewSysUserRoleDao(handlers ...gdb.ModelHandler) *SysUserRoleDao {
	return &SysUserRoleDao{
		group:    "default",
		table:    "sys_user_role",
		columns:  sysUserRoleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysUserRoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysUserRoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysUserRoleDao) Columns() SysUserRoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysUserRoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysUserRoleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysUserRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
