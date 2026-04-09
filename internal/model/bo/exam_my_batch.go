package bo

import "github.com/gogf/gf/v2/os/gtime"

// MyExamBatchItem 学员端「我的考试」列表一行（批次 + 试卷概要 + 时间窗状态）。
type MyExamBatchItem struct {
	BatchId                 int64       `json:"batch_id"`
	Title                   string      `json:"title"`
	MockExaminationPaperId  int64       `json:"mock_examination_paper_id"`  // 本行绑定的 Mock 卷 exam_batch_member.mock_examination_paper_id
	MockExaminationPaperIds []int64     `json:"mock_examination_paper_ids"` // 批次配置的全部可选卷
	PaperTitle              string      `json:"paper_title"`
	ExamStartAt             *gtime.Time `json:"exam_start_at"`
	ExamEndAt               *gtime.Time `json:"exam_end_at"`
	WindowStatus            string      `json:"window_status"` // upcoming | open | closed
}
