// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamOption is the golang structure for table exam_option.
type ExamOption struct {
	Id         int64       `json:"id"          orm:"id"          description:"主键"`                    // 主键
	QuestionId int64       `json:"question_id" orm:"question_id" description:"试题ID exam_question.id"` // 试题ID exam_question.id
	Flag       string      `json:"flag"        orm:"flag"        description:"选项标识 A/B/C/T/F"`        // 选项标识 A/B/C/T/F
	SortOrder  int         `json:"sort_order"  orm:"sort_order"  description:"对应 answers.index"`      // 对应 answers.index
	IsCorrect  int         `json:"is_correct"  orm:"is_correct"  description:"是否正确"`                  // 是否正确
	OptionType string      `json:"option_type" orm:"option_type" description:"text/image/pinyin 等"`   // text/image/pinyin 等
	Content    string      `json:"content"     orm:"content"     description:"文本或资源文件名"`              // 文本或资源文件名
	Creator    string      `json:"creator"     orm:"creator"     description:"创建者"`                   // 创建者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"创建时间"`                  // 创建时间
	Updater    string      `json:"updater"     orm:"updater"     description:"更新者"`                   // 更新者
	UpdateTime *gtime.Time `json:"update_time" orm:"update_time" description:"更新时间"`                  // 更新时间
	DeleteFlag int         `json:"delete_flag" orm:"delete_flag" description:"逻辑删除"`                  // 逻辑删除
}
