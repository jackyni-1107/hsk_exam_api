package sysnotification

import (
	"exam/internal/service/sysnotification"
)

type sSysnotification struct{}

func init() {
	sysnotification.RegisterSysnotification(New())
}

func New() *sSysnotification {
	return &sSysnotification{}
}
