package log

import (
	"context"
	"exam/api/admin/log/v1"
)

type ILoginLog interface {
	ILoginLogList
}

type ILoginLogList interface {
	LoginLogList(ctx context.Context, req *v1.LoginLogListReq) (res *v1.LoginLogListRes, err error)
}
