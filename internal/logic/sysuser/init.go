package sysuser

import "exam/internal/service/sysuser"

type sSysUser struct{}

func init() {
	sysuser.RegisterSysUser(New())
}

func New() *sSysUser {
	return &sSysUser{}
}
