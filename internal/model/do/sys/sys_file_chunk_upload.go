// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysFileChunkUpload is the golang structure of table sys_file_chunk_upload for DAO operations like Where/Data.
type SysFileChunkUpload struct {
	g.Meta         `orm:"table:sys_file_chunk_upload, do:true"`
	Id             any         // 主键ID
	UploadId       any         // 上传任务ID
	Filename       any         // 原始文件名
	TotalSize      any         // 总大小
	ChunkSize      any         // 分片大小
	TotalChunks    any         // 总分片数
	UploadedChunks any         // 已上传分片数
	Status         any         // 状态：0-上传中，1-已完成，2-已取消
	FileId         any         // 完成后关联的sys_file_storage.id
	Creator        any         // 上传者
	CreateTime     *gtime.Time // 创建时间
	UpdateTime     *gtime.Time // 更新时间
	ExpireTime     *gtime.Time // 过期时间（未完成的分片任务）
}
