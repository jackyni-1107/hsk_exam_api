// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttempt is the golang structure of table exam_attempt for DAO operations like Where/Data.
type ExamAttempt struct {
	g.Meta                 `orm:"table:exam_attempt, do:true"`
	Id                     any         // 主键
	MemberId               any         // sys_member.id
	ExamPaperId            any         // exam_paper.id
	MockExaminationPaperId any         // 冗余 mock_examination_paper.id
	ExamBatchId            any         // exam_batch.id，0=历史非批次
	MockLevelId            any         // mock_levels.id，0=历史非批次
	Status                 any         // 1=not_started 2=in_progress 3=submitted 4=ended
	DurationSeconds        any         // 开考时快照时长秒
	StartedAt              *gtime.Time // 开考时间
	DeadlineAt             *gtime.Time // 截止时间
	SubmittedAt            *gtime.Time // 交卷时间
	EndedAt                *gtime.Time // 结束/阅卷完成时间
	ObjectiveScore         any         // 客观题得分
	SubjectiveScore        any         // 主观题得分（未批阅时为空）
	TotalScore             any         // 总分（无主观或主观未批时可能为客观分）
	HasSubjective          any         // 本卷是否含主观题（交卷快照）
	Creator                any         // 创建者
	CreateTime             *gtime.Time // 创建时间
	Updater                any         // 更新者
	UpdateTime             *gtime.Time // 更新时间
	DeleteFlag             any         // 逻辑删除
}
