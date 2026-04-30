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
	ExamBatchId            any         // 冗余 exam_batch.id
	MockLevelId            any         // 冗余 mock_levels.id
	Status                 any         // 与 exam_attempt.status，列表阶段为已结束 4
	ObjectiveScore         any         //
	SubjectiveScore        any         //
	TotalScore             any         //
	SegmentScoreJson       any         // 按 segment_code 的整数分数字典 JSON
	HasSubjective          any         //
	StartedAt              *gtime.Time //
	SubmittedAt            *gtime.Time //
	EndedAt                *gtime.Time //
	CreateTime             *gtime.Time // 同步自会话创建时间
	UpdateTime             *gtime.Time //
	DeleteFlag             any         //
}
