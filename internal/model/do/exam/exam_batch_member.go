// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamBatchMember is the golang structure of table exam_batch_member for DAO operations like Where/Data.
type ExamBatchMember struct {
	g.Meta                 `orm:"table:exam_batch_member, do:true"`
	BatchId                any         // exam_batch.id
	MemberId               any         // sys_member.id
	MockExaminationPaperId any         // mock_examination_paper.id
	ExamPaperId            any         //
	Creator                any         // 导入操作者
	CreateTime             *gtime.Time // 导入时间
}
