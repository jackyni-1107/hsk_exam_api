package syslog

import (
	"context"
	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysLog) LoginLogList(ctx context.Context, page, size int, username, logType string, userType int, startTime, endTime string) ([]sysentity.SysLoginLog, int, error) {
	page, size = s.getPageSize(page, size)
	m := sysdao.SysLoginLog.Ctx(ctx)

	if username != "" {
		m = m.WhereLike("username", "%"+username+"%")
	}

	if logType != "" {
		m = m.Where("log_type", logType)
	}
	if userType > 0 {
		m = m.Where("user_type", userType)
	}

	if startTime != "" {
		m = m.WhereGTE("create_time", startTime)
	}
	if endTime != "" {
		m = m.WhereLTE("create_time", endTime)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取登录日志总数失败")
	}

	var list []sysentity.SysLoginLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取登录日志列表失败")
	}

	return list, total, nil
}
