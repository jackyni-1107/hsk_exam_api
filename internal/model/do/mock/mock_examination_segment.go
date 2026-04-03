// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MockExaminationSegment is the golang structure of table mock_examination_segment for DAO operations like Where/Data.
type MockExaminationSegment struct {
	g.Meta           `orm:"table:mock_examination_segment, do:true"`
	Id               any         // id 主键
	LevelId          any         // HSK等级ID
	SegmentCode      any         // 环节编码
	SegmentName      any         // 环节名称
	SegmentNameTrans any         // 环节名称国际化
	SegmentDesc      any         // segment说明
	ScoreFull        any         // 环节满分
	QuestionCount    any         // 题目数量
	Duration         any         // 环节时长 分钟数
	Seq              any         // 顺序号
	DeleteFlag       any         // 是否删除
	CreateTime       *gtime.Time // 创建时间
	UpdateTime       *gtime.Time // 更新时间
}
