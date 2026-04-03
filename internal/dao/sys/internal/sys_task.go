// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysTaskDao is the data access object for the table sys_task.
type SysTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysTaskColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysTaskColumns defines and stores column names for the table sys_task.
type SysTaskColumns struct {
	Id             string // 主键
	Name           string // 任务名称
	Code           string // 任务编码（唯一）
	Type           string // 类型：1-定时(cron)，2-延迟
	CronExpr       string // cron 表达式（type=1 时用）
	DelaySeconds   string // 延迟秒数（type=2 时用）
	Handler        string // 处理器（如 DemoHandler）
	Params         string // 参数 JSON
	RetryTimes     string // 重试次数
	RetryInterval  string // 重试间隔（秒）
	Concurrency    string // 并发度（0=不限制）
	AlertOnFail    string // 失败是否告警：0-否，1-是
	AlertReceivers string // 告警接收人（手机/邮箱，逗号分隔）
	Status         string // 状态：0-启用，1-停用
	Remark         string // 备注
	Creator        string // 创建者
	CreateTime     string // 创建时间
	Updater        string // 更新者
	UpdateTime     string // 更新时间
	DeleteFlag     string // 逻辑删除
}

// sysTaskColumns holds the columns for the table sys_task.
var sysTaskColumns = SysTaskColumns{
	Id:             "id",
	Name:           "name",
	Code:           "code",
	Type:           "type",
	CronExpr:       "cron_expr",
	DelaySeconds:   "delay_seconds",
	Handler:        "handler",
	Params:         "params",
	RetryTimes:     "retry_times",
	RetryInterval:  "retry_interval",
	Concurrency:    "concurrency",
	AlertOnFail:    "alert_on_fail",
	AlertReceivers: "alert_receivers",
	Status:         "status",
	Remark:         "remark",
	Creator:        "creator",
	CreateTime:     "create_time",
	Updater:        "updater",
	UpdateTime:     "update_time",
	DeleteFlag:     "delete_flag",
}

// NewSysTaskDao creates and returns a new DAO object for table data access.
func NewSysTaskDao(handlers ...gdb.ModelHandler) *SysTaskDao {
	return &SysTaskDao{
		group:    "default",
		table:    "sys_task",
		columns:  sysTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysTaskDao) Columns() SysTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
