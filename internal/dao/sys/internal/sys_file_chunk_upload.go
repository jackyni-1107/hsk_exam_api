// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysFileChunkUploadDao is the data access object for the table sys_file_chunk_upload.
type SysFileChunkUploadDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  SysFileChunkUploadColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// SysFileChunkUploadColumns defines and stores column names for the table sys_file_chunk_upload.
type SysFileChunkUploadColumns struct {
	Id             string // 主键ID
	UploadId       string // 上传任务ID
	Filename       string // 原始文件名
	TotalSize      string // 总大小
	ChunkSize      string // 分片大小
	TotalChunks    string // 总分片数
	UploadedChunks string // 已上传分片数
	Status         string // 状态：0-上传中，1-已完成，2-已取消
	FileId         string // 完成后关联的sys_file_storage.id
	Creator        string // 上传者
	CreateTime     string // 创建时间
	UpdateTime     string // 更新时间
	ExpireTime     string // 过期时间（未完成的分片任务）
}

// sysFileChunkUploadColumns holds the columns for the table sys_file_chunk_upload.
var sysFileChunkUploadColumns = SysFileChunkUploadColumns{
	Id:             "id",
	UploadId:       "upload_id",
	Filename:       "filename",
	TotalSize:      "total_size",
	ChunkSize:      "chunk_size",
	TotalChunks:    "total_chunks",
	UploadedChunks: "uploaded_chunks",
	Status:         "status",
	FileId:         "file_id",
	Creator:        "creator",
	CreateTime:     "create_time",
	UpdateTime:     "update_time",
	ExpireTime:     "expire_time",
}

// NewSysFileChunkUploadDao creates and returns a new DAO object for table data access.
func NewSysFileChunkUploadDao(handlers ...gdb.ModelHandler) *SysFileChunkUploadDao {
	return &SysFileChunkUploadDao{
		group:    "default",
		table:    "sys_file_chunk_upload",
		columns:  sysFileChunkUploadColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysFileChunkUploadDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysFileChunkUploadDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysFileChunkUploadDao) Columns() SysFileChunkUploadColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysFileChunkUploadDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysFileChunkUploadDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysFileChunkUploadDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
