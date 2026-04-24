package consts

// exam_result.status：1–4 与 exam_attempt.status 一致；5 表示全部算分流程已完成（不可再次执行 Finalize 所代表的客观算分任务；主观分已写入后亦不可重复评阅）。
const (
	ExamResultScoringComplete = 5
)
