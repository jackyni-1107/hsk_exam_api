package bo

import "github.com/gogf/gf/v2/os/gtime"

// ExamBatchAdminItem 管理端批次列表/详情行。
type ExamBatchAdminItem struct {
	Id                      int64       `json:"id"`
	MockExaminationPaperIds []int64     `json:"mock_examination_paper_ids"`
	Title                   string      `json:"title"`
	ExamStartAt             *gtime.Time `json:"exam_start_at"`
	ExamEndAt               *gtime.Time `json:"exam_end_at"`
	MemberCount             int         `json:"member_count"`
	CreateTime              *gtime.Time `json:"create_time"`
}

// ExamBatchMemberAdminRow 批次内已导入学员一行。
type ExamBatchMemberAdminRow struct {
	MemberId               int64       `json:"member_id"`
	MockExaminationPaperId int64       `json:"mock_examination_paper_id"`
	Username               string      `json:"username"`
	Nickname               string      `json:"nickname"`
	ImportTime             *gtime.Time `json:"import_time"`
}
