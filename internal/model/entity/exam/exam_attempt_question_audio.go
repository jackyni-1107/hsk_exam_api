// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttemptQuestionAudio is the golang structure for table exam_attempt_question_audio.
type ExamAttemptQuestionAudio struct {
	Id              int64       `json:"id"                orm:"id"                description:"主键"`                                    // 主键
	AttemptId       int64       `json:"attempt_id"        orm:"attempt_id"        description:"exam_attempt.id"`                       // exam_attempt.id
	ExamQuestionId  int64       `json:"exam_question_id"  orm:"exam_question_id"  description:"exam_question.id"`                      // exam_question.id
	MaxSegmentIndex int         `json:"max_segment_index" orm:"max_segment_index" description:"允许播放到该分片索引（含）；与 0..segment_count-1 求交"` // 允许播放到该分片索引（含）；与 0..segment_count-1 求交
	Creator         string      `json:"creator"           orm:"creator"           description:"创建者"`                                   // 创建者
	CreateTime      *gtime.Time `json:"create_time"       orm:"create_time"       description:"创建时间"`                                  // 创建时间
	Updater         string      `json:"updater"           orm:"updater"           description:"更新者"`                                   // 更新者
	UpdateTime      *gtime.Time `json:"update_time"       orm:"update_time"       description:"更新时间"`                                  // 更新时间
	DeleteFlag      int         `json:"delete_flag"       orm:"delete_flag"       description:"逻辑删除"`                                  // 逻辑删除
}
