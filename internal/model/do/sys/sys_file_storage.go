// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysFileStorage is the golang structure of table sys_file_storage for DAO operations like Where/Data.
type SysFileStorage struct {
	g.Meta      `orm:"table:sys_file_storage, do:true"`
	Id          any         // 主键ID
	StorageType any         // 存储类型：local/oss/s3/minio
	Bucket      any         // 桶名（OSS/S3/MinIO）
	Path        any         // 存储路径
	Filename    any         // 原始文件名
	MimeType    any         // MIME类型
	Size        any         // 文件大小(字节)
	Hash        any         // 文件哈希（用于秒传/去重）
	IsPrivate   any         // 是否私有：0-公开，1-私有
	Creator     any         // 上传者
	CreateTime  *gtime.Time // 创建时间
	DeleteFlag  any         // 逻辑删除
}
