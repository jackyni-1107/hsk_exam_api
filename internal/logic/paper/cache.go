package paper

import (
	"context"

	"exam/internal/utility/exampaper"
)

func invalidatePaperCaches(ctx context.Context, examPaperID, mockPaperID int64) {
	if examPaperID > 0 {
		InvalidatePaperForExamCache(ctx, examPaperID)
		invalidateAdminPaperDetailCacheByPaper(examPaperID)
	}
	if mockPaperID > 0 {
		exampaper.InvalidateByMockIDCache(mockPaperID)
	}
}
