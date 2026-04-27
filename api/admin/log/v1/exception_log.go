package v1

import "github.com/gogf/gf/v2/frame/g"

type ExceptionLogListReq struct {
	g.Meta    `path:"/log/exception-log/list" method:"get" tags:"异常日志" summary:"异常日志列表" permisson:"auto"`
	Page      int    `json:"page" dc:"页码"`
	Size      int    `json:"size" dc:"每页条数"`
	TraceId   string `json:"trace_id" dc:"链路追踪ID"`
	Path      string `json:"path" dc:"请求路径"`
	StartTime string `json:"start_time" dc:"开始时间"`
	EndTime   string `json:"end_time" dc:"结束时间"`
}

type ExceptionLogListRes struct {
	List  []*ExceptionLogItem `json:"list" dc:"列表"`
	Total int                 `json:"total" dc:"总数"`
}

type ExceptionLogItem struct {
	Id         int64  `json:"id" dc:"记录ID"`
	TraceId    string `json:"trace_id" dc:"链路追踪ID"`
	Path       string `json:"path" dc:"请求路径"`
	Method     string `json:"method" dc:"请求方法"`
	ErrorMsg   string `json:"error_msg" dc:"错误信息"`
	Stack      string `json:"stack" dc:"堆栈信息"`
	UserId     int64  `json:"user_id" dc:"用户ID"`
	Ip         string `json:"ip" dc:"IP地址"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}
