package consts

// 批次类型 exam_batch.batch_kind
const (
	ExamBatchKindFormal   = 0
	ExamBatchKindPractice = 1
	// ExamBatchKindFilterAll 管理端结果列表等：不按 batch_kind 过滤（含无批次会话）
	ExamBatchKindFilterAll = -1
)
