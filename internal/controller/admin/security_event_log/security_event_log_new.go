package security_event_log

import (
	"exam/api/admin/security_event_log"
)

type ControllerV1 struct{}

func NewV1() security_event_log.ISecurityEventLog {
	return &ControllerV1{}
}
