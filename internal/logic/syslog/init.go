package syslog

import (
	"exam/internal/service/syslog"
)

type sSysLog struct{}

func init() {
	syslog.RegisterSysLog(New())
}

func New() *sSysLog {
	return &sSysLog{}
}

// 内部通用分页逻辑，减少各文件重复代码
func (s *sSysLog) getPageSize(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return page, size
}
