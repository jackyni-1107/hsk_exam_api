// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package syslog

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysLog interface {
		AuditLogList(ctx context.Context, page int, size int, username string, path string, action string, logType string, traceId string, startTime string, endTime string) ([]sysentity.SysOperationAuditLog, int, error)
		AuditLogChangeDetails(ctx context.Context, operationLogId int64) ([]sysentity.SysAuditChangeDetail, error)
		ExceptionLogList(ctx context.Context, page int, size int, traceId string, path string, startTime string, endTime string) ([]sysentity.SysExceptionLog, int, error)
		LoginLogList(ctx context.Context, page int, size int, username string, logType string, userType int, startTime string, endTime string) ([]sysentity.SysLoginLog, int, error)
		SecurityEventLogList(ctx context.Context, page int, size int, eventType string, startTime string, endTime string) ([]sysentity.SysSecurityEventLog, int, error)
	}
)

var (
	localSysLog ISysLog
)

func SysLog() ISysLog {
	if localSysLog == nil {
		panic("implement not found for interface ISysLog, forgot register?")
	}
	return localSysLog
}

func RegisterSysLog(i ISysLog) {
	localSysLog = i
}
