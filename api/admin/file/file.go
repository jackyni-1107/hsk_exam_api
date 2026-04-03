package file

import (
	"context"

	v1 "exam/api/admin/file/v1"
)

type IFile interface {
	List(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error)
	Delete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error)
	StorageConfigList(ctx context.Context, req *v1.StorageConfigListReq) (res *v1.StorageConfigListRes, err error)
	StorageConfigCreate(ctx context.Context, req *v1.StorageConfigCreateReq) (res *v1.StorageConfigCreateRes, err error)
	StorageConfigUpdate(ctx context.Context, req *v1.StorageConfigUpdateReq) (res *v1.StorageConfigUpdateRes, err error)
	StorageConfigDelete(ctx context.Context, req *v1.StorageConfigDeleteReq) (res *v1.StorageConfigDeleteRes, err error)
	StorageConfigSetActive(ctx context.Context, req *v1.StorageConfigSetActiveReq) (res *v1.StorageConfigSetActiveRes, err error)
}
