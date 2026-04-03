// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysPasswordHistoryDao is the data access object for the table sys_password_history.
type SysPasswordHistoryDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  SysPasswordHistoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// SysPasswordHistoryColumns defines and stores column names for the table sys_password_history.
type SysPasswordHistoryColumns struct {
	Id           string //
	UserType     string // 1=admin 2=client
	UserId       string //
	PasswordHash string //
	CreatedAt    string //
}

// sysPasswordHistoryColumns holds the columns for the table sys_password_history.
var sysPasswordHistoryColumns = SysPasswordHistoryColumns{
	Id:           "id",
	UserType:     "user_type",
	UserId:       "user_id",
	PasswordHash: "password_hash",
	CreatedAt:    "created_at",
}

// NewSysPasswordHistoryDao creates and returns a new DAO object for table data access.
func NewSysPasswordHistoryDao(handlers ...gdb.ModelHandler) *SysPasswordHistoryDao {
	return &SysPasswordHistoryDao{
		group:    "default",
		table:    "sys_password_history",
		columns:  sysPasswordHistoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysPasswordHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysPasswordHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysPasswordHistoryDao) Columns() SysPasswordHistoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysPasswordHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysPasswordHistoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysPasswordHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
