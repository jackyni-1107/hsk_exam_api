// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package sysfile

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
	"io"
)

type (
	ISysfile interface {
		FileList(ctx context.Context, page int, size int, filename string) ([]sysentity.SysFileStorage, int, error)
		// FileUpload 将文件写入当前活动存储并登记 sys_file_storage，返回新记录 id、存储 path、展示用文件名。
		FileUpload(ctx context.Context, originalFilename string, size int64, contentType string, body io.ReadSeeker, isPrivate int, creator string) (id int64, objectPath string, displayName string, err error)
		// FileOpenDownload 按文件记录打开存储对象，供 HTTP 下载流式写出。
		FileOpenDownload(ctx context.Context, id int64) (filename string, mime string, size int64, body io.ReadCloser, err error)
		FileDelete(ctx context.Context, id int64) error
		StorageConfigList(ctx context.Context) ([]sysentity.SysFileStorageConfig, error)
		StorageConfigCreate(ctx context.Context, storageType string, name string, configJson string, creator string, cleanupBeforeDays int) (int64, error)
		StorageConfigUpdate(ctx context.Context, id int64, name string, configJson string, updater string, cleanupBeforeDays int) error
		StorageConfigDelete(ctx context.Context, id int64, updater string) error
		StorageConfigSetActive(ctx context.Context, id int64, updater string) error
	}
)

var (
	localSysfile ISysfile
)

func Sysfile() ISysfile {
	if localSysfile == nil {
		panic("implement not found for interface ISysfile, forgot register?")
	}
	return localSysfile
}

func RegisterSysfile(i ISysfile) {
	localSysfile = i
}
