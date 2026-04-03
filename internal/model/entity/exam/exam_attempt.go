// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttempt is the golang structure for table exam_attempt.
type ExamAttempt struct {
	Id                     int64       `json:"id"                        orm:"id"                        description:"主键"`                                              // 主键
	MemberId               int64       `json:"member_id"                 orm:"member_id"                 description:"sys_member.id"`                                   // sys_member.id
	ExamPaperId            int64       `json:"exam_paper_id"             orm:"exam_paper_id"             description:"exam_paper.id"`                                   // exam_paper.id
	MockExaminationPaperId int64       `json:"mock_examination_paper_id" orm:"mock_examination_paper_id" description:"冗余 mock_examination_paper.id"`                    // 冗余 mock_examination_paper.id
	ExamBatchId            int64       `json:"exam_batch_id"             orm:"exam_batch_id"             description:"exam_batch.id，0=历史非批次"`                           // exam_batch.id，0=历史非批次
	MockLevelId            int64       `json:"mock_level_id"             orm:"mock_level_id"             description:"mock_levels.id，0=历史非批次"`                          // mock_levels.id，0=历史非批次
	Status                 int         `json:"status"                    orm:"status"                    description:"1=not_started 2=in_progress 3=submitted 4=ended"` // 1=not_started 2=in_progress 3=submitted 4=ended
	DurationSeconds        int         `json:"duration_seconds"          orm:"duration_seconds"          description:"开考时快照时长秒"`                                        // 开考时快照时长秒
	StartedAt              *gtime.Time `json:"started_at"                orm:"started_at"                description:"开考时间"`                                            // 开考时间
	DeadlineAt             *gtime.Time `json:"deadline_at"               orm:"deadline_at"               description:"截止时间"`                                            // 截止时间
	SubmittedAt            *gtime.Time `json:"submitted_at"              orm:"submitted_at"              description:"交卷时间"`                                            // 交卷时间
	EndedAt                *gtime.Time `json:"ended_at"                  orm:"ended_at"                  description:"结束/阅卷完成时间"`                                       // 结束/阅卷完成时间
	ObjectiveScore         float64     `json:"objective_score"           orm:"objective_score"           description:"客观题得分"`                                           // 客观题得分
	SubjectiveScore        float64     `json:"subjective_score"          orm:"subjective_score"          description:"主观题得分（未批阅时为空）"`                                   // 主观题得分（未批阅时为空）
	TotalScore             float64     `json:"total_score"               orm:"total_score"               description:"总分（无主观或主观未批时可能为客观分）"`                             // 总分（无主观或主观未批时可能为客观分）
	HasSubjective          int         `json:"has_subjective"            orm:"has_subjective"            description:"本卷是否含主观题（交卷快照）"`                                  // 本卷是否含主观题（交卷快照）
	Creator                string      `json:"creator"                   orm:"creator"                   description:"创建者"`                                             // 创建者
	CreateTime             *gtime.Time `json:"create_time"               orm:"create_time"               description:"创建时间"`                                            // 创建时间
	Updater                string      `json:"updater"                   orm:"updater"                   description:"更新者"`                                             // 更新者
	UpdateTime             *gtime.Time `json:"update_time"               orm:"update_time"               description:"更新时间"`                                            // 更新时间
	DeleteFlag             int         `json:"delete_flag"               orm:"delete_flag"               description:"逻辑删除"`                                            // 逻辑删除
}
