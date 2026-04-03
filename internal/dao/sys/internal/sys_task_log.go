// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysTaskLogDao is the data access object for the table sys_task_log.
type SysTaskLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysTaskLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysTaskLogColumns defines and stores column names for the table sys_task_log.
type SysTaskLogColumns struct {
	Id          string // 主键
	TaskId      string // 任务ID
	RunId       string // 执行批次ID
	TriggerType string // 1-定时，2-延迟，3-手动
	Status      string // 0-执行中，1-成功，2-失败
	StartTime   string // 开始时间
	EndTime     string // 结束时间
	DurationMs  string // 耗时(ms)
	RetryCount  string // 已重试次数
	ErrorMsg    string // 错误信息
	Result      string // 执行结果
	Node        string // 执行节点（hostname）
	CreateTime  string // 创建时间
}

// sysTaskLogColumns holds the columns for the table sys_task_log.
var sysTaskLogColumns = SysTaskLogColumns{
	Id:          "id",
	TaskId:      "task_id",
	RunId:       "run_id",
	TriggerType: "trigger_type",
	Status:      "status",
	StartTime:   "start_time",
	EndTime:     "end_time",
	DurationMs:  "duration_ms",
	RetryCount:  "retry_count",
	ErrorMsg:    "error_msg",
	Result:      "result",
	Node:        "node",
	CreateTime:  "create_time",
}

// NewSysTaskLogDao creates and returns a new DAO object for table data access.
func NewSysTaskLogDao(handlers ...gdb.ModelHandler) *SysTaskLogDao {
	return &SysTaskLogDao{
		group:    "default",
		table:    "sys_task_log",
		columns:  sysTaskLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysTaskLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysTaskLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysTaskLogDao) Columns() SysTaskLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysTaskLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysTaskLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysTaskLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
