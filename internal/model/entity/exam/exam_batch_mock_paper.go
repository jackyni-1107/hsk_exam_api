package exam

// ExamBatchMockPaper is the golang structure for table exam_batch_mock_paper.
type ExamBatchMockPaper struct {
	BatchId                int64 `json:"batch_id"                  orm:"batch_id"                  description:"exam_batch.id"`             // exam_batch.id
	MockExaminationPaperId int64 `json:"mock_examination_paper_id" orm:"mock_examination_paper_id" description:"mock_examination_paper.id"` // mock_examination_paper.id
}
