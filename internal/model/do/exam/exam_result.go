// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamResult is the golang structure of table exam_result for DAO operations like Where/Data.
type ExamResult struct {
	g.Meta                 `orm:"table:exam_result, do:true"`
	AttemptId              any         // exam_attempt.id
	MemberId               any         // sys_member.id
	ExamPaperId            any         // exam_paper.id
	MockExaminationPaperId any         // 冗余 mock_examination_paper.id
	Status                 any         // 与 exam_attempt.status，列表阶段为已结束 4
	ObjectiveScore         any         //
	SubjectiveScore        any         //
	TotalScore             any         //
	HasSubjective          any         //
	StartedAt              *gtime.Time //
	SubmittedAt            *gtime.Time //
	EndedAt                *gtime.Time //
	CreateTime             *gtime.Time // 同步自会话创建时间
	UpdateTime             *gtime.Time //
	DeleteFlag             any         //
}
