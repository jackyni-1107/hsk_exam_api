// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysNotificationTemplateDao is the data access object for the table sys_notification_template.
type SysNotificationTemplateDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  SysNotificationTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// SysNotificationTemplateColumns defines and stores column names for the table sys_notification_template.
type SysNotificationTemplateColumns struct {
	Id                       string // 主键ID
	Code                     string // 模板编码
	Name                     string // 模板名称
	Channel                  string // 渠道：sms/email/template
	ChannelConfigId          string // 绑定的通知渠道配置ID（sys_notification_channel_config.id）
	TemplateType             string // 模板类型：1=系统模板 2=第三方模板
	Content                  string // 模板内容，支持变量 {{var}}
	ThirdPartyTemplateId     string // 第三方模板ID
	ThirdPartyTemplateParams string // 第三方模板参数(JSON)
	Variables                string // 变量列表，逗号分隔
	Status                   string // 状态：0-启用，1-停用
	Remark                   string // 备注
	Creator                  string // 创建者
	CreateTime               string // 创建时间
	Updater                  string // 更新者
	UpdateTime               string // 更新时间
	DeleteFlag               string // 逻辑删除
}

// sysNotificationTemplateColumns holds the columns for the table sys_notification_template.
var sysNotificationTemplateColumns = SysNotificationTemplateColumns{
	Id:                       "id",
	Code:                     "code",
	Name:                     "name",
	Channel:                  "channel",
	ChannelConfigId:          "channel_config_id",
	TemplateType:             "template_type",
	Content:                  "content",
	ThirdPartyTemplateId:     "third_party_template_id",
	ThirdPartyTemplateParams: "third_party_template_params",
	Variables:                "variables",
	Status:                   "status",
	Remark:                   "remark",
	Creator:                  "creator",
	CreateTime:               "create_time",
	Updater:                  "updater",
	UpdateTime:               "update_time",
	DeleteFlag:               "delete_flag",
}

// NewSysNotificationTemplateDao creates and returns a new DAO object for table data access.
func NewSysNotificationTemplateDao(handlers ...gdb.ModelHandler) *SysNotificationTemplateDao {
	return &SysNotificationTemplateDao{
		group:    "default",
		table:    "sys_notification_template",
		columns:  sysNotificationTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysNotificationTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysNotificationTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysNotificationTemplateDao) Columns() SysNotificationTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysNotificationTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysNotificationTemplateDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysNotificationTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
