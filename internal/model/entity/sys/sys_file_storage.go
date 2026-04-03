// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysFileStorage is the golang structure for table sys_file_storage.
type SysFileStorage struct {
	Id          int64       `json:"id"           orm:"id"           description:"主键ID"`                    // 主键ID
	StorageType string      `json:"storage_type" orm:"storage_type" description:"存储类型：local/oss/s3/minio"` // 存储类型：local/oss/s3/minio
	Bucket      string      `json:"bucket"       orm:"bucket"       description:"桶名（OSS/S3/MinIO）"`        // 桶名（OSS/S3/MinIO）
	Path        string      `json:"path"         orm:"path"         description:"存储路径"`                    // 存储路径
	Filename    string      `json:"filename"     orm:"filename"     description:"原始文件名"`                   // 原始文件名
	MimeType    string      `json:"mime_type"    orm:"mime_type"    description:"MIME类型"`                  // MIME类型
	Size        int64       `json:"size"         orm:"size"         description:"文件大小(字节)"`                // 文件大小(字节)
	Hash        string      `json:"hash"         orm:"hash"         description:"文件哈希（用于秒传/去重）"`           // 文件哈希（用于秒传/去重）
	IsPrivate   int         `json:"is_private"   orm:"is_private"   description:"是否私有：0-公开，1-私有"`          // 是否私有：0-公开，1-私有
	Creator     string      `json:"creator"      orm:"creator"      description:"上传者"`                     // 上传者
	CreateTime  *gtime.Time `json:"create_time"  orm:"create_time"  description:"创建时间"`                    // 创建时间
	DeleteFlag  int         `json:"delete_flag"  orm:"delete_flag"  description:"逻辑删除"`                    // 逻辑删除
}
