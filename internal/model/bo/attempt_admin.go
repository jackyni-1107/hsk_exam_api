package bo

import (
	"github.com/gogf/gf/v2/os/gtime"

	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"
)

// AttemptAdminListRow 管理端会话列表行（与 Raw 列别名一致）。
type AttemptAdminListRow struct {
	Id                 int64       `json:"id"`
	MemberId           int64       `json:"member_id" orm:"member_id"`
	ExaminationPaperId int64       `json:"examination_paper_id"`
	ExamBatchId        int64       `json:"exam_batch_id"`
	MockLevelId        int64       `json:"mock_level_id"`
	Status             int         `json:"status"`
	ObjectiveScore     float64     `json:"objective_score"`
	SubjectiveScore    float64     `json:"subjective_score"`
	TotalScore         float64     `json:"total_score"`
	HasSubjective      int         `json:"has_subjective"`
	StartedAt          *gtime.Time `json:"started_at"`
	SubmittedAt        *gtime.Time `json:"submitted_at"`
	EndedAt            *gtime.Time `json:"ended_at"`
	CreateTime         *gtime.Time `json:"create_time"`
	Username           string      `json:"username"`
	Nickname           string      `json:"nickname"`
	PaperTitle         string      `json:"paper_title"`
	PaperLevel         string      `json:"paper_level"`
	RemotePaperId      string      `json:"remote_paper_id"`
}

// AttemptAdminAnswerRow 单题答题展示行。
type AttemptAdminAnswerRow struct {
	Answer           examentity.ExamAttemptAnswer
	Question         examentity.ExamQuestion
	Section          *examentity.ExamSection
	Options          []examentity.ExamOption
	ObjectiveCorrect *bool
}

// AttemptAdminDetailView 管理端会话详情。
type AttemptAdminDetailView struct {
	// ResultStatus 为 exam_result.status（无结果行时为 0）；5=全部算分完成，1–4 与 attempt 阶段对齐或已结束待主观。
	ResultStatus int `json:"result_status"`
	Attempt      examentity.ExamAttempt
	User         sysentity.SysMember
	Paper        examentity.ExamPaper
	Answers      []AttemptAdminAnswerRow
}

// SubjectiveScoreItem 主观题得分项。
type SubjectiveScoreItem struct {
	QuestionID int64
	Score      float64
}
