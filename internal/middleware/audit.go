package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysdo "exam/internal/model/do/sys"
	"exam/internal/utility"
)

// 操作类型
const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
	actionQuery  = "query"
)

// 日志类型
const (
	logTypeOperation = "operation"
	logTypeAPIAccess = "api_access"
)

type operationLogIdKeyType string

const operationLogIdKey operationLogIdKeyType = "operation_log_id"

// SetOperationLogId 将 operation_log_id 写入 context
func SetOperationLogId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, operationLogIdKey, id)
}

// GetOperationLogId 从 context 获取 operation_log_id
func GetOperationLogId(ctx context.Context) int64 {
	if v := ctx.Value(operationLogIdKey); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// Audit 操作审计中间件：记录请求/响应到 Sys_operation_audit_log
func Audit(r *ghttp.Request) {
	start := time.Now()
	requestData := sanitizeRequest(r)
	ctxData := GetCtxData(r.GetCtx())
	userId := int64(consts.AnonymousUserId)
	username := ""
	userType := consts.UserTypeAdmin
	if ctxData != nil {
		userId = ctxData.UserId
		username = ctxData.Username
		userType = ctxData.UserType
	}
	action := inferAction(r)
	module := inferModule(r.URL.Path)
	method := r.Method
	path := r.URL.Path
	ip := r.GetClientIp()
	userAgent := r.Header.Get("User-Agent")
	logType := inferLogType(method)
	traceId := GetTraceId(r.GetCtx())
	deviceInfo := utility.ParseDeviceInfo(userAgent)

	// 同步插入以获取 ID，供 RecordChange 关联
	res, err := sysdao.SysOperationAuditLog.Ctx(r.GetCtx()).Insert(sysdo.SysOperationAuditLog{
		UserId:      userId,
		Username:    username,
		UserType:    userType,
		Module:      module,
		Action:      action,
		LogType:     logType,
		Method:      method,
		Path:        path,
		RequestData: requestData,
		Ip:          ip,
		UserAgent:   userAgent,
		TraceId:     traceId,
		DeviceInfo:  deviceInfo,
	})
	if err == nil && res != nil {
		if id, e := res.LastInsertId(); e == nil && id > 0 {
			r.SetCtx(SetOperationLogId(r.GetCtx(), id))
		}
	}

	r.Middleware.Next()

	// 异步更新 response 和 duration
	responseData := getResponseData(r)
	durationMs := int(time.Since(start).Milliseconds())
	opLogId := GetOperationLogId(r.GetCtx())
	go func() {
		if opLogId > 0 {
			_, _ = sysdao.SysOperationAuditLog.Ctx(context.Background()).Where("id", opLogId).Data(sysdo.SysOperationAuditLog{
				ResponseData: responseData,
				DurationMs:   durationMs,
			}).Update()
		}
	}()
}

func inferLogType(method string) string {
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return logTypeOperation
	default:
		return logTypeAPIAccess
	}
}

func inferAction(r *ghttp.Request) string {
	// 可从 g.Meta 扩展获取，此处按 method 推断
	switch r.Method {
	case "POST":
		return actionCreate
	case "PUT", "PATCH":
		return actionUpdate
	case "DELETE":
		return actionDelete
	default:
		return actionQuery
	}
}

func inferModule(path string) string {
	if route, ok := parseAdminRoute(path); ok {
		return route.module()
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func sanitizeRequest(r *ghttp.Request) string {
	// 优先读取 Body（POST/PUT/PATCH 等）
	body := r.GetBodyString()
	if body != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(body), &m); err != nil {
			return truncate(body, 4096)
		}
		if _, ok := m["password"]; ok {
			m["password"] = "***"
		}
		if _, ok := m["Password"]; ok {
			m["Password"] = "***"
		}
		data, _ := json.Marshal(m)
		return truncate(string(data), 4096)
	}
	// GET 等无 Body 时，记录 Query 参数
	if q := r.URL.RawQuery; q != "" {
		params := make(map[string]interface{})
		for k, v := range r.URL.Query() {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}
		if _, ok := params["password"]; ok {
			params["password"] = "***"
		}
		data, _ := json.Marshal(params)
		return truncate(string(data), 4096)
	}
	return ""
}

func getResponseData(r *ghttp.Request) string {
	res := r.GetHandlerResponse()
	if res == nil {
		return ""
	}
	data, err := json.Marshal(res)
	if err != nil {
		return ""
	}
	return truncate(string(data), 4096)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
