package sysuser

import "exam/internal/service/user"

type sSysUser struct{}

func init() {
	user.RegisterUser(New())
}

func New() *sSysUser {
	return &sSysUser{}
}
