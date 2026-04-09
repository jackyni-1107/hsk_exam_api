// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamBatch is the golang structure for table exam_batch.
type ExamBatch struct {
	Id          int64       `json:"id"            orm:"id"            description:"主键"`           // 主键
	Title       string      `json:"title"         orm:"title"         description:"批次名称"`         // 批次名称
	ExamStartAt *gtime.Time `json:"exam_start_at" orm:"exam_start_at" description:"考试允许开始时间"`     // 考试允许开始时间
	ExamEndAt   *gtime.Time `json:"exam_end_at"   orm:"exam_end_at"   description:"考试允许结束时间"`     // 考试允许结束时间
	Creator     string      `json:"creator"       orm:"creator"       description:"创建者"`          // 创建者
	CreateTime  *gtime.Time `json:"create_time"   orm:"create_time"   description:"创建时间"`         // 创建时间
	Updater     string      `json:"updater"       orm:"updater"       description:"更新者"`          // 更新者
	UpdateTime  *gtime.Time `json:"update_time"   orm:"update_time"   description:"更新时间"`         // 更新时间
	DeleteFlag  int         `json:"delete_flag"   orm:"delete_flag"   description:"逻辑删除：0-否，1-是"` // 逻辑删除：0-否，1-是
}
