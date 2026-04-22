// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

// ExamBatchPaper is the golang structure for table exam_batch_paper.
type ExamBatchPaper struct {
	BatchId                int64 `json:"batch_id"                  orm:"batch_id"                  description:"exam_batch.id"`             // exam_batch.id
	MockExaminationPaperId int64 `json:"mock_examination_paper_id" orm:"mock_examination_paper_id" description:"mock_examination_paper.id"` // mock_examination_paper.id
	ExamPaperId            int64 `json:"exam_paper_id"             orm:"exam_paper_id"             description:""`                          //
}
