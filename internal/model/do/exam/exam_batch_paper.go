// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ExamBatchPaper is the golang structure of table exam_batch_paper for DAO operations like Where/Data.
type ExamBatchPaper struct {
	g.Meta                 `orm:"table:exam_batch_paper, do:true"`
	BatchId                any // exam_batch.id
	MockExaminationPaperId any // mock_examination_paper.id
	ExamPaperId            any //
}
