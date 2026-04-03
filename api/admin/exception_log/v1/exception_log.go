package v1

import "github.com/gogf/gf/v2/frame/g"

type ExceptionLogListReq struct {
	g.Meta    `path:"/exception-log/list" method:"get" tags:"异常日志" summary:"异常日志列表"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	TraceId   string `json:"trace_id"`
	Path      string `json:"path"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ExceptionLogListRes struct {
	List  []*ExceptionLogItem `json:"list"`
	Total int                 `json:"total"`
}

type ExceptionLogItem struct {
	Id         int64  `json:"id"`
	TraceId    string `json:"trace_id"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	ErrorMsg   string `json:"error_msg"`
	Stack      string `json:"stack"`
	UserId     int64  `json:"user_id"`
	Ip         string `json:"ip"`
	CreateTime string `json:"create_time"`
}
