package v1

import "github.com/gogf/gf/v2/frame/g"

type TaskListReq struct {
	g.Meta `path:"/task/list" method:"get" tags:"任务" summary:"任务列表"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Type   int    `json:"type"`
	Status *int   `json:"status"`
	Handler string `json:"handler"`
}

type TaskListRes struct {
	List  []*TaskItem `json:"list"`
	Total int         `json:"total"`
}

type TaskItem struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Type           int    `json:"type"`
	CronExpr       string `json:"cron_expr"`
	DelaySeconds   int    `json:"delay_seconds"`
	Handler        string `json:"handler"`
	Params         string `json:"params"`
	RetryTimes     int    `json:"retry_times"`
	RetryInterval  int    `json:"retry_interval"`
	Concurrency    int    `json:"concurrency"`
	AlertOnFail    int    `json:"alert_on_fail"`
	AlertReceivers string `json:"alert_receivers"`
	Status         int    `json:"status"`
	Remark         string `json:"remark"`
	CreateTime     string `json:"create_time"`
}

type TaskCreateReq struct {
	g.Meta         `path:"/task" method:"post" tags:"任务" summary:"创建任务"`
	Name           string `json:"name" v:"required"`
	Code           string `json:"code" v:"required"`
	Type           int    `json:"type" v:"required"`
	CronExpr       string `json:"cron_expr"`
	DelaySeconds   int    `json:"delay_seconds"`
	Handler        string `json:"handler" v:"required"`
	Params         string `json:"params"`
	RetryTimes     int    `json:"retry_times"`
	RetryInterval  int    `json:"retry_interval"`
	Concurrency    int    `json:"concurrency"`
	AlertOnFail    int    `json:"alert_on_fail"`
	AlertReceivers string `json:"alert_receivers"`
	Status         int    `json:"status"`
	Remark         string `json:"remark"`
}

type TaskCreateRes struct {
	Id int64 `json:"id"`
}

type TaskUpdateReq struct {
	g.Meta         `path:"/task/{id}" method:"put" tags:"任务" summary:"更新任务"`
	Id             int64  `json:"id" in:"path" v:"required|min:1"`
	Name           string `json:"name" v:"required"`
	Code           string `json:"code" v:"required"`
	Type           int    `json:"type" v:"required"`
	CronExpr       string `json:"cron_expr"`
	DelaySeconds   int    `json:"delay_seconds"`
	Handler        string `json:"handler" v:"required"`
	Params         string `json:"params"`
	RetryTimes     int    `json:"retry_times"`
	RetryInterval  int    `json:"retry_interval"`
	Concurrency    int    `json:"concurrency"`
	AlertOnFail    int    `json:"alert_on_fail"`
	AlertReceivers string `json:"alert_receivers"`
	Status         int    `json:"status"`
	Remark         string `json:"remark"`
}

type TaskUpdateRes struct{}

type TaskDeleteReq struct {
	g.Meta `path:"/task/{id}" method:"delete" tags:"任务" summary:"删除任务"`
	Id     int64 `json:"id" in:"path" v:"required|min:1"`
}

type TaskDeleteRes struct{}

type TaskRunReq struct {
	g.Meta `path:"/task/run" method:"post" tags:"任务" summary:"手动执行任务"`
	Id     int64 `json:"id" v:"required|min:1"`
}

type TaskRunRes struct {
	RunId string `json:"run_id"`
}

type TaskLogListReq struct {
	g.Meta  `path:"/task/log" method:"get" tags:"任务" summary:"任务执行日志"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	TaskId  int64  `json:"task_id"`
	RunId   string `json:"run_id"`
	Status  *int   `json:"status"`
}

type TaskLogListRes struct {
	List  []*TaskLogItem `json:"list"`
	Total int            `json:"total"`
}

type TaskLogItem struct {
	Id          int64  `json:"id"`
	TaskId      int64  `json:"task_id"`
	RunId       string `json:"run_id"`
	TriggerType int    `json:"trigger_type"`
	Status      int    `json:"status"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	DurationMs  int    `json:"duration_ms"`
	RetryCount  int    `json:"retry_count"`
	ErrorMsg    string `json:"error_msg"`
	Result      string `json:"result"`
	Node        string `json:"node"`
	CreateTime  string `json:"create_time"`
}
