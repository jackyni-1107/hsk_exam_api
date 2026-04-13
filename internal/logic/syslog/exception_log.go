package syslog

import (
	"context"
	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sSysLog) ExceptionLogList(ctx context.Context, page, size int, traceId, path, startTime, endTime string) ([]sysentity.SysExceptionLog, int, error) {
	page, size = s.getPageSize(page, size)

	m := sysdao.SysExceptionLog.Ctx(ctx)

	// 模糊查询
	if path != "" {
		m = m.WhereLike("path", "%"+path+"%")
	}

	// 精确查询与范围查询
	m = m.Where(g.Map{
		"trace_id": traceId,
	})

	if startTime != "" {
		m = m.WhereGTE("create_time", startTime)
	}
	if endTime != "" {
		m = m.WhereLTE("create_time", endTime)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取异常日志总数失败")
	}

	var list []sysentity.SysExceptionLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取异常日志列表失败")
	}

	return list, total, nil
}
