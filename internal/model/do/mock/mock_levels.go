// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MockLevels is the golang structure of table mock_levels for DAO operations like Where/Data.
type MockLevels struct {
	g.Meta             `orm:"table:mock_levels, do:true"`
	Id                 any         // 主键
	LevelId            any         // 在指定类型下的等级
	LevelType          any         // 等级类型 （1 hsk 2 hskk 3 yct）
	TypeName           any         // 类型名称
	LevelName          any         // hsk等级名称
	AppLevelName       any         // app专用等级名称
	LevelNameTrans     any         // 等级名称多语言
	SignUpUrl          any         // 报名地址
	LevelDesc          any         // 等级简介
	CnDesc             any         // 中文等级描述
	EnDesc             any         // 英文等级描述
	LevelExplainAudio  any         // 模拟考试音音频地址
	DeleteFlag         any         // 逻辑删除标识
	MarkFlag           any         // 是否需要批改
	CreateTime         *gtime.Time // 创建时间
	UpdateTime         *gtime.Time // 更新时间
	ExamShowStatus     any         // 模拟考显示状态，0表示不显示，1表示显示
	HomeworkShowStatus any         // 作业显示状态，0表示不显示，1表示显示
	PaperAnswerTime    any         // 该等级下试卷作答时间
	ResourceType       any         // 1：考试列表 2：文件资源列表
}
