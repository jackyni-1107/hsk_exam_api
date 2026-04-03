package login_log

import (
	"exam/api/admin/login_log"
)

type ControllerV1 struct{}

func NewV1() login_log.ILoginLog {
	return &ControllerV1{}
}
