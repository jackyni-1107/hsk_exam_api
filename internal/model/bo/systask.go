package bo

type TaskCreateInput struct {
	Name           string
	Code           string
	CronExpr       string
	Handler        string
	Params         string
	AlertReceivers string
	Remark         string
	Creator        string
	Type           int
	DelaySeconds   int
	RetryTimes     int
	RetryInterval  int
	Concurrency    int
	AlertOnFail    int
	Status         int
}

type TaskUpdateInput struct {
	Id             int64
	Name           string
	Code           string
	CronExpr       string
	Handler        string
	Params         string
	AlertReceivers string
	Remark         string
	Updater        string
	Type           int
	DelaySeconds   int
	RetryTimes     int
	RetryInterval  int
	Concurrency    int
	AlertOnFail    int
	Status         int
}

type TaskRuntimeStats struct {
	DelayQueueSize         int
	DelayDueCount          int
	DelayScannerActive     bool
	DelayScannerTTLMillis  int64
	DelayOldestDueAtMillis int64
}
