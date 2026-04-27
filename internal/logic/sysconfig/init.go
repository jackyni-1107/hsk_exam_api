package sysconfig

import "exam/internal/service/sysconfig"

type sSysConfig struct{}

func init() {
	sysconfig.RegisterSysConfig(New())
}

func New() *sSysConfig {
	return &sSysConfig{}
}
