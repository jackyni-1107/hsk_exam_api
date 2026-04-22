package exam

type ImportParams struct {
	MockExaminationPaperId int64
	Title                  string
	AudioHlsPrefix         string
	ConflictMode           string
	Creator                string
}

type ImportResult struct {
	ExaminationPaperID         int64
	Conflict                   bool
	ExistingExaminationPaperID int64
	SectionCount               int
	QuestionCount              int
}
