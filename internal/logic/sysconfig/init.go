package sysconfig

import (
	"exam/internal/service/sysconfig"

	"github.com/gogf/gf/v2/os/gcache"
)

type sSysConfig struct {
	cache *gcache.Cache // 预留缓存对象
}

func init() {
	sysconfig.RegisterSysConfig(New())
}

func New() *sSysConfig {
	return &sSysConfig{
		cache: gcache.New(), // 初始化缓存
	}
}
