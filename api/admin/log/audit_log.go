package log

import (
	"context"
	"exam/api/admin/log/v1"
)

type ILog interface {
	IAuditLog
	IExceptionLog
	ILoginLog
	ISecurityEventLog
}

type IAuditLog interface {
	AuditLogList(ctx context.Context, req *v1.AuditLogListReq) (res *v1.AuditLogListRes, err error)
	AuditLogChangeDetails(ctx context.Context, req *v1.AuditLogChangeDetailsReq) (res *v1.AuditLogChangeDetailsRes, err error)
}
