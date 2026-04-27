// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package audit

import (
	"context"
	"exam/internal/model/bo"
)

type (
	IAudit interface {
		// RecordLoginSuccess 记录登录成功
		RecordLoginSuccess(ctx context.Context, userId int64, username string, userType int, ip string, userAgent string, traceId string)
		// RecordLoginFail 记录登录失败
		RecordLoginFailure(ctx context.Context, userId int64, username string, userType int, ip string, userAgent string, failReason string, traceId string)
		// RecordLogout 记录登出
		RecordLogout(ctx context.Context, userId int64, username string, userType int, ip string, userAgent string, traceId string)
		// RecordSecurityEvent 记录安全事件
		RecordSecurityEvent(ctx context.Context, eventType string, userId int64, ip string, userAgent string, detail string, traceId string)
		// RecordException 记录异常日志
		RecordException(ctx context.Context, path string, method string, errorMsg string, stack string, userId int64, ip string, traceId string)
		CreateOperationLog(ctx context.Context, in bo.OperationAuditLogCreateInput) (int64, error)
		FinishOperationLog(ctx context.Context, operationLogId int64, responseData string, durationMs int) error
		// RecordChange 记录变更明细，过滤敏感字段；operationLogId 为 0 时跳过
		RecordChange(ctx context.Context, tableName string, recordId int64, operationLogId int64, beforeMap map[string]interface{}, afterMap map[string]interface{})
	}
)

var (
	localAudit IAudit
)

func Audit() IAudit {
	if localAudit == nil {
		panic("implement not found for interface IAudit, forgot register?")
	}
	return localAudit
}

func RegisterAudit(i IAudit) {
	localAudit = i
}
