package bo

import "github.com/gogf/gf/v2/os/gtime"

// MyExamBatchItem 学员端「我的考试」列表一行（批次 + 试卷概要 + 时间窗状态）。
type MyExamBatchItem struct {
	BatchId                int64       `json:"batch_id"`
	Title                  string      `json:"title"`
	ExamPaperId            int64       `json:"exam_paper_id"`             // exam_batch_member.exam_paper_id
	MockExaminationPaperId int64       `json:"mock_examination_paper_id"` // 冗余：exam_paper.mock_examination_paper_id，兼容旧客户端
	PaperTitle             string      `json:"paper_title"`
	ExamStartAt            *gtime.Time `json:"exam_start_at"`
	ExamEndAt              *gtime.Time `json:"exam_end_at"`
	AttemptId              int64       `json:"attempt_id"`
	WindowStatus           string      `json:"window_status"` // upcoming | open | closed
}
