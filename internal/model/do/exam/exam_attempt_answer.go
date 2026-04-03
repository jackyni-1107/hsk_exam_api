// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttemptAnswer is the golang structure of table exam_attempt_answer for DAO operations like Where/Data.
type ExamAttemptAnswer struct {
	g.Meta         `orm:"table:exam_attempt_answer, do:true"`
	Id             any         // 主键
	AttemptId      any         // exam_attempt.id
	ExamQuestionId any         // exam_question.id
	AnswerJson     any         // 用户答案JSON
	AwardedScore   any         // 主观题人工得分；NULL 表示未评
	Version        any         // 乐观锁版本
	Creator        any         // 创建者
	CreateTime     *gtime.Time // 创建时间
	Updater        any         // 更新者
	UpdateTime     *gtime.Time // 更新时间
	DeleteFlag     any         // 逻辑删除
}
