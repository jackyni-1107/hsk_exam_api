package mock

import (
	"context"

	mockentity "exam/internal/model/entity/mock"
)

type MockExaminationPaperWithImport struct {
	mockentity.MockExaminationPaper
	Imported bool `json:"imported"`
}

type IMock interface {
	MockLevelsList(ctx context.Context) ([]mockentity.MockLevels, error)
	ExaminationPaperList(ctx context.Context, levelId int64, importStatus string) ([]*MockExaminationPaperWithImport, error)
	ExaminationPaperDetail(ctx context.Context, id int64) (*MockExaminationPaperWithImport, error)
}

var localMock IMock

func Mock() IMock {
	if localMock == nil {
		panic("implement not found for interface IMock, forgot register?")
	}
	return localMock
}

func RegisterMock(i IMock) {
	localMock = i
}
