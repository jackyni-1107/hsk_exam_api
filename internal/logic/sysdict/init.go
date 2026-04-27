package sysdict

import "exam/internal/service/sysdict"

type sSysDict struct{}

func init() {
	sysdict.RegisterSysDict(New())
}

func New() *sSysDict {
	return &sSysDict{}
}
