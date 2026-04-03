// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamResult is the golang structure for table exam_result.
type ExamResult struct {
	AttemptId              int64       `json:"attempt_id"                orm:"attempt_id"                description:"exam_attempt.id"`                  // exam_attempt.id
	MemberId               int64       `json:"member_id"                 orm:"member_id"                 description:"sys_member.id"`                    // sys_member.id
	ExamPaperId            int64       `json:"exam_paper_id"             orm:"exam_paper_id"             description:"exam_paper.id"`                    // exam_paper.id
	MockExaminationPaperId int64       `json:"mock_examination_paper_id" orm:"mock_examination_paper_id" description:"冗余 mock_examination_paper.id"`     // 冗余 mock_examination_paper.id
	Status                 int         `json:"status"                    orm:"status"                    description:"与 exam_attempt.status，列表阶段为已结束 4"` // 与 exam_attempt.status，列表阶段为已结束 4
	ObjectiveScore         float64     `json:"objective_score"           orm:"objective_score"           description:""`                                 //
	SubjectiveScore        float64     `json:"subjective_score"          orm:"subjective_score"          description:""`                                 //
	TotalScore             float64     `json:"total_score"               orm:"total_score"               description:""`                                 //
	HasSubjective          int         `json:"has_subjective"            orm:"has_subjective"            description:""`                                 //
	StartedAt              *gtime.Time `json:"started_at"                orm:"started_at"                description:""`                                 //
	SubmittedAt            *gtime.Time `json:"submitted_at"              orm:"submitted_at"              description:""`                                 //
	EndedAt                *gtime.Time `json:"ended_at"                  orm:"ended_at"                  description:""`                                 //
	CreateTime             *gtime.Time `json:"create_time"               orm:"create_time"               description:"同步自会话创建时间"`                        // 同步自会话创建时间
	UpdateTime             *gtime.Time `json:"update_time"               orm:"update_time"               description:""`                                 //
	DeleteFlag             int         `json:"delete_flag"               orm:"delete_flag"               description:""`                                 //
}
