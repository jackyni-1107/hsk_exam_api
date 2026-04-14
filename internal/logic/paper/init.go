package paper

import (
	"context"

	"exam/internal/service/paper"
)

type sPaper struct{}

func init() {
	paper.RegisterPaper(New())
}

func New() *sPaper {
	return &sPaper{}
}

func (s *sPaper) InvalidatePaperForExamCache(ctx context.Context, examPaperId int64) {
	InvalidatePaperForExamCache(ctx, examPaperId)
}

func (s *sPaper) InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64) {
	InvalidatePaperSectionForExamCache(ctx, examPaperId, sectionId)
}
