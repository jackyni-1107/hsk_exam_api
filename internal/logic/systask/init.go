package systask

import "exam/internal/service/systask"

type sSysTask struct{}

func init() {
	systask.RegisterSysTask(New())
}

func New() *sSysTask {
	return &sSysTask{}
}
