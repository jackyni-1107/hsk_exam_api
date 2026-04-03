package v1

import "github.com/gogf/gf/v2/frame/g"

type BatchListReq struct {
	g.Meta                 `path:"/exam/batch/list" method:"get" tags:"考试批次" summary:"批次分页列表"`
	MockExaminationPaperId int64 `json:"mock_examination_paper_id" dc:"按 mock 卷 id 筛选，0 表示不限"`
	Page                   int   `json:"page" dc:"页码"`
	Size                   int   `json:"size" dc:"每页条数"`
}

type BatchListRes struct {
	List  []*BatchListItem `json:"list"`
	Total int              `json:"total"`
}

type BatchListItem struct {
	Id                     int64   `json:"id"`
	MockExaminationPaperId int64   `json:"mock_examination_paper_id"`
	Title                  string  `json:"title"`
	ExamStartAt            string  `json:"exam_start_at"`
	ExamEndAt              string  `json:"exam_end_at"`
	MockLevelIds           []int64 `json:"mock_level_ids"`
	MemberCount            int     `json:"member_count"`
	CreateTime             string  `json:"create_time"`
}

type BatchDetailReq struct {
	g.Meta `path:"/exam/batch/{id}" method:"get" tags:"考试批次" summary:"批次详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
}

type BatchDetailRes struct {
	Batch BatchListItem `json:"batch"`
}

type BatchCreateReq struct {
	g.Meta                 `path:"/exam/batch" method:"post" tags:"考试批次" summary:"创建批次"`
	MockExaminationPaperId int64   `json:"mock_examination_paper_id" v:"required|min:1#err.invalid_params"`
	Title                  string  `json:"title" dc:"批次名称"`
	ExamStartAt            string  `json:"exam_start_at" v:"required#err.invalid_params" dc:"考试允许开始时间，RFC3339 或常见日期时间字符串"`
	ExamEndAt              string  `json:"exam_end_at" v:"required#err.invalid_params" dc:"考试允许结束时间"`
	MockLevelIds           []int64 `json:"mock_level_ids" v:"required#err.invalid_params" dc:"mock_levels.id，可多选"`
}

type BatchCreateRes struct {
	Id int64 `json:"id"`
}

type BatchUpdateReq struct {
	g.Meta       `path:"/exam/batch/{id}" method:"put" tags:"考试批次" summary:"更新批次"`
	Id           int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
	Title        string  `json:"title"`
	ExamStartAt  string  `json:"exam_start_at" v:"required#err.invalid_params"`
	ExamEndAt    string  `json:"exam_end_at" v:"required#err.invalid_params"`
	MockLevelIds []int64 `json:"mock_level_ids" v:"required#err.invalid_params"`
}

type BatchUpdateRes struct{}

type BatchDeleteReq struct {
	g.Meta `path:"/exam/batch/{id}" method:"delete" tags:"考试批次" summary:"删除批次（逻辑删除）"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
}

type BatchDeleteRes struct{}

type BatchMembersImportReq struct {
	g.Meta    `path:"/exam/batch/{id}/members/import" method:"post" tags:"考试批次" summary:"向批次导入学员（sys_member.id）"`
	Id        int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
	MemberIds []int64 `json:"member_ids" v:"required#err.invalid_params" dc:"学员主键列表，重复会自动去重"`
}

type BatchMembersImportRes struct {
	Inserted int `json:"inserted" dc:"新写入条数（已存在的主键不计入）"`
}

type BatchMemberListReq struct {
	g.Meta `path:"/exam/batch/{id}/members/list" method:"get" tags:"考试批次" summary:"批次内学员分页列表"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
	Page   int   `json:"page"`
	Size   int   `json:"size"`
}

type BatchMemberListRes struct {
	List  []*BatchMemberListItem `json:"list"`
	Total int                    `json:"total"`
}

type BatchMemberListItem struct {
	MemberId   int64  `json:"member_id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ImportTime string `json:"import_time"`
}

type BatchMembersRemoveReq struct {
	g.Meta    `path:"/exam/batch/{id}/members/remove" method:"post" tags:"考试批次" summary:"从批次移除学员"`
	Id        int64   `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
	MemberIds []int64 `json:"member_ids" v:"required#err.invalid_params"`
}

type BatchMembersRemoveRes struct {
	Removed int `json:"removed"`
}
