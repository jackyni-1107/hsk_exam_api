package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginLogListReq struct {
	g.Meta    `path:"/login-log/list" method:"get" tags:"登录日志" summary:"登录日志列表"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	Username  string `json:"username"`
	LogType   string `json:"log_type"`
	UserType  int    `json:"user_type"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type LoginLogListRes struct {
	List  []*LoginLogItem `json:"list"`
	Total int             `json:"total"`
}

type LoginLogItem struct {
	Id         int64  `json:"id"`
	LogType    string `json:"log_type"`
	UserId     int64  `json:"user_id"`
	Username   string `json:"username"`
	UserType   int    `json:"user_type"`
	Ip         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	DeviceInfo string `json:"device_info"`
	TraceId    string `json:"trace_id"`
	FailReason string `json:"fail_reason"`
	CreateTime string `json:"create_time"`
}
