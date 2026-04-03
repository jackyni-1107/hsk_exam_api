// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamQuestionBlock is the golang structure for table exam_question_block.
type ExamQuestionBlock struct {
	Id                      int64       `json:"id"                        orm:"id"                        description:"主键"`                            // 主键
	SectionId               int64       `json:"section_id"                orm:"section_id"                description:"大题ID exam_section.id"`          // 大题ID exam_section.id
	BlockOrder              int         `json:"block_order"               orm:"block_order"               description:"对应 topic JSON 中 items 下标"`      // 对应 topic JSON 中 items 下标
	GroupIndex              int         `json:"group_index"               orm:"group_index"               description:"套题外层 index（若存在）"`               // 套题外层 index（若存在）
	QuestionDescriptionJson string      `json:"question_description_json" orm:"question_description_json" description:"块级 question_description_obj 等"` // 块级 question_description_obj 等
	Creator                 string      `json:"creator"                   orm:"creator"                   description:"创建者"`                           // 创建者
	CreateTime              *gtime.Time `json:"create_time"               orm:"create_time"               description:"创建时间"`                          // 创建时间
	Updater                 string      `json:"updater"                   orm:"updater"                   description:"更新者"`                           // 更新者
	UpdateTime              *gtime.Time `json:"update_time"               orm:"update_time"               description:"更新时间"`                          // 更新时间
	DeleteFlag              int         `json:"delete_flag"               orm:"delete_flag"               description:"逻辑删除"`                          // 逻辑删除
}
