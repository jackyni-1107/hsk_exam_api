// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamSection is the golang structure of table exam_section for DAO operations like Where/Data.
type ExamSection struct {
	g.Meta         `orm:"table:exam_section, do:true"`
	Id             any         // 主键
	ExamPaperId            any         // 试卷ID exam_paper.id
	MockExaminationPaperId any         // 冗余 mock_examination_paper.id
	SortOrder              any         // 在 index.items 中的顺序
	TopicTitle     any         // topic_title
	TopicSubtitle  any         // topic_subtitle
	TopicType      any         // 题型代码 pt/xp/xt/...
	PartCode       any         // 大题内 part 序号
	SegmentCode    any         // listen/read
	TopicItemsFile any         // topic_items 文件名，如 pt.json
	TopicJson      any         // 该 topic 文件全文快照
	Creator        any         // 创建者
	CreateTime     *gtime.Time // 创建时间
	Updater        any         // 更新者
	UpdateTime     *gtime.Time // 更新时间
	DeleteFlag     any         // 逻辑删除
}
