package config

import (
	"context"

	v1 "exam/api/admin/config/v1"
)

type IConfig interface {
	ConfigList(ctx context.Context, req *v1.ConfigListReq) (res *v1.ConfigListRes, err error)
	ConfigCreate(ctx context.Context, req *v1.ConfigCreateReq) (res *v1.ConfigCreateRes, err error)
	ConfigUpdate(ctx context.Context, req *v1.ConfigUpdateReq) (res *v1.ConfigUpdateRes, err error)
	ConfigDelete(ctx context.Context, req *v1.ConfigDeleteReq) (res *v1.ConfigDeleteRes, err error)
	ConfigGet(ctx context.Context, req *v1.ConfigGetReq) (res *v1.ConfigGetRes, err error)
	DictTypeList(ctx context.Context, req *v1.DictTypeListReq) (res *v1.DictTypeListRes, err error)
	DictTypeCreate(ctx context.Context, req *v1.DictTypeCreateReq) (res *v1.DictTypeCreateRes, err error)
	DictTypeUpdate(ctx context.Context, req *v1.DictTypeUpdateReq) (res *v1.DictTypeUpdateRes, err error)
	DictTypeDelete(ctx context.Context, req *v1.DictTypeDeleteReq) (res *v1.DictTypeDeleteRes, err error)
	DictDataList(ctx context.Context, req *v1.DictDataListReq) (res *v1.DictDataListRes, err error)
	DictDataCreate(ctx context.Context, req *v1.DictDataCreateReq) (res *v1.DictDataCreateRes, err error)
	DictDataUpdate(ctx context.Context, req *v1.DictDataUpdateReq) (res *v1.DictDataUpdateRes, err error)
	DictDataDelete(ctx context.Context, req *v1.DictDataDeleteReq) (res *v1.DictDataDeleteRes, err error)
}
