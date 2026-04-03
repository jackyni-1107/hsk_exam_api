// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MockExaminationPart is the golang structure for table mock_examination_part.
type MockExaminationPart struct {
	Id                      int64       `json:"id"                        orm:"id"                        description:"id主键"`      // id主键
	Code                    int         `json:"code"                      orm:"code"                      description:"部分编号"`      // 部分编号
	SegmentId               int64       `json:"segment_id"                orm:"segment_id"                description:"环节id"`      // 环节id
	QuestionCount           int         `json:"question_count"            orm:"question_count"            description:"题目数量"`      // 题目数量
	ObjectiveQuestionCount  int         `json:"objective_question_count"  orm:"objective_question_count"  description:"客观题数量"`     // 客观题数量
	SubjectiveQuestionCount int         `json:"subjective_question_count" orm:"subjective_question_count" description:"主观题数量"`     // 主观题数量
	QuestionScore           float64     `json:"question_score"            orm:"question_score"            description:"每个题目得分"`    // 每个题目得分
	PartScore               float64     `json:"part_score"                orm:"part_score"                description:"该题型分数"`     // 该题型分数
	AnswerTime              int         `json:"answer_time"               orm:"answer_time"               description:"每个题型的回答时间"` // 每个题型的回答时间
	PartName                string      `json:"part_name"                 orm:"part_name"                 description:"只用于hskk"`   // 只用于hskk
	PartNameTrans           string      `json:"part_name_trans"           orm:"part_name_trans"           description:"部分名称多语言"`   // 部分名称多语言
	DeleteFlag              int         `json:"delete_flag"               orm:"delete_flag"               description:"是否删除"`      // 是否删除
	CreateTime              *gtime.Time `json:"create_time"               orm:"create_time"               description:"创建时间"`      // 创建时间
	UpdateTime              *gtime.Time `json:"update_time"               orm:"update_time"               description:"更新时间"`      // 更新时间
}
