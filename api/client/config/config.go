package config

import (
	"context"

	v1 "exam/api/client/config/v1"
)

type IConfig interface {
	MultiGet(ctx context.Context, req *v1.MultiGetReq) (res *v1.MultiGetRes, err error)
}
