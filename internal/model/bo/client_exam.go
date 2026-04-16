package bo

import examentity "exam/internal/model/entity/exam"

// SaveAnswerItem 批量保存中的一题。
type SaveAnswerItem struct {
	QuestionID      int64
	AnswerJSON      string
	ExpectedVersion *int
}

// AttemptView 会话详情（接口返回）。
type AttemptView struct {
	Attempt          examentity.ExamAttempt
	ServerTime       string
	DeadlineReached  bool
	SegmentCode      string
	RemainingSeconds *int // 当前 segment_code 对应环节剩余时间（秒）
}

// AnswerPayload 客户端答题 JSON。
type AnswerPayload struct {
	SelectedOptionIDs []int64 `json:"selected_option_ids"`
	Text              string  `json:"text"`
}

// QuestionScoreMeta 阅卷用题目元数据。
type QuestionScoreMeta struct {
	QuestionID    int64
	IsExample     int
	IsSubjective  int
	Score         float64
	CorrectOptIDs []int64
}

// RandomAnswerDraftItem 随机填答案草稿（不入库）。Answer 与 PUT …/answers 的 items[].answer 一致：客观题为选项 id（数字），主观题为文本字符串。
type RandomAnswerDraftItem struct {
	QuestionID int64
	Answer     any
}

// AttemptAnswerClientItem GET /exam/attempts/{id}/answers 单条，与保存接口 items 元素语义一致。
type AttemptAnswerClientItem struct {
	QuestionID int64
	Answer     any
}
