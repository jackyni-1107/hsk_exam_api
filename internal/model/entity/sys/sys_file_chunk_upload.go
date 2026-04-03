// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysFileChunkUpload is the golang structure for table sys_file_chunk_upload.
type SysFileChunkUpload struct {
	Id             int64       `json:"id"              orm:"id"              description:"主键ID"`                      // 主键ID
	UploadId       string      `json:"upload_id"       orm:"upload_id"       description:"上传任务ID"`                    // 上传任务ID
	Filename       string      `json:"filename"        orm:"filename"        description:"原始文件名"`                     // 原始文件名
	TotalSize      int64       `json:"total_size"      orm:"total_size"      description:"总大小"`                       // 总大小
	ChunkSize      int         `json:"chunk_size"      orm:"chunk_size"      description:"分片大小"`                      // 分片大小
	TotalChunks    int         `json:"total_chunks"    orm:"total_chunks"    description:"总分片数"`                      // 总分片数
	UploadedChunks int         `json:"uploaded_chunks" orm:"uploaded_chunks" description:"已上传分片数"`                    // 已上传分片数
	Status         int         `json:"status"          orm:"status"          description:"状态：0-上传中，1-已完成，2-已取消"`      // 状态：0-上传中，1-已完成，2-已取消
	FileId         int64       `json:"file_id"         orm:"file_id"         description:"完成后关联的sys_file_storage.id"` // 完成后关联的sys_file_storage.id
	Creator        string      `json:"creator"         orm:"creator"         description:"上传者"`                       // 上传者
	CreateTime     *gtime.Time `json:"create_time"     orm:"create_time"     description:"创建时间"`                      // 创建时间
	UpdateTime     *gtime.Time `json:"update_time"     orm:"update_time"     description:"更新时间"`                      // 更新时间
	ExpireTime     *gtime.Time `json:"expire_time"     orm:"expire_time"     description:"过期时间（未完成的分片任务）"`            // 过期时间（未完成的分片任务）
}
