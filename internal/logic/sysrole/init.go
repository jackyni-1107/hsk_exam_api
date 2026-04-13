package sysrole

import (
	"exam/internal/service/sysrole"
)

type sSysRole struct{}

func init() {
	sysrole.RegisterSysRole(New())
}

func New() *sSysRole {
	return &sSysRole{}
}
