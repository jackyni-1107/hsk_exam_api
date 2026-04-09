package exam

import (
	"exam/internal/dao/exam/internal"
)

type examBatchMockPaperDao struct {
	*internal.ExamBatchMockPaperDao
}

var (
	ExamBatchMockPaper = examBatchMockPaperDao{internal.NewExamBatchMockPaperDao()}
)
