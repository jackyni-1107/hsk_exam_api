package v1

import "github.com/gogf/gf/v2/frame/g"

type AuditLogListReq struct {
	g.Meta    `path:"/audit-log/list" method:"get" tags:"审计日志" summary:"操作审计列表"`
	Page      int    `json:"page" dc:"页码"`
	Size      int    `json:"size" dc:"每页条数"`
	Username  string `json:"username" dc:"操作用户名"`
	Path      string `json:"path" dc:"请求路径"`
	Action    string `json:"action" dc:"操作动作"`
	LogType   string `json:"log_type" dc:"日志类型"`
	TraceId   string `json:"trace_id" dc:"链路追踪ID"`
	StartTime string `json:"start_time" dc:"开始时间"`
	EndTime   string `json:"end_time" dc:"结束时间"`
}

type AuditLogListRes struct {
	List  []*AuditLogItem `json:"list" dc:"列表"`
	Total int             `json:"total" dc:"总数"`
}

type AuditLogItem struct {
	Id           int64  `json:"id" dc:"记录ID"`
	UserId       int64  `json:"user_id" dc:"用户ID"`
	Username     string `json:"username" dc:"用户名"`
	UserType     int    `json:"user_type" dc:"用户类型"`
	Module       string `json:"module" dc:"模块"`
	Action       string `json:"action" dc:"操作动作"`
	LogType      string `json:"log_type" dc:"日志类型"`
	Method       string `json:"method" dc:"请求方法"`
	Path         string `json:"path" dc:"请求路径"`
	RequestData  string `json:"request_data" dc:"请求数据"`
	ResponseData string `json:"response_data" dc:"响应数据"`
	Ip           string `json:"ip" dc:"IP地址"`
	UserAgent    string `json:"user_agent" dc:"用户代理"`
	TraceId      string `json:"trace_id" dc:"链路追踪ID"`
	DeviceInfo   string `json:"device_info" dc:"设备信息"`
	DurationMs   int    `json:"duration_ms" dc:"耗时(毫秒)"`
	CreateTime   string `json:"create_time" dc:"创建时间"`
}

type AuditLogChangeDetailsReq struct {
	g.Meta `path:"/audit-log/{id}/change-details" method:"get" tags:"审计日志" summary:"变更明细"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"记录ID"`
}

type AuditLogChangeDetailsRes struct {
	List []*AuditChangeDetailItem `json:"list" dc:"变更明细列表"`
}

type AuditChangeDetailItem struct {
	Id          int64  `json:"id" dc:"记录ID"`
	TableName   string `json:"table_name" dc:"表名"`
	RecordId    int64  `json:"record_id" dc:"关联记录ID"`
	FieldName   string `json:"field_name" dc:"字段名"`
	BeforeValue string `json:"before_value" dc:"变更前值"`
	AfterValue  string `json:"after_value" dc:"变更后值"`
	CreateTime  string `json:"create_time" dc:"创建时间"`
}
