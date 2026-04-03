package audit_log

import (
	"context"

	v1 "exam/api/admin/audit_log/v1"
)

type IAuditLog interface {
	AuditLogList(ctx context.Context, req *v1.AuditLogListReq) (res *v1.AuditLogListRes, err error)
	AuditLogChangeDetails(ctx context.Context, req *v1.AuditLogChangeDetailsReq) (res *v1.AuditLogChangeDetailsRes, err error)
}
