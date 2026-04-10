package syslog

import (
	"context"

	sysentity "exam/internal/model/entity/sys"
)

type ISysLog interface {
	AuditLogList(ctx context.Context, page, size int, username, path, action, logType, traceId, startTime, endTime string) ([]sysentity.SysOperationAuditLog, int, error)
	AuditLogChangeDetails(ctx context.Context, operationLogId int64) ([]sysentity.SysAuditChangeDetail, error)
	LoginLogList(ctx context.Context, page, size int, username, logType string, userType int, startTime, endTime string) ([]sysentity.SysLoginLog, int, error)
	ExceptionLogList(ctx context.Context, page, size int, traceId, path, startTime, endTime string) ([]sysentity.SysExceptionLog, int, error)
	SecurityEventLogList(ctx context.Context, page, size int, eventType, startTime, endTime string) ([]sysentity.SysSecurityEventLog, int, error)
}

var localSysLog ISysLog

func SysLog() ISysLog {
	if localSysLog == nil {
		panic("implement not found for interface ISysLog, forgot register?")
	}
	return localSysLog
}

func RegisterSysLog(i ISysLog) {
	localSysLog = i
}
