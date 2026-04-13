package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FileListReq struct {
	g.Meta   `path:"/file/list" method:"get" tags:"文件" summary:"文件列表"`
	Page     int    `json:"page" dc:"页码"`
	Size     int    `json:"size" dc:"每页条数"`
	Filename string `json:"filename" dc:"文件名"`
}

type FileListRes struct {
	List  []*FileItem `json:"list" dc:"列表"`
	Total int         `json:"total" dc:"总数"`
}

type FileItem struct {
	Id         int64  `json:"id" dc:"文件ID"`
	Filename   string `json:"filename" dc:"文件名"`
	Path       string `json:"path" dc:"存储路径"`
	Size       int64  `json:"size" dc:"文件大小(字节)"`
	MimeType   string `json:"mime_type" dc:"MIME类型"`
	IsPrivate  int    `json:"is_private" dc:"是否私有：0公开 1私有"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type FileDeleteReq struct {
	g.Meta `path:"/file/{id}" method:"delete" tags:"文件" summary:"删除文件"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"文件ID"`
}

type FileDeleteRes struct{}

type FileUploadReq struct {
	g.Meta    `path:"/file/upload" method:"post" tags:"文件" summary:"上传文件"`
	IsPrivate int               `json:"is_private" form:"is_private" d:"0" dc:"是否私有：0公开 1私有"`
	File      *ghttp.UploadFile `json:"file" type:"file" dc:"文件（表单字段名 file）"`
}

type FileUploadRes struct {
	Id       int64  `json:"id" dc:"文件ID"`
	Path     string `json:"path" dc:"存储路径"`
	Filename string `json:"filename" dc:"原始文件名"`
	Size     int64  `json:"size" dc:"大小(字节)"`
	MimeType string `json:"mime_type" dc:"MIME"`
}

type StorageConfigListReq struct {
	g.Meta `path:"/file/storage-config/list" method:"get" tags:"文件存储" summary:"存储配置列表"`
}

type StorageConfigListRes struct {
	List []*StorageConfigItem `json:"list" dc:"存储配置列表"`
}

type StorageConfigItem struct {
	Id                int64  `json:"id" dc:"配置ID"`
	StorageType       string `json:"storage_type" dc:"存储类型"`
	Name              string `json:"name" dc:"配置名称"`
	IsActive          int    `json:"is_active" dc:"是否当前使用：0否 1是"`
	ConfigJson        string `json:"config_json" dc:"配置JSON"`
	CleanupBeforeDays int    `json:"cleanup_before_days" dc:"自动清理天数"`
	CreateTime        string `json:"create_time" dc:"创建时间"`
}

type StorageConfigCreateReq struct {
	g.Meta            `path:"/file/storage-config" method:"post" tags:"文件存储" summary:"新增存储配置"`
	StorageType       string `json:"storage_type" v:"required#err.invalid_params" dc:"存储类型"`
	Name              string `json:"name" v:"required#err.invalid_params" dc:"配置名称"`
	ConfigJson        string `json:"config_json" v:"required#err.invalid_params" dc:"配置JSON"`
	CleanupBeforeDays int    `json:"cleanup_before_days" dc:"自动清理天数"`
}

type StorageConfigCreateRes struct {
	Id int64 `json:"id" dc:"配置ID"`
}

type StorageConfigUpdateReq struct {
	g.Meta            `path:"/file/storage-config/{id}" method:"put" tags:"文件存储" summary:"更新存储配置"`
	Id                int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
	Name              string `json:"name" dc:"配置名称"`
	ConfigJson        string `json:"config_json" dc:"配置JSON"`
	CleanupBeforeDays int    `json:"cleanup_before_days" dc:"自动清理天数"`
}

type StorageConfigUpdateRes struct{}

type StorageConfigDeleteReq struct {
	g.Meta `path:"/file/storage-config/{id}" method:"delete" tags:"文件存储" summary:"删除存储配置"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
}

type StorageConfigDeleteRes struct{}

type StorageConfigSetActiveReq struct {
	g.Meta `path:"/file/storage-config/{id}/set-active" method:"post" tags:"文件存储" summary:"设为当前存储"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
}

type StorageConfigSetActiveRes struct{}
