package sysfile

import svc "exam/internal/service/sysfile"

func init() {
	svc.RegisterSysfile(new(sSysfile))
}
