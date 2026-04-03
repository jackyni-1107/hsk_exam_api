package v1

import "github.com/gogf/gf/v2/frame/g"

type AuditLogListReq struct {
	g.Meta    `path:"/audit-log/list" method:"get" tags:"审计日志" summary:"操作审计列表"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	Username  string `json:"username"`
	Path      string `json:"path"`
	Action    string `json:"action"`
	LogType   string `json:"log_type"`
	TraceId   string `json:"trace_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type AuditLogListRes struct {
	List  []*AuditLogItem `json:"list"`
	Total int             `json:"total"`
}

type AuditLogItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	UserType     int    `json:"user_type"`
	Module       string `json:"module"`
	Action       string `json:"action"`
	LogType      string `json:"log_type"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	RequestData  string `json:"request_data"`
	ResponseData string `json:"response_data"`
	Ip           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
	TraceId      string `json:"trace_id"`
	DeviceInfo   string `json:"device_info"`
	DurationMs   int    `json:"duration_ms"`
	CreateTime   string `json:"create_time"`
}

type AuditLogChangeDetailsReq struct {
	g.Meta `path:"/audit-log/{id}/change-details" method:"get" tags:"审计日志" summary:"变更明细"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type AuditLogChangeDetailsRes struct {
	List []*AuditChangeDetailItem `json:"list"`
}

type AuditChangeDetailItem struct {
	Id          int64  `json:"id"`
	TableName   string `json:"table_name"`
	RecordId    int64  `json:"record_id"`
	FieldName   string `json:"field_name"`
	BeforeValue string `json:"before_value"`
	AfterValue  string `json:"after_value"`
	CreateTime  string `json:"create_time"`
}
