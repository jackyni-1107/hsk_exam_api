// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MockLevels is the golang structure for table mock_levels.
type MockLevels struct {
	Id                 int64       `json:"id"                   orm:"id"                   description:"主键"`                        // 主键
	LevelId            int         `json:"level_id"             orm:"level_id"             description:"在指定类型下的等级"`                 // 在指定类型下的等级
	LevelType          int         `json:"level_type"           orm:"level_type"           description:"等级类型 （1 hsk 2 hskk 3 yct）"` // 等级类型 （1 hsk 2 hskk 3 yct）
	TypeName           string      `json:"type_name"            orm:"type_name"            description:"类型名称"`                      // 类型名称
	LevelName          string      `json:"level_name"           orm:"level_name"           description:"hsk等级名称"`                   // hsk等级名称
	AppLevelName       string      `json:"app_level_name"       orm:"app_level_name"       description:"app专用等级名称"`                 // app专用等级名称
	LevelNameTrans     string      `json:"level_name_trans"     orm:"level_name_trans"     description:"等级名称多语言"`                   // 等级名称多语言
	SignUpUrl          string      `json:"sign_up_url"          orm:"sign_up_url"          description:"报名地址"`                      // 报名地址
	LevelDesc          string      `json:"level_desc"           orm:"level_desc"           description:"等级简介"`                      // 等级简介
	CnDesc             string      `json:"cn_desc"              orm:"cn_desc"              description:"中文等级描述"`                    // 中文等级描述
	EnDesc             string      `json:"en_desc"              orm:"en_desc"              description:"英文等级描述"`                    // 英文等级描述
	LevelExplainAudio  string      `json:"level_explain_audio"  orm:"level_explain_audio"  description:"模拟考试音音频地址"`                 // 模拟考试音音频地址
	DeleteFlag         int         `json:"delete_flag"          orm:"delete_flag"          description:"逻辑删除标识"`                    // 逻辑删除标识
	MarkFlag           int         `json:"mark_flag"            orm:"mark_flag"            description:"是否需要批改"`                    // 是否需要批改
	CreateTime         *gtime.Time `json:"create_time"          orm:"create_time"          description:"创建时间"`                      // 创建时间
	UpdateTime         *gtime.Time `json:"update_time"          orm:"update_time"          description:"更新时间"`                      // 更新时间
	ExamShowStatus     int         `json:"exam_show_status"     orm:"exam_show_status"     description:"模拟考显示状态，0表示不显示，1表示显示"`      // 模拟考显示状态，0表示不显示，1表示显示
	HomeworkShowStatus int         `json:"homework_show_status" orm:"homework_show_status" description:"作业显示状态，0表示不显示，1表示显示"`       // 作业显示状态，0表示不显示，1表示显示
	PaperAnswerTime    int         `json:"paper_answer_time"    orm:"paper_answer_time"    description:"该等级下试卷作答时间"`                // 该等级下试卷作答时间
	ResourceType       int         `json:"resource_type"        orm:"resource_type"        description:"1：考试列表 2：文件资源列表"`           // 1：考试列表 2：文件资源列表
}
