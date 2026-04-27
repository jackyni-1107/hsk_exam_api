package v1

import "github.com/gogf/gf/v2/frame/g"

type SecurityEventLogListReq struct {
	g.Meta    `path:"/log/security-event-log/list" method:"get" tags:"安全事件" summary:"安全事件列表"`
	Page      int    `json:"page" dc:"页码"`
	Size      int    `json:"size" dc:"每页条数"`
	EventType string `json:"event_type" dc:"事件类型"`
	StartTime string `json:"start_time" dc:"开始时间"`
	EndTime   string `json:"end_time" dc:"结束时间"`
}

type SecurityEventLogListRes struct {
	List  []*SecurityEventLogItem `json:"list" dc:"列表"`
	Total int                     `json:"total" dc:"总数"`
}

type SecurityEventLogItem struct {
	Id         int64  `json:"id" dc:"记录ID"`
	EventType  string `json:"event_type" dc:"事件类型"`
	UserId     int64  `json:"user_id" dc:"用户ID"`
	Ip         string `json:"ip" dc:"IP地址"`
	UserAgent  string `json:"user_agent" dc:"用户代理"`
	Detail     string `json:"detail" dc:"事件详情"`
	TraceId    string `json:"trace_id" dc:"链路追踪ID"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}
