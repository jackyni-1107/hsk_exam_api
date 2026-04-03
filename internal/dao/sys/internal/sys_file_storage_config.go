// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysFileStorageConfigDao is the data access object for the table sys_file_storage_config.
type SysFileStorageConfigDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  SysFileStorageConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// SysFileStorageConfigColumns defines and stores column names for the table sys_file_storage_config.
type SysFileStorageConfigColumns struct {
	Id                string // 主键ID
	StorageType       string // 存储类型：local/oss/s3/minio
	Name              string // 配置名称
	IsActive          string // 是否启用：0-否，1-是（全局仅一个）
	ConfigJson        string // JSON配置。local: {"base_path":"./storage"}；oss/s3/minio: {"endpoint":"","bucket":"","access_key":"","secret_key":"","region":""}
	CleanupBeforeDays string // 文件清理策略：清理多少天前的孤立文件
	Creator           string // 创建者
	CreateTime        string // 创建时间
	Updater           string // 更新者
	UpdateTime        string // 更新时间
	DeleteFlag        string // 逻辑删除
}

// sysFileStorageConfigColumns holds the columns for the table sys_file_storage_config.
var sysFileStorageConfigColumns = SysFileStorageConfigColumns{
	Id:                "id",
	StorageType:       "storage_type",
	Name:              "name",
	IsActive:          "is_active",
	ConfigJson:        "config_json",
	CleanupBeforeDays: "cleanup_before_days",
	Creator:           "creator",
	CreateTime:        "create_time",
	Updater:           "updater",
	UpdateTime:        "update_time",
	DeleteFlag:        "delete_flag",
}

// NewSysFileStorageConfigDao creates and returns a new DAO object for table data access.
func NewSysFileStorageConfigDao(handlers ...gdb.ModelHandler) *SysFileStorageConfigDao {
	return &SysFileStorageConfigDao{
		group:    "default",
		table:    "sys_file_storage_config",
		columns:  sysFileStorageConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysFileStorageConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysFileStorageConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysFileStorageConfigDao) Columns() SysFileStorageConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysFileStorageConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysFileStorageConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysFileStorageConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
