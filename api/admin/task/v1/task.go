package v1

import "github.com/gogf/gf/v2/frame/g"

type TaskListReq struct {
	g.Meta  `path:"/task/list" method:"get" tags:"任务" summary:"任务列表" permission:"task:list"`
	Page    int    `json:"page" dc:"页码"`
	Size    int    `json:"size" dc:"每页条数"`
	Name    string `json:"name" dc:"任务名称"`
	Code    string `json:"code" dc:"任务编码"`
	Type    int    `json:"type" dc:"任务类型"`
	Status  *int   `json:"status" dc:"状态"`
	Handler string `json:"handler" dc:"处理器"`
}

type TaskListRes struct {
	List  []*TaskItem `json:"list" dc:"列表"`
	Total int         `json:"total" dc:"总数"`
}

type TaskRuntimeStatsReq struct {
	g.Meta `path:"/task/runtime" method:"get" tags:"任务" summary:"任务运行时态统计"`
}

type TaskRuntimeStatsRes struct {
	DelayQueueSize        int    `json:"delay_queue_size" dc:"延迟队列总数"`
	DelayDueCount         int    `json:"delay_due_count" dc:"已到期待执行数量"`
	DelayScannerActive    bool   `json:"delay_scanner_active" dc:"延迟扫描器是否活跃"`
	DelayScannerTTLMillis int64  `json:"delay_scanner_ttl_millis" dc:"延迟扫描器锁剩余生存时间"`
	DelayOldestDueAt      string `json:"delay_oldest_due_at" dc:"队列最早到期时间"`
}

type TaskItem struct {
	Id             int64  `json:"id" dc:"任务ID"`
	Name           string `json:"name" dc:"任务名称"`
	Code           string `json:"code" dc:"任务编码"`
	Type           int    `json:"type" dc:"任务类型"`
	CronExpr       string `json:"cron_expr" dc:"Cron表达式"`
	DelaySeconds   int    `json:"delay_seconds" dc:"延迟秒数"`
	Handler        string `json:"handler" dc:"处理器"`
	Params         string `json:"params" dc:"参数(JSON)"`
	RetryTimes     int    `json:"retry_times" dc:"重试次数"`
	RetryInterval  int    `json:"retry_interval" dc:"重试间隔(秒)"`
	Concurrency    int    `json:"concurrency" dc:"并发数"`
	AlertOnFail    int    `json:"alert_on_fail" dc:"失败告警：0否 1是"`
	AlertReceivers string `json:"alert_receivers" dc:"告警接收人"`
	Status         int    `json:"status" dc:"状态：0启用 1停用"`
	Remark         string `json:"remark" dc:"备注"`
	CreateTime     string `json:"create_time" dc:"创建时间"`
}

type TaskCreateReq struct {
	g.Meta         `path:"/task" method:"post" tags:"任务" summary:"创建任务" permission:"task:create"`
	Name           string `json:"name" v:"required" dc:"任务名称"`
	Code           string `json:"code" v:"required" dc:"任务编码"`
	Type           int    `json:"type" v:"required" dc:"任务类型"`
	CronExpr       string `json:"cron_expr" dc:"Cron表达式"`
	DelaySeconds   int    `json:"delay_seconds" dc:"延迟秒数"`
	Handler        string `json:"handler" v:"required" dc:"处理器"`
	Params         string `json:"params" dc:"参数(JSON)"`
	RetryTimes     int    `json:"retry_times" dc:"重试次数"`
	RetryInterval  int    `json:"retry_interval" dc:"重试间隔(秒)"`
	Concurrency    int    `json:"concurrency" dc:"并发数"`
	AlertOnFail    int    `json:"alert_on_fail" dc:"失败告警：0否 1是"`
	AlertReceivers string `json:"alert_receivers" dc:"告警接收人"`
	Status         int    `json:"status" dc:"状态：0启用 1停用"`
	Remark         string `json:"remark" dc:"备注"`
}

type TaskCreateRes struct {
	Id int64 `json:"id" dc:"任务ID"`
}

type TaskUpdateReq struct {
	g.Meta         `path:"/task/{id}" method:"put" tags:"任务" summary:"更新任务" permission:"task:update"`
	Id             int64  `json:"id" in:"path" v:"required|min:1" dc:"任务ID"`
	Name           string `json:"name" v:"required" dc:"任务名称"`
	Code           string `json:"code" v:"required" dc:"任务编码"`
	Type           int    `json:"type" v:"required" dc:"任务类型"`
	CronExpr       string `json:"cron_expr" dc:"Cron表达式"`
	DelaySeconds   int    `json:"delay_seconds" dc:"延迟秒数"`
	Handler        string `json:"handler" v:"required" dc:"处理器"`
	Params         string `json:"params" dc:"参数(JSON)"`
	RetryTimes     int    `json:"retry_times" dc:"重试次数"`
	RetryInterval  int    `json:"retry_interval" dc:"重试间隔(秒)"`
	Concurrency    int    `json:"concurrency" dc:"并发数"`
	AlertOnFail    int    `json:"alert_on_fail" dc:"失败告警：0否 1是"`
	AlertReceivers string `json:"alert_receivers" dc:"告警接收人"`
	Status         int    `json:"status" dc:"状态：0启用 1停用"`
	Remark         string `json:"remark" dc:"备注"`
}

type TaskUpdateRes struct{}

type TaskDeleteReq struct {
	g.Meta `path:"/task/{id}" method:"delete" tags:"任务" summary:"删除任务" permission:"task:delete"`
	Id     int64 `json:"id" in:"path" v:"required|min:1" dc:"任务ID"`
}

type TaskDeleteRes struct{}

type TaskRunReq struct {
	g.Meta `path:"/task/run" method:"post" tags:"任务" summary:"手动执行任务" permission:"task:run"`
	Id     int64 `json:"id" v:"required|min:1" dc:"任务ID"`
}

type TaskRunRes struct {
	RunId string `json:"run_id" dc:"执行ID"`
}

type TaskLogListReq struct {
	g.Meta `path:"/task/log" method:"get" tags:"任务" summary:"任务执行日志" permission:"task:log"`
	Page   int    `json:"page" dc:"页码"`
	Size   int    `json:"size" dc:"每页条数"`
	TaskId int64  `json:"task_id" dc:"任务ID"`
	RunId  string `json:"run_id" dc:"执行ID"`
	Status *int   `json:"status" dc:"状态"`
}

type TaskLogListRes struct {
	List  []*TaskLogItem `json:"list" dc:"列表"`
	Total int            `json:"total" dc:"总数"`
}

type TaskLogItem struct {
	Id          int64  `json:"id" dc:"记录ID"`
	TaskId      int64  `json:"task_id" dc:"任务ID"`
	RunId       string `json:"run_id" dc:"执行ID"`
	TriggerType int    `json:"trigger_type" dc:"触发类型"`
	Status      int    `json:"status" dc:"执行状态"`
	StartTime   string `json:"start_time" dc:"开始时间"`
	EndTime     string `json:"end_time" dc:"结束时间"`
	DurationMs  int    `json:"duration_ms" dc:"耗时(毫秒)"`
	RetryCount  int    `json:"retry_count" dc:"重试次数"`
	ErrorMsg    string `json:"error_msg" dc:"错误信息"`
	Result      string `json:"result" dc:"执行结果"`
	Node        string `json:"node" dc:"执行节点"`
	CreateTime  string `json:"create_time" dc:"创建时间"`
}
