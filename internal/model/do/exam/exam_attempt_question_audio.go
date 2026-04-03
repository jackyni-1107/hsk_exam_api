// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttemptQuestionAudio is the golang structure of table exam_attempt_question_audio for DAO operations like Where/Data.
type ExamAttemptQuestionAudio struct {
	g.Meta          `orm:"table:exam_attempt_question_audio, do:true"`
	Id              any         // 主键
	AttemptId       any         // exam_attempt.id
	ExamQuestionId  any         // exam_question.id
	MaxSegmentIndex any         // 允许播放到该分片索引（含）；与 0..segment_count-1 求交
	Creator         any         // 创建者
	CreateTime      *gtime.Time // 创建时间
	Updater         any         // 更新者
	UpdateTime      *gtime.Time // 更新时间
	DeleteFlag      any         // 逻辑删除
}
