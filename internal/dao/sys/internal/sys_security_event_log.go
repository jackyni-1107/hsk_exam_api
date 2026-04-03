// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysSecurityEventLogDao is the data access object for the table sys_security_event_log.
type SysSecurityEventLogDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  SysSecurityEventLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// SysSecurityEventLogColumns defines and stores column names for the table sys_security_event_log.
type SysSecurityEventLogColumns struct {
	Id         string // 主键ID
	EventType  string // 事件类型：token_invalid/permission_denied/brute_force/suspicious_ip
	UserId     string // 用户ID
	Ip         string // 客户端IP
	UserAgent  string // User-Agent
	Detail     string // 详情
	TraceId    string // 链路追踪ID
	CreateTime string // 创建时间
}

// sysSecurityEventLogColumns holds the columns for the table sys_security_event_log.
var sysSecurityEventLogColumns = SysSecurityEventLogColumns{
	Id:         "id",
	EventType:  "event_type",
	UserId:     "user_id",
	Ip:         "ip",
	UserAgent:  "user_agent",
	Detail:     "detail",
	TraceId:    "trace_id",
	CreateTime: "create_time",
}

// NewSysSecurityEventLogDao creates and returns a new DAO object for table data access.
func NewSysSecurityEventLogDao(handlers ...gdb.ModelHandler) *SysSecurityEventLogDao {
	return &SysSecurityEventLogDao{
		group:    "default",
		table:    "sys_security_event_log",
		columns:  sysSecurityEventLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysSecurityEventLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysSecurityEventLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysSecurityEventLogDao) Columns() SysSecurityEventLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysSecurityEventLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysSecurityEventLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysSecurityEventLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
