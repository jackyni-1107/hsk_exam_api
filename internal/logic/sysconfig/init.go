package sysconfig

import (
	"exam/internal/service/sysconfig"
	"github.com/gogf/gf/v2/os/gcache"
)

type sSysconfig struct {
	cache *gcache.Cache // 预留缓存对象
}

func init() {
	sysconfig.RegisterSysconfig(New())
}

func New() *sSysconfig {
	return &sSysconfig{
		cache: gcache.New(), // 初始化缓存
	}
}
