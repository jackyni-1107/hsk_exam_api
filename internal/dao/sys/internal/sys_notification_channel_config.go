// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysNotificationChannelConfigDao is the data access object for the table sys_notification_channel_config.
type SysNotificationChannelConfigDao struct {
	table    string                              // table is the underlying table name of the DAO.
	group    string                              // group is the database configuration group name of the current DAO.
	columns  SysNotificationChannelConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                  // handlers for customized model modification.
}

// SysNotificationChannelConfigColumns defines and stores column names for the table sys_notification_channel_config.
type SysNotificationChannelConfigColumns struct {
	Id         string // 主键ID
	Channel    string // 渠道类型：email-邮件，sms-短信
	Provider   string // 提供商：email用smtp；sms用aliyun/tencent
	Name       string // 配置名称
	IsActive   string // 是否启用：0-否，1-是（同渠道仅一个可启用）
	ConfigJson string // JSON配置，不同provider结构不同
	Creator    string // 创建者
	CreateTime string // 创建时间
	Updater    string // 更新者
	UpdateTime string // 更新时间
	DeleteFlag string // 逻辑删除
}

// sysNotificationChannelConfigColumns holds the columns for the table sys_notification_channel_config.
var sysNotificationChannelConfigColumns = SysNotificationChannelConfigColumns{
	Id:         "id",
	Channel:    "channel",
	Provider:   "provider",
	Name:       "name",
	IsActive:   "is_active",
	ConfigJson: "config_json",
	Creator:    "creator",
	CreateTime: "create_time",
	Updater:    "updater",
	UpdateTime: "update_time",
	DeleteFlag: "delete_flag",
}

// NewSysNotificationChannelConfigDao creates and returns a new DAO object for table data access.
func NewSysNotificationChannelConfigDao(handlers ...gdb.ModelHandler) *SysNotificationChannelConfigDao {
	return &SysNotificationChannelConfigDao{
		group:    "default",
		table:    "sys_notification_channel_config",
		columns:  sysNotificationChannelConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysNotificationChannelConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysNotificationChannelConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysNotificationChannelConfigDao) Columns() SysNotificationChannelConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysNotificationChannelConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysNotificationChannelConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysNotificationChannelConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
