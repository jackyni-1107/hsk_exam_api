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
	RemainingSeconds *int // 进行中且存在批次结束时间时：exam_batch.exam_end_at 与最近一次保存答案时间的差（秒），下限 0
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
