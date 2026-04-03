package bo

import (
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/model/entity"
)

// AttemptAdminListRow 管理端会话列表行（与 Raw 列别名一致）。
type AttemptAdminListRow struct {
	Id                 int64       `json:"id"`
	ClientUserId       int64       `json:"client_user_id"`
	ExaminationPaperId int64       `json:"examination_paper_id"`
	Status             int         `json:"status"`
	ObjectiveScore  float64     `json:"objective_score"`
	SubjectiveScore float64     `json:"subjective_score"`
	TotalScore      float64     `json:"total_score"`
	HasSubjective   int         `json:"has_subjective"`
	StartedAt       *gtime.Time `json:"started_at"`
	SubmittedAt     *gtime.Time `json:"submitted_at"`
	EndedAt         *gtime.Time `json:"ended_at"`
	CreateTime      *gtime.Time `json:"create_time"`
	Username        string      `json:"username"`
	Nickname        string      `json:"nickname"`
	PaperTitle      string      `json:"paper_title"`
	PaperLevel      string      `json:"paper_level"`
	RemotePaperId   string      `json:"remote_paper_id"`
}

// AttemptAdminAnswerRow 单题答题展示行。
type AttemptAdminAnswerRow struct {
	Answer           entity.ExamAttemptAnswer
	Question         entity.ExamQuestion
	Section          *entity.ExamSection
	Options          []entity.ExamOption
	ObjectiveCorrect *bool
}

// AttemptAdminDetailView 管理端会话详情。
type AttemptAdminDetailView struct {
	Attempt entity.ExamAttempt
	User    entity.ClientUser
	Paper   entity.ExamPaper
	Answers []AttemptAdminAnswerRow
}

// SubjectiveScoreItem 主观题得分项。
type SubjectiveScoreItem struct {
	QuestionID int64
	Score      float64
}
