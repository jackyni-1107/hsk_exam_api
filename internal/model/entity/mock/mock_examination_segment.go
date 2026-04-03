// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package mock

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MockExaminationSegment is the golang structure for table mock_examination_segment.
type MockExaminationSegment struct {
	Id               int64       `json:"id"                 orm:"id"                 description:"id 主键"`     // id 主键
	LevelId          int64       `json:"level_id"           orm:"level_id"           description:"HSK等级ID"`   // HSK等级ID
	SegmentCode      string      `json:"segment_code"       orm:"segment_code"       description:"环节编码"`      // 环节编码
	SegmentName      string      `json:"segment_name"       orm:"segment_name"       description:"环节名称"`      // 环节名称
	SegmentNameTrans string      `json:"segment_name_trans" orm:"segment_name_trans" description:"环节名称国际化"`   // 环节名称国际化
	SegmentDesc      string      `json:"segment_desc"       orm:"segment_desc"       description:"segment说明"` // segment说明
	ScoreFull        int         `json:"score_full"         orm:"score_full"         description:"环节满分"`      // 环节满分
	QuestionCount    int         `json:"question_count"     orm:"question_count"     description:"题目数量"`      // 题目数量
	Duration         int         `json:"duration"           orm:"duration"           description:"环节时长 分钟数"`  // 环节时长 分钟数
	Seq              int         `json:"seq"                orm:"seq"                description:"顺序号"`       // 顺序号
	DeleteFlag       int         `json:"delete_flag"        orm:"delete_flag"        description:"是否删除"`      // 是否删除
	CreateTime       *gtime.Time `json:"create_time"        orm:"create_time"        description:"创建时间"`      // 创建时间
	UpdateTime       *gtime.Time `json:"update_time"        orm:"update_time"        description:"更新时间"`      // 更新时间
}
