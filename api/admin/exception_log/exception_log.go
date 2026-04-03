package exception_log

import (
	"context"

	"exam/api/admin/exception_log/v1"
)

type IExceptionLog interface {
	IExceptionLogList
}

type IExceptionLogList interface {
	ExceptionLogList(ctx context.Context, req *v1.ExceptionLogListReq) (res *v1.ExceptionLogListRes, err error)
}
