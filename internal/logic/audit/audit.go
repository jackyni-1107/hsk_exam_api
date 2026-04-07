package audit

import (
	"context"
	"exam/internal/consts"
	"fmt"

	sysdao "exam/internal/dao/sys"
	sysdo "exam/internal/model/do/sys"
	"exam/internal/utility"
)

type sAudit struct{}

func New() *sAudit {
	return &sAudit{}
}

// RecordLoginSuccess 记录登录成功
func (s *sAudit) RecordLoginSuccess(ctx context.Context, userId int64, username string, userType int, ip, userAgent, traceId string) {
	deviceInfo := utility.ParseDeviceInfo(userAgent)
	go func() {
		_, _ = sysdao.SysLoginLog.Ctx(context.Background()).Insert(sysdo.SysLoginLog{
			LogType:    consts.AuditLogTypeLoginSuccess,
			UserId:     userId,
			Username:   username,
			UserType:   userType,
			Ip:         ip,
			UserAgent:  userAgent,
			DeviceInfo: deviceInfo,
			TraceId:    traceId,
		})
	}()
}

// RecordLoginFail 记录登录失败
func (s *sAudit) RecordLoginFail(ctx context.Context, userId int64, username string, userType int, ip, userAgent, failReason, traceId string) {
	deviceInfo := utility.ParseDeviceInfo(userAgent)
	go func() {
		_, _ = sysdao.SysLoginLog.Ctx(context.Background()).Insert(sysdo.SysLoginLog{
			LogType:    consts.AuditLogTypeLoginFail,
			UserId:     userId,
			Username:   username,
			UserType:   userType,
			Ip:         ip,
			UserAgent:  userAgent,
			DeviceInfo: deviceInfo,
			TraceId:    traceId,
			FailReason: failReason,
		})
	}()
}

// RecordLogout 记录登出
func (s *sAudit) RecordLogout(ctx context.Context, userId int64, username string, userType int, ip, userAgent, traceId string) {
	deviceInfo := utility.ParseDeviceInfo(userAgent)
	go func() {
		_, _ = sysdao.SysLoginLog.Ctx(context.Background()).Insert(sysdo.SysLoginLog{
			LogType:    consts.AuditLogTypeLogout,
			UserId:     userId,
			Username:   username,
			UserType:   userType,
			Ip:         ip,
			UserAgent:  userAgent,
			DeviceInfo: deviceInfo,
			TraceId:    traceId,
		})
	}()
}

// RecordSecurityEvent 记录安全事件
func (s *sAudit) RecordSecurityEvent(ctx context.Context, eventType string, userId int64, ip, userAgent, detail, traceId string) {
	go func() {
		_, _ = sysdao.SysSecurityEventLog.Ctx(context.Background()).Insert(sysdo.SysSecurityEventLog{
			EventType: eventType,
			UserId:    userId,
			Ip:        ip,
			UserAgent: userAgent,
			Detail:    detail,
			TraceId:   traceId,
		})
	}()
}

// RecordException 记录异常日志
func (s *sAudit) RecordException(ctx context.Context, path, method, errorMsg, stack string, userId int64, ip, traceId string) {
	go func() {
		_, _ = sysdao.SysExceptionLog.Ctx(context.Background()).Insert(sysdo.SysExceptionLog{
			TraceId:  traceId,
			Path:     path,
			Method:   method,
			ErrorMsg: errorMsg,
			Stack:    stack,
			UserId:   userId,
			Ip:       ip,
		})
	}()
}

// 敏感字段（不记录变更前后值）
var sensitiveFields = map[string]bool{
	"password": true, "Password": true, "password_hash": true,
}

// RecordChange 记录变更明细，过滤敏感字段；operationLogId 为 0 时跳过
func (s *sAudit) RecordChange(ctx context.Context, tableName string, recordId int64, operationLogId int64, beforeMap, afterMap map[string]interface{}) {
	if operationLogId <= 0 {
		return
	}
	allKeys := make(map[string]bool)
	for k := range beforeMap {
		allKeys[k] = true
	}
	for k := range afterMap {
		allKeys[k] = true
	}
	var details []sysdo.SysAuditChangeDetail
	for k := range allKeys {
		if sensitiveFields[k] {
			continue
		}
		beforeVal := beforeMap[k]
		afterVal := afterMap[k]
		beforeStr := toStr(beforeVal)
		afterStr := toStr(afterVal)
		if beforeStr == afterStr {
			continue
		}
		details = append(details, sysdo.SysAuditChangeDetail{
			OperationLogId: operationLogId,
			TableName:      tableName,
			RecordId:       recordId,
			FieldName:      k,
			BeforeValue:    beforeStr,
			AfterValue:     afterStr,
		})
	}
	if len(details) == 0 {
		return
	}
	go func() {
		for _, d := range details {
			_, _ = sysdao.SysAuditChangeDetail.Ctx(context.Background()).Insert(d)
		}
	}()
}

func toStr(v interface{}) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", x)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", x)
	case float32, float64:
		return fmt.Sprintf("%v", x)
	case bool:
		return fmt.Sprintf("%t", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}
