package log

import (
	"context"
	"exam/api/admin/log/v1"
)

type ISecurityEventLog interface {
	ISecurityEventLogList
}

type ISecurityEventLogList interface {
	SecurityEventLogList(ctx context.Context, req *v1.SecurityEventLogListReq) (res *v1.SecurityEventLogListRes, err error)
}
