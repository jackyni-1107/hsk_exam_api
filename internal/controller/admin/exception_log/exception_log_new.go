package exception_log

import (
	"exam/api/admin/exception_log"
)

type ControllerV1 struct{}

func NewV1() exception_log.IExceptionLog {
	return &ControllerV1{}
}
