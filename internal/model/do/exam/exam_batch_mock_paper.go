package exam

import "github.com/gogf/gf/v2/frame/g"

// ExamBatchMockPaper is the golang structure of table exam_batch_mock_paper for DAO operations like Where/Data.
type ExamBatchMockPaper struct {
	g.Meta                 `orm:"table:exam_batch_mock_paper, do:true"`
	BatchId                any // exam_batch.id
	MockExaminationPaperId any // mock_examination_paper.id
}
