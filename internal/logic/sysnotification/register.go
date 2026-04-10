package sysnotification

import svc "exam/internal/service/sysnotification"

func init() {
	svc.RegisterSysnotification(new(sSysnotification))
}
