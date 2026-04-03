// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamQuestionBlock is the golang structure of table exam_question_block for DAO operations like Where/Data.
type ExamQuestionBlock struct {
	g.Meta                  `orm:"table:exam_question_block, do:true"`
	Id                      any         // 主键
	SectionId               any         // 大题ID exam_section.id
	BlockOrder              any         // 对应 topic JSON 中 items 下标
	GroupIndex              any         // 套题外层 index（若存在）
	QuestionDescriptionJson any         // 块级 question_description_obj 等
	Creator                 any         // 创建者
	CreateTime              *gtime.Time // 创建时间
	Updater                 any         // 更新者
	UpdateTime              *gtime.Time // 更新时间
	DeleteFlag              any         // 逻辑删除
}
