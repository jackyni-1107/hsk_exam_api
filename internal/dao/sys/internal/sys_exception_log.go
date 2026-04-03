// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysExceptionLogDao is the data access object for the table sys_exception_log.
type SysExceptionLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  SysExceptionLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// SysExceptionLogColumns defines and stores column names for the table sys_exception_log.
type SysExceptionLogColumns struct {
	Id         string // 主键ID
	TraceId    string // 链路追踪ID
	Path       string // 请求路径
	Method     string // HTTP方法
	ErrorMsg   string // 错误信息
	Stack      string // 堆栈
	UserId     string // 用户ID
	Ip         string // 客户端IP
	CreateTime string // 创建时间
}

// sysExceptionLogColumns holds the columns for the table sys_exception_log.
var sysExceptionLogColumns = SysExceptionLogColumns{
	Id:         "id",
	TraceId:    "trace_id",
	Path:       "path",
	Method:     "method",
	ErrorMsg:   "error_msg",
	Stack:      "stack",
	UserId:     "user_id",
	Ip:         "ip",
	CreateTime: "create_time",
}

// NewSysExceptionLogDao creates and returns a new DAO object for table data access.
func NewSysExceptionLogDao(handlers ...gdb.ModelHandler) *SysExceptionLogDao {
	return &SysExceptionLogDao{
		group:    "default",
		table:    "sys_exception_log",
		columns:  sysExceptionLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysExceptionLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysExceptionLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysExceptionLogDao) Columns() SysExceptionLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysExceptionLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysExceptionLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysExceptionLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
