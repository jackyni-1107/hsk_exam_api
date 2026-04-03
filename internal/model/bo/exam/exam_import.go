package exam

type ImportParams struct {
	MockExaminationPaperId int64
	IndexURL               string
	IndexJSON              string
	Level                  string
	PaperID                string
	SourceBaseURL          string
	AudioHlsPrefix         string
	ConflictMode           string
	NewPaperID             string
	Creator                string
}

type ImportResult struct {
	ExaminationPaperID         int64
	Conflict                   bool
	ExistingExaminationPaperID int64
	SectionCount               int
	QuestionCount              int
}
