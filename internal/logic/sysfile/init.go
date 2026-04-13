package sysfile

import (
	"exam/internal/service/sysfile"
)

type sSysfile struct{}

func init() {
	sysfile.RegisterSysfile(New())
}

func New() *sSysfile {
	return &sSysfile{}
}
