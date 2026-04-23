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

// ---------- Task Redis 键（hskexam:task:业务名:唯一标识） ----------

const (
	TaskClusterExecLockKeyFmt     = "hskexam:task:cluster_exec:%d"
	TaskClusterExecLockTTLSeconds = 300
	TaskSemKeyFmt                 = "hskexam:task:sem:%d"
	TaskSemExpireSeconds          = 3600
	TaskDelayQueueKey             = "hskexam:task:delay_queue"
	TaskDelayScannerLockKey       = "hskexam:task:delay_queue:scanner_lock"
	TaskDelayScannerLockTTLMillis = 1500
	TaskDelayPollIntervalSeconds  = 1
	TaskDelayBatchSize            = 100
)
