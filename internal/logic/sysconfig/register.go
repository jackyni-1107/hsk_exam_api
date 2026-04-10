package sysconfig

import svc "exam/internal/service/sysconfig"

func init() {
	svc.RegisterSysconfig(new(sSysconfig))
}
