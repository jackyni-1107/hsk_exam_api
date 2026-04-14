package bo

import (
	"exam/internal/model/entity/mock"
)

// MockExaminationPaperWithImport 移动到这里
type MockExaminationPaperWithImport struct {
	mock.MockExaminationPaper
	Imported bool `json:"imported"`
}
