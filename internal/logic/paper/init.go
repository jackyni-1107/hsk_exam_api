package paper

import (
	"context"

	"exam/internal/service/exam"
)

type sExam struct{}

func init() {
	exam.RegisterExam(New())
}

func New() *sExam {
	return &sExam{}
}

func (s *sExam) InvalidatePaperForExamCache(ctx context.Context, examPaperId int64) {
	InvalidatePaperForExamCache(ctx, examPaperId)
}

func (s *sExam) InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64) {
	InvalidatePaperSectionForExamCache(ctx, examPaperId, sectionId)
}
