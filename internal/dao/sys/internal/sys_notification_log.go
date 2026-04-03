// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysNotificationLogDao is the data access object for the table sys_notification_log.
type SysNotificationLogDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  SysNotificationLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// SysNotificationLogColumns defines and stores column names for the table sys_notification_log.
type SysNotificationLogColumns struct {
	Id           string // 主键ID
	TemplateCode string // 模板编码
	Channel      string // 渠道：sms/email/template
	Recipient    string // 接收者（手机号/邮箱/用户ID）
	Content      string // 发送内容
	Status       string // 状态：0-待发送，1-成功，2-失败
	ErrorMsg     string // 失败原因
	CreateTime   string // 创建时间
}

// sysNotificationLogColumns holds the columns for the table sys_notification_log.
var sysNotificationLogColumns = SysNotificationLogColumns{
	Id:           "id",
	TemplateCode: "template_code",
	Channel:      "channel",
	Recipient:    "recipient",
	Content:      "content",
	Status:       "status",
	ErrorMsg:     "error_msg",
	CreateTime:   "create_time",
}

// NewSysNotificationLogDao creates and returns a new DAO object for table data access.
func NewSysNotificationLogDao(handlers ...gdb.ModelHandler) *SysNotificationLogDao {
	return &SysNotificationLogDao{
		group:    "default",
		table:    "sys_notification_log",
		columns:  sysNotificationLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysNotificationLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysNotificationLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysNotificationLogDao) Columns() SysNotificationLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysNotificationLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysNotificationLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysNotificationLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
