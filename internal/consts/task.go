package consts

// 任务类型（sys_task.type）
const (
	TaskTypeCron  = 1 // 定时 cron
	TaskTypeDelay = 2 // 延迟
)

// 任务启用状态（sys_task.status）
const (
	TaskStatusEnabled  = 0
	TaskStatusDisabled = 1
)

// 单次运行状态（sys_task_log.status）
const (
	TaskRunStatusRunning = 0
	TaskRunStatusSuccess = 1
	TaskRunStatusFailed  = 2
)

// 触发类型（sys_task_log.trigger_type）
const (
	TriggerTypeCron   = 1
	TriggerTypeDelay  = 2
	TriggerTypeManual = 3
)
