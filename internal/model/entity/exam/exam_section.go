// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamSection is the golang structure for table exam_section.
type ExamSection struct {
	Id                     int64       `json:"id"                        orm:"id"                        description:"主键"`                           // 主键
	ExamPaperId            int64       `json:"exam_paper_id"             orm:"exam_paper_id"             description:"试卷ID exam_paper.id"`           // 试卷ID exam_paper.id
	MockExaminationPaperId int64       `json:"mock_examination_paper_id" orm:"mock_examination_paper_id" description:"冗余 mock_examination_paper.id"` // 冗余 mock_examination_paper.id
	SortOrder              int         `json:"sort_order"                orm:"sort_order"                description:"在 index.items 中的顺序"`           // 在 index.items 中的顺序
	TopicTitle             string      `json:"topic_title"               orm:"topic_title"               description:"topic_title"`                  // topic_title
	TopicSubtitle          string      `json:"topic_subtitle"            orm:"topic_subtitle"            description:"topic_subtitle"`               // topic_subtitle
	TopicType              string      `json:"topic_type"                orm:"topic_type"                description:"题型代码 pt/xp/xt/..."`            // 题型代码 pt/xp/xt/...
	PartCode               int         `json:"part_code"                 orm:"part_code"                 description:"大题内 part 序号"`                  // 大题内 part 序号
	SegmentCode            string      `json:"segment_code"              orm:"segment_code"              description:"listen/read"`                  // listen/read
	TopicItemsFile         string      `json:"topic_items_file"          orm:"topic_items_file"          description:"topic_items 文件名，如 pt.json"`    // topic_items 文件名，如 pt.json
	TopicJson              string      `json:"topic_json"                orm:"topic_json"                description:"该 topic 文件全文快照"`               // 该 topic 文件全文快照
	Creator                string      `json:"creator"                   orm:"creator"                   description:"创建者"`                          // 创建者
	CreateTime             *gtime.Time `json:"create_time"               orm:"create_time"               description:"创建时间"`                         // 创建时间
	Updater                string      `json:"updater"                   orm:"updater"                   description:"更新者"`                          // 更新者
	UpdateTime             *gtime.Time `json:"update_time"               orm:"update_time"               description:"更新时间"`                         // 更新时间
	DeleteFlag             int         `json:"delete_flag"               orm:"delete_flag"               description:"逻辑删除"`                         // 逻辑删除
}
