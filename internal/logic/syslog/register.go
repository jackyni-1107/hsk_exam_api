package syslog

import svc "exam/internal/service/syslog"

func init() {
	svc.RegisterSysLog(new(sSysLog))
}
