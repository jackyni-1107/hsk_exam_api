package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type BatchListReq struct {
	g.Meta      `path:"/exam/batch/list" method:"get" tags:"考试批次" summary:"批次分页列表" permission:"exam:batch:list"`
	ExamPaperId int64  `json:"exam_paper_id" dc:"按 exam_paper.id 筛选，0 表示不限"`
	TimeFrom    string `json:"time_from" dc:"可选，与 time_to 组成查询区间；筛选与批次 [exam_start_at, exam_end_at] 有交集的批次（可只传一侧）"`
	TimeTo      string `json:"time_to" dc:"可选，见 time_from；同时传时须 time_from <= time_to"`
	Page        int    `json:"page" dc:"页码"`
	Size        int    `json:"size" dc:"每页条数"`
}

type BatchListRes struct {
	List  []*BatchListItem `json:"list" dc:"列表"`
	Total int              `json:"total" dc:"总数"`
}

type BatchListItem struct {
	Id                    int64   `json:"id" dc:"批次ID"`
	ExamPaperIds          []int64 `json:"exam_paper_ids" dc:"关联的 exam_paper.id 列表"`
	Title                 string  `json:"title" dc:"批次名称"`
	ExamStartAt           string  `json:"exam_start_at" dc:"考试开始时间"`
	ExamEndAt             string  `json:"exam_end_at" dc:"考试结束时间"`
	BatchKind             int     `json:"batch_kind" dc:"0=正式 1=练习/模拟"`
	AllowMultipleAttempts int     `json:"allow_multiple_attempts" dc:"1=同用户同卷可多次新建会话"`
	MaxAttemptsPerMember  int     `json:"max_attempts_per_member" dc:"可重复时每人每卷上限，0=不限制"`
	SkipScoring           int     `json:"skip_scoring" dc:"1=交卷后不写入正式成绩"`
	AutoSubmitOnDeadline  int     `json:"auto_submit_on_deadline" dc:"0=不因个人考试倒计时自动交卷或拦截保存"`
	MemberCount           int     `json:"member_count" dc:"学员数"`
	CreateTime            string  `json:"create_time" dc:"创建时间"`
}

type BatchDetailReq struct {
	g.Meta `path:"/exam/batch/{id}" method:"get" tags:"考试批次" summary:"批次详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
}

type BatchDetailRes struct {
	Batch BatchListItem `json:"batch" dc:"批次信息"`
}

type BatchCreateReq struct {
	g.Meta                `path:"/exam/batch" method:"post" tags:"考试批次" summary:"创建批次" permission:"exam:batch:create"`
	Title                 string  `json:"title" dc:"批次名称"`
	ExamStartAt           string  `json:"exam_start_at" v:"required#err.invalid_params" dc:"考试允许开始时间，RFC3339 或常见日期时间字符串"`
	ExamEndAt             string  `json:"exam_end_at" v:"required#err.invalid_params" dc:"考试允许结束时间"`
	ExamPaperIds          []int64 `json:"exam_paper_ids" v:"required#err.invalid_params" dc:"exam_paper.id 多选"`
	BatchKind             int     `json:"batch_kind" dc:"0=正式 1=练习，默认 0"`
	AllowMultipleAttempts int     `json:"allow_multiple_attempts" dc:"1=允许多次新建答题会话，默认 0"`
	MaxAttemptsPerMember  int     `json:"max_attempts_per_member" dc:"可重复时每人每卷上限，0=不限制"`
	SkipScoring           int     `json:"skip_scoring" dc:"1=跳过正式算分与成绩写入，默认 0"`
	AutoSubmitOnDeadline  *int    `json:"auto_submit_on_deadline,omitempty" dc:"不传则默认 1；0=不因个人 deadline 自动交卷或拦截保存"`
}

type BatchCreateRes struct {
	Id int64 `json:"id" dc:"批次ID"`
}

type BatchUpdateReq struct {
	g.Meta                `path:"/exam/batch/{id}" method:"put" tags:"考试批次" summary:"更新批次" permission:"exam:batch:update"`
	Id                    int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	Title                 string  `json:"title" dc:"批次名称"`
	ExamStartAt           string  `json:"exam_start_at" v:"required#err.invalid_params" dc:"考试允许开始时间"`
	ExamEndAt             string  `json:"exam_end_at" v:"required#err.invalid_params" dc:"考试允许结束时间"`
	ExamPaperIds          []int64 `json:"exam_paper_ids" v:"required#err.invalid_params" dc:"exam_paper.id 列表"`
	BatchKind             int     `json:"batch_kind" dc:"0=正式 1=练习"`
	AllowMultipleAttempts int     `json:"allow_multiple_attempts" dc:"1=允许多次新建答题会话"`
	MaxAttemptsPerMember  int     `json:"max_attempts_per_member" dc:"可重复时每人每卷上限，0=不限制"`
	SkipScoring           int     `json:"skip_scoring" dc:"1=跳过正式算分"`
	AutoSubmitOnDeadline  *int    `json:"auto_submit_on_deadline,omitempty" dc:"不传则保持库中值；0=不因个人 deadline 自动交卷或拦截保存"`
}

type BatchUpdateRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/exam/batch/{id}" method:"delete" tags:"考试批次" summary:"删除批次" permission:"exam:batch:delete"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
}

type BatchDeleteRes struct{}

type BatchMembersImportReq struct {
	g.Meta      `path:"/exam/batch/{id}/members/import" method:"post" tags:"考试批次" summary:"向批次导入学员（sys_member.id）"`
	Id          int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	ExamPaperId int64   `json:"exam_paper_id" v:"required|min:1#err.invalid_params" dc:"须属于本批次已配置的 exam_paper"`
	MemberIds   []int64 `json:"member_ids" v:"required#err.invalid_params" dc:"学员主键列表，重复会自动去重"`
}

type BatchMembersImportRes struct {
	Inserted int `json:"inserted" dc:"新写入条数（已存在的主键不计入）"`
}

type BatchMembersImportFileReq struct {
	g.Meta      `path:"/exam/batch/{id}/members/import-file" method:"post" tags:"考试批次" summary:"导入学员" permission:"exam:batch:import-file"`
	Id          int64             `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	ExamPaperId int64             `json:"exam_paper_id" v:"required|min:1#err.invalid_params" dc:"须属于本批次已配置的 exam_paper"`
	File        *ghttp.UploadFile `json:"file" type:"file" dc:"CSV 文件（支持列：用户名/username 或 会员ID/member_id）"`
}

type BatchMembersImportFileRes struct {
	Total    int      `json:"total" dc:"CSV 非空数据行数"`
	Success  int      `json:"success" dc:"成功导入条数"`
	Failed   int      `json:"failed" dc:"失败条数"`
	Inserted int      `json:"inserted" dc:"实际新写入条数"`
	Errors   []string `json:"errors" dc:"错误明细（最多 100 条）"`
}

type BatchMemberListReq struct {
	g.Meta `path:"/exam/batch/{id}/members/list" method:"get" tags:"考试批次" summary:"批次内学员分页列表"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	Page   int   `json:"page" dc:"页码"`
	Size   int   `json:"size" dc:"每页条数"`
}

type BatchMemberListRes struct {
	List  []*BatchMemberListItem `json:"list" dc:"列表"`
	Total int                    `json:"total" dc:"总数"`
}

type BatchMemberListItem struct {
	MemberId    int64  `json:"member_id" dc:"学员ID"`
	ExamPaperId int64  `json:"exam_paper_id" dc:"exam_paper.id"`
	Username    string `json:"username" dc:"学员账号"`
	Nickname    string `json:"nickname" dc:"学员昵称"`
	ImportTime  string `json:"import_time" dc:"导入时间"`
}

type BatchMembersRemoveReq struct {
	g.Meta      `path:"/exam/batch/{id}/members/remove" method:"post" tags:"考试批次" summary:"从批次移除学员" permission:"exam:batch:remove"`
	Id          int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	ExamPaperId int64   `json:"exam_paper_id" v:"required|min:1#err.invalid_params" dc:"exam_paper.id"`
	MemberIds   []int64 `json:"member_ids" v:"required#err.invalid_params" dc:"学员主键列表"`
}

type BatchMembersRemoveRes struct {
	Removed int `json:"removed" dc:"移除条数"`
}
