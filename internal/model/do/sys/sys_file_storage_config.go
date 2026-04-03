// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysFileStorageConfig is the golang structure of table sys_file_storage_config for DAO operations like Where/Data.
type SysFileStorageConfig struct {
	g.Meta            `orm:"table:sys_file_storage_config, do:true"`
	Id                any         // 主键ID
	StorageType       any         // 存储类型：local/oss/s3/minio
	Name              any         // 配置名称
	IsActive          any         // 是否启用：0-否，1-是（全局仅一个）
	ConfigJson        any         // JSON配置。local: {"base_path":"./storage"}；oss/s3/minio: {"endpoint":"","bucket":"","access_key":"","secret_key":"","region":""}
	CleanupBeforeDays any         // 文件清理策略：清理多少天前的孤立文件
	Creator           any         // 创建者
	CreateTime        *gtime.Time // 创建时间
	Updater           any         // 更新者
	UpdateTime        *gtime.Time // 更新时间
	DeleteFlag        any         // 逻辑删除
}
