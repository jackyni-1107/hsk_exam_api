// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttemptAnswer is the golang structure for table exam_attempt_answer.
type ExamAttemptAnswer struct {
	Id             int64       `json:"id"               orm:"id"               description:"主键"`                // 主键
	AttemptId      int64       `json:"attempt_id"       orm:"attempt_id"       description:"exam_attempt.id"`   // exam_attempt.id
	ExamQuestionId int64       `json:"exam_question_id" orm:"exam_question_id" description:"exam_question.id"`  // exam_question.id
	AnswerJson     string      `json:"answer_json"      orm:"answer_json"      description:"用户答案JSON"`          // 用户答案JSON
	AwardedScore   float64     `json:"awarded_score"    orm:"awarded_score"    description:"主观题人工得分；NULL 表示未评"` // 主观题人工得分；NULL 表示未评
	Version        int         `json:"version"          orm:"version"          description:"乐观锁版本"`             // 乐观锁版本
	Creator        string      `json:"creator"          orm:"creator"          description:"创建者"`               // 创建者
	CreateTime     *gtime.Time `json:"create_time"      orm:"create_time"      description:"创建时间"`              // 创建时间
	Updater        string      `json:"updater"          orm:"updater"          description:"更新者"`               // 更新者
	UpdateTime     *gtime.Time `json:"update_time"      orm:"update_time"      description:"更新时间"`              // 更新时间
	DeleteFlag     int         `json:"delete_flag"      orm:"delete_flag"      description:"逻辑删除"`              // 逻辑删除
}
