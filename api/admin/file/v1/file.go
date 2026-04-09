package v1

import "github.com/gogf/gf/v2/frame/g"

type FileListReq struct {
	g.Meta   `path:"/file/list" method:"get" tags:"文件" summary:"文件列表"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Filename string `json:"filename"`
}

type FileListRes struct {
	List  []*FileItem `json:"list"`
	Total int         `json:"total"`
}

type FileItem struct {
	Id         int64  `json:"id"`
	Filename   string `json:"filename"`
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	MimeType   string `json:"mime_type"`
	IsPrivate  int    `json:"is_private"`
	CreateTime string `json:"create_time"`
}

type FileDeleteReq struct {
	g.Meta `path:"/file/{id}" method:"delete" tags:"文件" summary:"删除文件"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type FileDeleteRes struct{}

type StorageConfigListReq struct {
	g.Meta `path:"/file/storage-config/list" method:"get" tags:"文件存储" summary:"存储配置列表"`
}

type StorageConfigListRes struct {
	List []*StorageConfigItem `json:"list"`
}

type StorageConfigItem struct {
	Id                int64  `json:"id"`
	StorageType       string `json:"storage_type"`
	Name              string `json:"name"`
	IsActive          int    `json:"is_active"`
	ConfigJson        string `json:"config_json"`
	CleanupBeforeDays int    `json:"cleanup_before_days"`
	CreateTime        string `json:"create_time"`
}

type StorageConfigCreateReq struct {
	g.Meta            `path:"/file/storage-config" method:"post" tags:"文件存储" summary:"新增存储配置"`
	StorageType       string `json:"storage_type" v:"required#err.invalid_params"`
	Name              string `json:"name" v:"required#err.invalid_params"`
	ConfigJson        string `json:"config_json" v:"required#err.invalid_params"`
	CleanupBeforeDays int    `json:"cleanup_before_days"`
}

type StorageConfigCreateRes struct {
	Id int64 `json:"id"`
}

type StorageConfigUpdateReq struct {
	g.Meta            `path:"/file/storage-config/{id}" method:"put" tags:"文件存储" summary:"更新存储配置"`
	Id                int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	Name              string `json:"name"`
	ConfigJson        string `json:"config_json"`
	CleanupBeforeDays int    `json:"cleanup_before_days"`
}

type StorageConfigUpdateRes struct{}

type StorageConfigDeleteReq struct {
	g.Meta `path:"/file/storage-config/{id}" method:"delete" tags:"文件存储" summary:"删除存储配置"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type StorageConfigDeleteRes struct{}

type StorageConfigSetActiveReq struct {
	g.Meta `path:"/file/storage-config/{id}/set-active" method:"post" tags:"文件存储" summary:"设为当前存储"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type StorageConfigSetActiveRes struct{}
