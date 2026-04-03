// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ExamBatchMockLevel is the golang structure of table exam_batch_mock_level for DAO operations like Where/Data.
type ExamBatchMockLevel struct {
	g.Meta      `orm:"table:exam_batch_mock_level, do:true"`
	BatchId     any // exam_batch.id
	MockLevelId any // mock_levels.id
}
