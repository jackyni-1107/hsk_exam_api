// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysOperationAuditLogDao is the data access object for the table sys_operation_audit_log.
type SysOperationAuditLogDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  SysOperationAuditLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// SysOperationAuditLogColumns defines and stores column names for the table sys_operation_audit_log.
type SysOperationAuditLogColumns struct {
	Id           string // 主键ID
	UserId       string // 用户ID
	Username     string // 用户名
	UserType     string // 用户类型：1-后台用户，2-前台用户
	Module       string // 模块
	Action       string // 操作类型：create/update/delete/query
	LogType      string // 日志类型：operation-操作, api_access-API访问
	Method       string // HTTP方法
	Path         string // 请求路径
	RequestData  string // 请求数据
	ResponseData string // 响应数据（仅create/update/delete记录）
	Ip           string // 客户端IP
	UserAgent    string // User-Agent
	TraceId      string // 链路追踪ID
	DeviceInfo   string // 设备信息JSON：device_type, os, browser
	DurationMs   string // 耗时(毫秒)
	CreateTime   string // 创建时间
}

// sysOperationAuditLogColumns holds the columns for the table sys_operation_audit_log.
var sysOperationAuditLogColumns = SysOperationAuditLogColumns{
	Id:           "id",
	UserId:       "user_id",
	Username:     "username",
	UserType:     "user_type",
	Module:       "module",
	Action:       "action",
	LogType:      "log_type",
	Method:       "method",
	Path:         "path",
	RequestData:  "request_data",
	ResponseData: "response_data",
	Ip:           "ip",
	UserAgent:    "user_agent",
	TraceId:      "trace_id",
	DeviceInfo:   "device_info",
	DurationMs:   "duration_ms",
	CreateTime:   "create_time",
}

// NewSysOperationAuditLogDao creates and returns a new DAO object for table data access.
func NewSysOperationAuditLogDao(handlers ...gdb.ModelHandler) *SysOperationAuditLogDao {
	return &SysOperationAuditLogDao{
		group:    "default",
		table:    "sys_operation_audit_log",
		columns:  sysOperationAuditLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysOperationAuditLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysOperationAuditLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysOperationAuditLogDao) Columns() SysOperationAuditLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysOperationAuditLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysOperationAuditLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysOperationAuditLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
