package security_event_log

import (
	"context"

	"exam/api/admin/security_event_log/v1"
)

type ISecurityEventLog interface {
	ISecurityEventLogList
}

type ISecurityEventLogList interface {
	SecurityEventLogList(ctx context.Context, req *v1.SecurityEventLogListReq) (res *v1.SecurityEventLogListRes, err error)
}
