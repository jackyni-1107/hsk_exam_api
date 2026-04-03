// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamOption is the golang structure of table exam_option for DAO operations like Where/Data.
type ExamOption struct {
	g.Meta     `orm:"table:exam_option, do:true"`
	Id         any         // 主键
	QuestionId any         // 试题ID exam_question.id
	Flag       any         // 选项标识 A/B/C/T/F
	SortOrder  any         // 对应 answers.index
	IsCorrect  any         // 是否正确
	OptionType any         // text/image/pinyin 等
	Content    any         // 文本或资源文件名
	Creator    any         // 创建者
	CreateTime *gtime.Time // 创建时间
	Updater    any         // 更新者
	UpdateTime *gtime.Time // 更新时间
	DeleteFlag any         // 逻辑删除
}
