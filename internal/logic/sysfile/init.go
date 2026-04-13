package sysfile

import (
	"exam/internal/service/sysfile"
)

type sSysFile struct{}

func init() {
	sysfile.RegisterSysFile(New())
}

func New() *sSysFile {
	return &sSysFile{}
}
