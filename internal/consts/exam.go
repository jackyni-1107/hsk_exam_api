package consts

const (
	ExamImportConflictFail      = "fail"
	ExamImportConflictOverwrite = "overwrite"
	// ExamImportConflictNew 与 Overwrite 实现相同：不删 exam_paper 行、不迁移会话表（见 docs/exam-paper-import.md）。
	ExamImportConflictNew = "new"
)
