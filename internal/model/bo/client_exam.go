package bo

import examentity "exam/internal/model/entity/exam"

// SaveAnswerItem 批量保存中的一题。OptionID 与 Text 互斥：客观题传 OptionID，主观题/填空传 Text。
type SaveAnswerItem struct {
	QuestionID      int64
	OptionID        int64
	Text            string
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

// AnswerPayload 客户端答题 JSON 存储载荷。OptionID 与 Text 互斥，零值表示未作答。
type AnswerPayload struct {
	OptionID int64  `json:"o_id,omitempty"`
	Text     string `json:"text,omitempty"`
}

// QuestionScoreMeta 阅卷用题目元数据。
type QuestionScoreMeta struct {
	QuestionID    int64
	SegmentCode   string
	IsExample     int
	IsSubjective  int
	Score         float64
	CorrectOptIDs []int64
}

// RandomAnswerDraftItem 随机填答案草稿（不入库）。OptionID 与 Text 互斥。
type RandomAnswerDraftItem struct {
	QuestionID int64
	OptionID   int64
	Text       string
}

// AttemptAnswerClientItem GET /exam/attempts/{id}/answers 单条，与保存接口 items 元素语义一致。
type AttemptAnswerClientItem struct {
	QuestionID int64
	OptionID   int64
	Text       string
}

// AttemptCheatEventRecordInput 会话作弊事件上报入参。
type AttemptCheatEventRecordInput struct {
	AttemptID   int64
	UserID      int64
	EventType   string
	SegmentCode string
	Detail      string
	IP          string
	UserAgent   string
}
