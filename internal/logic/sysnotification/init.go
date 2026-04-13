package sysnotification

import "exam/internal/service/sysnotification"

type sSysNotification struct{}

func init() {
	sysnotification.RegisterSysNotification(New())
}

func New() *sSysNotification {
	return &sSysNotification{}
}
