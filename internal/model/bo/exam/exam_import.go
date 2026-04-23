package exam

type ImportParams struct {
	MockExaminationPaperId int64
	Title                  string
	AudioHlsPrefix         string
	ConflictMode           string
	Creator                string
}

type ImportResult struct {
	MockExaminationPaperID         int64
	Conflict                       bool
	ExistingMockExaminationPaperID int64
	SectionCount                   int
	QuestionCount                  int
}
