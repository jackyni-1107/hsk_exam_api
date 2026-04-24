package bo

import "github.com/gogf/gf/v2/os/gtime"

// ExamBatchAdminItem 管理端批次列表/详情行。
type ExamBatchAdminItem struct {
	Id                    int64       `json:"id"`
	ExamPaperIds          []int64     `json:"exam_paper_ids"`
	Title                 string      `json:"title"`
	ExamStartAt           *gtime.Time `json:"exam_start_at"`
	ExamEndAt             *gtime.Time `json:"exam_end_at"`
	BatchKind             int         `json:"batch_kind"`
	AllowMultipleAttempts int         `json:"allow_multiple_attempts"`
	MaxAttemptsPerMember  int         `json:"max_attempts_per_member"`
	SkipScoring           int         `json:"skip_scoring"`
	AutoSubmitOnDeadline  int         `json:"auto_submit_on_deadline"`
	MemberCount           int         `json:"member_count"`
	CreateTime            *gtime.Time `json:"create_time"`
}

// ExamBatchPolicyInput 创建/更新批次时的策略字段（与 exam_batch 列一致）。
type ExamBatchPolicyInput struct {
	BatchKind             int
	AllowMultipleAttempts int
	MaxAttemptsPerMember  int
	SkipScoring           int
	AutoSubmitOnDeadline  int
}

// ExamBatchMemberAdminRow 批次内已导入学员一行。
type ExamBatchMemberAdminRow struct {
	MemberId    int64       `json:"member_id"`
	ExamPaperId int64       `json:"exam_paper_id"`
	Username    string      `json:"username"`
	Nickname    string      `json:"nickname"`
	ImportTime  *gtime.Time `json:"import_time"`
}
