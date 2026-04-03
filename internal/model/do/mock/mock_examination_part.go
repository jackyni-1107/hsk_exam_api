// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MockExaminationPart is the golang structure of table mock_examination_part for DAO operations like Where/Data.
type MockExaminationPart struct {
	g.Meta                  `orm:"table:mock_examination_part, do:true"`
	Id                      any         // id主键
	Code                    any         // 部分编号
	SegmentId               any         // 环节id
	QuestionCount           any         // 题目数量
	ObjectiveQuestionCount  any         // 客观题数量
	SubjectiveQuestionCount any         // 主观题数量
	QuestionScore           any         // 每个题目得分
	PartScore               any         // 该题型分数
	AnswerTime              any         // 每个题型的回答时间
	PartName                any         // 只用于hskk
	PartNameTrans           any         // 部分名称多语言
	DeleteFlag              any         // 是否删除
	CreateTime              *gtime.Time // 创建时间
	UpdateTime              *gtime.Time // 更新时间
}
