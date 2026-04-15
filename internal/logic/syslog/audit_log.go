package syslog

import (
	"context"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysLog) AuditLogList(ctx context.Context, page, size int, username, path, action, logType, traceId, startTime, endTime string) ([]sysentity.SysOperationAuditLog, int, error) {
	page, size = s.getPageSize(page, size)

	// 构建基础模型
	m := dao.SystemOperationAuditLog.Ctx(ctx)

	// 模糊查询
	if username != "" {
		m = m.WhereLike("username", "%"+username+"%")
	}
	if path != "" {
		m = m.WhereLike("path", "%"+path+"%")
	}

	if action != "" {
		m = m.Where("action", action)
	}
	if logType != "" {
		m = m.Where("log_type", logType)
	}
	if traceId != "" {
		m = m.Where("trace_id", traceId)
	}

	// 时间范围
	if startTime != "" {
		m = m.WhereGTE("create_time", startTime)
	}
	if endTime != "" {
		m = m.WhereLTE("create_time", endTime)
	}

	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取审计日志总数失败")
	}

	var list []sysentity.SysOperationAuditLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "获取审计日志列表失败")
	}

	return list, total, nil
}

func (s *sSysLog) AuditLogChangeDetails(ctx context.Context, operationLogId int64) ([]sysentity.SysAuditChangeDetail, error) {
	var list []sysentity.SysAuditChangeDetail
	err := sysdao.SysAuditChangeDetail.Ctx(ctx).Where("operation_log_id", operationLogId).OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, nil
}
