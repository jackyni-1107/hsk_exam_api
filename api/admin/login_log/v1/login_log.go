package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginLogListReq struct {
	g.Meta    `path:"/login-log/list" method:"get" tags:"登录日志" summary:"登录日志列表" permission:"login_log:list"`
	Page      int    `json:"page" dc:"页码"`
	Size      int    `json:"size" dc:"每页条数"`
	Username  string `json:"username" dc:"用户名"`
	LogType   string `json:"log_type" dc:"日志类型"`
	UserType  int    `json:"user_type" dc:"用户类型"`
	StartTime string `json:"start_time" dc:"开始时间"`
	EndTime   string `json:"end_time" dc:"结束时间"`
}

type LoginLogListRes struct {
	List  []*LoginLogItem `json:"list" dc:"列表"`
	Total int             `json:"total" dc:"总数"`
}

type LoginLogItem struct {
	Id         int64  `json:"id" dc:"记录ID"`
	LogType    string `json:"log_type" dc:"日志类型"`
	UserId     int64  `json:"user_id" dc:"用户ID"`
	Username   string `json:"username" dc:"用户名"`
	UserType   int    `json:"user_type" dc:"用户类型"`
	Ip         string `json:"ip" dc:"IP地址"`
	UserAgent  string `json:"user_agent" dc:"用户代理"`
	DeviceInfo string `json:"device_info" dc:"设备信息"`
	TraceId    string `json:"trace_id" dc:"链路追踪ID"`
	FailReason string `json:"fail_reason" dc:"失败原因"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}
