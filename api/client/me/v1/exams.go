package v1

import "github.com/gogf/gf/v2/frame/g"

type ExamsReq struct {
	g.Meta `path:"/me/exams" method:"get" tags:"客户端-当前用户" summary:"我的考试"`
}

type ExamsRes struct {
	List []ExamBatchItem `json:"list"`
}

type ExamBatchItem struct {
	BatchId                 int64   `json:"batch_id"`
	Title                   string  `json:"title" dc:"批次名称"`
	MockExaminationPaperId  int64   `json:"mock_examination_paper_id" dc:"本行绑定的 Mock 卷"`
	MockExaminationPaperIds []int64 `json:"mock_examination_paper_ids" dc:"批次配置的全部可选卷"`
	PaperTitle              string  `json:"paper_title"`
	ExamStartAt             string  `json:"exam_start_at"`
	ExamEndAt               string  `json:"exam_end_at"`
	WindowStatus            string  `json:"window_status" dc:"upcoming=未开始 open=进行中 closed=已结束或未开放"`
}
