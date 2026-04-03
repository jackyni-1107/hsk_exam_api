package v1

import "github.com/gogf/gf/v2/frame/g"

type SecurityEventLogListReq struct {
	g.Meta    `path:"/security-event-log/list" method:"get" tags:"安全事件" summary:"安全事件列表"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	EventType string `json:"event_type"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type SecurityEventLogListRes struct {
	List  []*SecurityEventLogItem `json:"list"`
	Total int                     `json:"total"`
}

type SecurityEventLogItem struct {
	Id         int64  `json:"id"`
	EventType  string `json:"event_type"`
	UserId     int64  `json:"user_id"`
	Ip         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Detail     string `json:"detail"`
	TraceId    string `json:"trace_id"`
	CreateTime string `json:"create_time"`
}
