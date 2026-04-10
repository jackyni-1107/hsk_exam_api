package systask

import svc "exam/internal/service/systask"

func init() {
	svc.RegisterSysTask(new(sSysTask))
}
