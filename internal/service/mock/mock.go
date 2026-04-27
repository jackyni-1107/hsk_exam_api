// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package mock

import (
	"context"
	"exam/internal/model/bo"
	mockentity "exam/internal/model/entity/mock"
)

type (
	IMock interface {
		MockLevelsList(ctx context.Context) ([]mockentity.MockLevels, error)
		MockLevelNameByID(ctx context.Context, id int64) (string, error)
		ExaminationPaperList(ctx context.Context, levelId int64, importStatus string) ([]bo.MockExaminationPaperWithImport, error)
		ExaminationPaperDetail(ctx context.Context, id int64) (*bo.MockExaminationPaperWithImport, error)
	}
)

var (
	localMock IMock
)

func Mock() IMock {
	if localMock == nil {
		panic("implement not found for interface IMock, forgot register?")
	}
	return localMock
}

func RegisterMock(i IMock) {
	localMock = i
}
