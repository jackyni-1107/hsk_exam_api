package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
	"exam/internal/model/bo"
	auditsvc "exam/internal/service/audit"
	"exam/internal/utility"
)

const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
	actionQuery  = "query"
)

const (
	logTypeOperation = "operation"
	logTypeAPIAccess = "api_access"
)

type operationLogIdKeyType string

const operationLogIdKey operationLogIdKeyType = "operation_log_id"

func SetOperationLogId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, operationLogIdKey, id)
}

func GetOperationLogId(ctx context.Context) int64 {
	if v := ctx.Value(operationLogIdKey); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func Audit(r *ghttp.Request) {
	start := time.Now()
	requestData := sanitizeRequest(r)
	ctxData := GetCtxData(r.GetCtx())
	userID := int64(consts.AnonymousUserId)
	username := ""
	userType := consts.UserTypeAdmin
	if ctxData != nil {
		userID = ctxData.UserId
		username = ctxData.Username
		userType = ctxData.UserType
	}
	method := r.Method
	opLogID, err := auditsvc.Audit().CreateOperationLog(r.GetCtx(), bo.OperationAuditLogCreateInput{
		UserId:      userID,
		Username:    username,
		UserType:    userType,
		Module:      inferModule(r.URL.Path),
		Action:      inferAction(r),
		LogType:     inferLogType(method),
		Method:      method,
		Path:        r.URL.Path,
		RequestData: requestData,
		Ip:          r.GetClientIp(),
		UserAgent:   r.Header.Get("User-Agent"),
		TraceId:     GetTraceId(r.GetCtx()),
		DeviceInfo:  utility.ParseDeviceInfo(r.Header.Get("User-Agent")),
	})
	if err == nil && opLogID > 0 {
		r.SetCtx(SetOperationLogId(r.GetCtx(), opLogID))
	}

	r.Middleware.Next()

	responseData := getResponseData(r)
	durationMs := int(time.Since(start).Milliseconds())
	opLogID = GetOperationLogId(r.GetCtx())
	go func() {
		if opLogID > 0 {
			_ = auditsvc.Audit().FinishOperationLog(context.Background(), opLogID, responseData, durationMs)
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
