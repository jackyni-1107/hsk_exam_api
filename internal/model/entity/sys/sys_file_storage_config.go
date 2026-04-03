// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysFileStorageConfig is the golang structure for table sys_file_storage_config.
type SysFileStorageConfig struct {
	Id                int64       `json:"id"                  orm:"id"                  description:"主键ID"`                                                                                                                                                  // 主键ID
	StorageType       string      `json:"storage_type"        orm:"storage_type"        description:"存储类型：local/oss/s3/minio"`                                                                                                                               // 存储类型：local/oss/s3/minio
	Name              string      `json:"name"                orm:"name"                description:"配置名称"`                                                                                                                                                  // 配置名称
	IsActive          int         `json:"is_active"           orm:"is_active"           description:"是否启用：0-否，1-是（全局仅一个）"`                                                                                                                                   // 是否启用：0-否，1-是（全局仅一个）
	ConfigJson        string      `json:"config_json"         orm:"config_json"         description:"JSON配置。local: {\"base_path\":\"./storage\"}；oss/s3/minio: {\"endpoint\":\"\",\"bucket\":\"\",\"access_key\":\"\",\"secret_key\":\"\",\"region\":\"\"}"` // JSON配置。local: {"base_path":"./storage"}；oss/s3/minio: {"endpoint":"","bucket":"","access_key":"","secret_key":"","region":""}
	CleanupBeforeDays int         `json:"cleanup_before_days" orm:"cleanup_before_days" description:"文件清理策略：清理多少天前的孤立文件"`                                                                                                                                    // 文件清理策略：清理多少天前的孤立文件
	Creator           string      `json:"creator"             orm:"creator"             description:"创建者"`                                                                                                                                                   // 创建者
	CreateTime        *gtime.Time `json:"create_time"         orm:"create_time"         description:"创建时间"`                                                                                                                                                  // 创建时间
	Updater           string      `json:"updater"             orm:"updater"             description:"更新者"`                                                                                                                                                   // 更新者
	UpdateTime        *gtime.Time `json:"update_time"         orm:"update_time"         description:"更新时间"`                                                                                                                                                  // 更新时间
	DeleteFlag        int         `json:"delete_flag"         orm:"delete_flag"         description:"逻辑删除"`                                                                                                                                                  // 逻辑删除
}
