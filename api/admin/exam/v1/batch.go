package v1

import "github.com/gogf/gf/v2/frame/g"

type BatchListReq struct {
	g.Meta                 `path:"/exam/batch/list" method:"get" tags:"考试批次" summary:"批次分页列表"`
	MockExaminationPaperId int64 `json:"mock_examination_paper_id" dc:"按 mock 卷 id 筛选，0 表示不限"`
	Page                   int   `json:"page" dc:"页码"`
	Size                   int   `json:"size" dc:"每页条数"`
}

type BatchListRes struct {
	List  []*BatchListItem `json:"list" dc:"列表"`
	Total int              `json:"total" dc:"总数"`
}

type BatchListItem struct {
	Id                      int64   `json:"id" dc:"批次ID"`
	MockExaminationPaperIds []int64 `json:"mock_examination_paper_ids" dc:"关联的 mock 卷 id 列表"`
	Title                   string  `json:"title" dc:"批次名称"`
	ExamStartAt             string  `json:"exam_start_at" dc:"考试开始时间"`
	ExamEndAt               string  `json:"exam_end_at" dc:"考试结束时间"`
	MemberCount             int     `json:"member_count" dc:"学员数"`
	CreateTime              string  `json:"create_time" dc:"创建时间"`
}

type BatchDetailReq struct {
	g.Meta `path:"/exam/batch/{id}" method:"get" tags:"考试批次" summary:"批次详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
}

type BatchDetailRes struct {
	Batch BatchListItem `json:"batch" dc:"批次信息"`
}

type BatchCreateReq struct {
	g.Meta                  `path:"/exam/batch" method:"post" tags:"考试批次" summary:"创建批次"`
	Title                   string  `json:"title" dc:"批次名称"`
	ExamStartAt             string  `json:"exam_start_at" v:"required#err.invalid_params" dc:"考试允许开始时间，RFC3339 或常见日期时间字符串"`
	ExamEndAt               string  `json:"exam_end_at" v:"required#err.invalid_params" dc:"考试允许结束时间"`
	MockExaminationPaperIds []int64 `json:"mock_examination_paper_ids" v:"required#err.invalid_params" dc:"mock 卷 id 多选；须已在考试侧导入 exam_paper"`
}

type BatchCreateRes struct {
	Id int64 `json:"id" dc:"批次ID"`
}

type BatchUpdateReq struct {
	g.Meta                  `path:"/exam/batch/{id}" method:"put" tags:"考试批次" summary:"更新批次"`
	Id                      int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	Title                   string  `json:"title" dc:"批次名称"`
	ExamStartAt             string  `json:"exam_start_at" v:"required#err.invalid_params" dc:"考试允许开始时间"`
	ExamEndAt               string  `json:"exam_end_at" v:"required#err.invalid_params" dc:"考试允许结束时间"`
	MockExaminationPaperIds []int64 `json:"mock_examination_paper_ids" v:"required#err.invalid_params" dc:"mock 卷 id 列表"`
}

type BatchUpdateRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/exam/batch/{id}" method:"delete" tags:"考试批次" summary:"删除批次（逻辑删除）"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
}

type BatchDeleteRes struct{}

type BatchMembersImportReq struct {
	g.Meta                 `path:"/exam/batch/{id}/members/import" method:"post" tags:"考试批次" summary:"向批次导入学员（sys_member.id）"`
	Id                     int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	MockExaminationPaperId int64   `json:"mock_examination_paper_id" v:"required|min:1#err.invalid_params" dc:"须属于本批次已配置的 mock 卷"`
	MemberIds              []int64 `json:"member_ids" v:"required#err.invalid_params" dc:"学员主键列表，重复会自动去重"`
}

type BatchMembersImportRes struct {
	Inserted int `json:"inserted" dc:"新写入条数（已存在的主键不计入）"`
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
	MemberId               int64  `json:"member_id" dc:"学员ID"`
	MockExaminationPaperId int64  `json:"mock_examination_paper_id" dc:"mock 卷 ID"`
	Username               string `json:"username" dc:"学员账号"`
	Nickname               string `json:"nickname" dc:"学员昵称"`
	ImportTime             string `json:"import_time" dc:"导入时间"`
}

type BatchMembersRemoveReq struct {
	g.Meta                 `path:"/exam/batch/{id}/members/remove" method:"post" tags:"考试批次" summary:"从批次移除学员（指定 Mock 卷）"`
	Id                     int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"批次ID"`
	MockExaminationPaperId int64   `json:"mock_examination_paper_id" v:"required|min:1#err.invalid_params" dc:"mock 卷 ID"`
	MemberIds              []int64 `json:"member_ids" v:"required#err.invalid_params" dc:"学员主键列表"`
}

type BatchMembersRemoveRes struct {
	Removed int `json:"removed" dc:"移除条数"`
}
