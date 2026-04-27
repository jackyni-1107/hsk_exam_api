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
}
