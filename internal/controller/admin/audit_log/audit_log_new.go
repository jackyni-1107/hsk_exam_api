package audit_log

import (
	"exam/api/admin/audit_log"
)

type ControllerV1 struct{}

func NewV1() audit_log.IAuditLog {
	return &ControllerV1{}
}
