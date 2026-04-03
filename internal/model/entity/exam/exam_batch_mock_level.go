// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

// ExamBatchMockLevel is the golang structure for table exam_batch_mock_level.
type ExamBatchMockLevel struct {
	BatchId     int64 `json:"batch_id"      orm:"batch_id"      description:"exam_batch.id"`  // exam_batch.id
	MockLevelId int64 `json:"mock_level_id" orm:"mock_level_id" description:"mock_levels.id"` // mock_levels.id
}
