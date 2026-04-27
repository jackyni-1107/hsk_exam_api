package log

import apilog "exam/api/admin/log"

type ControllerV1 struct{}

func NewV1() apilog.ILog {
	return &ControllerV1{}
}
