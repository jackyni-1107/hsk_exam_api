package syslog

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
)

func (s *sSysLog) SecurityEventLogList(ctx context.Context, page, size int, eventType, startTime, endTime string) ([]sysentity.SysSecurityEventLog, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	m := sysdao.SysSecurityEventLog.Ctx(ctx)
	if eventType != "" {
		m = m.Where("event_type", eventType)
	}
	if startTime != "" {
		m = m.WhereGTE("create_time", startTime)
	}
	if endTime != "" {
		m = m.WhereLTE("create_time", endTime)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysSecurityEventLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}
