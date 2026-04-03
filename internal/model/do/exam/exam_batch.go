// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamBatch is the golang structure of table exam_batch for DAO operations like Where/Data.
type ExamBatch struct {
	g.Meta                 `orm:"table:exam_batch, do:true"`
	Id                     any         // 主键
	MockExaminationPaperId any         // mock 卷 id，与 exam_paper 业务主键一致
	Title                  any         // 批次名称
	ExamStartAt            *gtime.Time // 考试允许开始时间
	ExamEndAt              *gtime.Time // 考试允许结束时间
	Creator                any         // 创建者
	CreateTime             *gtime.Time // 创建时间
	Updater                any         // 更新者
	UpdateTime             *gtime.Time // 更新时间
	DeleteFlag             any         // 逻辑删除：0-否，1-是
}
