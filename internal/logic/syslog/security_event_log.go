package syslog

import (
	"context"
	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sSysLog) SecurityEventLogList(ctx context.Context, page, size int, eventType, startTime, endTime string) ([]sysentity.SysSecurityEventLog, int, error) {
	page, size = s.getPageSize(page, size)

	m := sysdao.SysSecurityEventLog.Ctx(ctx)

	// 使用 g.Map 自动过滤空值查询
	m = m.Where(g.Map{
		"event_type": eventType,
	})

	// 时间范围过滤
	if startTime != "" {
		m = m.WhereGTE("create_time", startTime)
	}
	if endTime != "" {
		m = m.WhereLTE("create_time", endTime)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取安全事件日志总数失败")
	}

	var list []sysentity.SysSecurityEventLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取安全事件日志列表失败")
	}

	return list, total, nil
}
