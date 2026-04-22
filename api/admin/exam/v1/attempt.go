package v1

import "github.com/gogf/gf/v2/frame/g"

type AttemptListReq struct {
	g.Meta             `path:"/exam/attempt/list" method:"get" tags:"考试结果" summary:"答题会话列表"`
	Page               int    `json:"page" dc:"页码"`
	Size               int    `json:"size" dc:"每页条数"`
	Level              string `json:"level" dc:"试卷级别，如 hsk1"`
	ExaminationPaperId int64  `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	ExamBatchId        int64  `json:"exam_batch_id" dc:"考试批次 id，0 不限"`
	Status             int    `json:"status" dc:"会话状态 1-4，0 表示不限"`
	Username           string `json:"username" dc:"学员账号（模糊）"`
}

type AttemptListRes struct {
	List  []*AttemptListItem `json:"list" dc:"列表"`
	Total int                `json:"total" dc:"总数"`
}

type AttemptListItem struct {
	Id                 int64   `json:"id" dc:"会话ID"`
	MemberId           int64   `json:"member_id" dc:"学员ID"`
	Username           string  `json:"username" dc:"学员账号"`
	Nickname           string  `json:"nickname" dc:"学员昵称"`
	ExaminationPaperId int64   `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	ExamBatchId        int64   `json:"exam_batch_id" dc:"考试批次ID"`
	MockLevelId        int64   `json:"mock_level_id" dc:"Mock等级ID"`
	PaperTitle         string  `json:"paper_title" dc:"试卷标题"`
	PaperLevel         string  `json:"paper_level" dc:"试卷级别"`
	RemotePaperId      string  `json:"remote_paper_id" dc:"远程试卷ID"`
	Status             int     `json:"status" dc:"会话状态"`
	ObjectiveScore     float64 `json:"objective_score" dc:"客观题得分"`
	SubjectiveScore    float64 `json:"subjective_score" dc:"主观题得分"`
	TotalScore         float64 `json:"total_score" dc:"总分"`
	HasSubjective      int     `json:"has_subjective" dc:"是否含主观题：0否 1是"`
	SubjectiveGraded   int     `json:"subjective_graded" dc:"主观题是否已评过（1=是，仅允许评一次）"`
	StartedAt          string  `json:"started_at" dc:"开考时间"`
	SubmittedAt        string  `json:"submitted_at" dc:"交卷时间"`
	EndedAt            string  `json:"ended_at" dc:"结束时间"`
	CreateTime         string  `json:"create_time" dc:"创建时间"`
}

type AttemptDetailReq struct {
	g.Meta `path:"/exam/attempt/{id}" method:"get" tags:"考试结果" summary:"答题会话详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"会话ID"`
}

type AttemptDetailRes struct {
	Attempt AttemptDetailAttempt  `json:"attempt" dc:"会话信息"`
	User    AttemptDetailUser     `json:"user" dc:"用户信息"`
	Paper   AttemptDetailPaper    `json:"paper" dc:"试卷信息"`
	Answers []AttemptDetailAnswer `json:"answers" dc:"答题列表"`
}

type AttemptDetailOption struct {
	Id         int64  `json:"id" dc:"选项ID"`
	Flag       string `json:"flag" dc:"选项标识"`
	SortOrder  int    `json:"sort_order" dc:"排序"`
	IsCorrect  int    `json:"is_correct" dc:"是否正确：0否 1是"`
	OptionType string `json:"option_type" dc:"选项类型"`
	Content    string `json:"content" dc:"选项内容"`
}

type AttemptDetailAttempt struct {
	Id                 int64   `json:"id" dc:"会话ID"`
	MemberId           int64   `json:"member_id" dc:"学员ID"`
	ExaminationPaperId int64   `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	Status             int     `json:"status" dc:"会话状态"`
	DurationSeconds    int     `json:"duration_seconds" dc:"考试时长(秒)"`
	ObjectiveScore     float64 `json:"objective_score" dc:"客观题得分"`
	SubjectiveScore    float64 `json:"subjective_score" dc:"主观题得分"`
	TotalScore         float64 `json:"total_score" dc:"总分"`
	HasSubjective      int     `json:"has_subjective" dc:"是否含主观题：0否 1是"`
	StartedAt          string  `json:"started_at" dc:"开考时间"`
	DeadlineAt         string  `json:"deadline_at" dc:"截止时间"`
	SubmittedAt        string  `json:"submitted_at" dc:"交卷时间"`
	EndedAt            string  `json:"ended_at" dc:"结束时间"`
	CreateTime         string  `json:"create_time" dc:"创建时间"`
}

type AttemptDetailUser struct {
	Id       int64  `json:"id" dc:"用户ID"`
	Username string `json:"username" dc:"用户名"`
	Nickname string `json:"nickname" dc:"昵称"`
}

type AttemptDetailPaper struct {
	Id            int64  `json:"id" dc:"mock_examination_paper.id"`
	Name          string `json:"name" dc:"mock_examination_paper.name"`
	Level         string `json:"level" dc:"试卷级别"`
	PaperId       string `json:"paper_id" dc:"远程试卷ID"`
	Title         string `json:"title" dc:"exam_paper.title"`
	ExamPaperId   int64  `json:"exam_paper_id" dc:"exam_paper.id"`
	SourceBaseUrl string `json:"source_base_url" dc:"资源基址（拼接题目包内相对路径）"`
}

type AttemptDetailAnswer struct {
	QuestionId       int64                 `json:"question_id" dc:"题目ID"`
	QuestionNo       int                   `json:"question_no" dc:"题号"`
	StemText         string                `json:"stem_text" dc:"题干文本"`
	ScreenTextJson   string                `json:"screen_text_json" dc:"屏幕文本 JSON（与试卷详情一致）"`
	IsExample        int                   `json:"is_example" dc:"是否例题：0否 1是"`
	IsSubjective     int                   `json:"is_subjective" dc:"是否主观题：0否 1是"`
	Score            float64               `json:"score" dc:"题目分值"`
	AnswerJson       string                `json:"answer_json" dc:"作答内容(JSON)"`
	AwardedScore     *float64              `json:"awarded_score" dc:"主观题人工得分，未评为 null"`
	ObjectiveCorrect *bool                 `json:"objective_correct" dc:"客观非例题时是否选对，主观题/例题为 null"`
	SectionId        int64                 `json:"section_id" dc:"所属大题 ID"`
	SectionTitle     string                `json:"section_title" dc:"所属大题标题"`
	AnalysisText     string                `json:"analysis_text" dc:"解析文案（从 analysis_json 抽取）"`
	Options          []AttemptDetailOption `json:"options" dc:"题目选项列表（含是否正确标记）"`
}

type AttemptSubjectiveScoresReq struct {
	g.Meta `path:"/exam/attempt/{id}/subjective-scores" method:"put" tags:"考试结果" summary:"保存主观题得分"`
	Id     int64                        `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"会话ID"`
	Items  []AttemptSubjectiveScoreItem `json:"items" v:"required#err.invalid_params" dc:"评分项列表"`
}

type AttemptSubjectiveScoreItem struct {
	QuestionId int64   `json:"question_id" v:"required|min:1" dc:"题目ID"`
	Score      float64 `json:"score" v:"min:0" dc:"得分"`
}

type AttemptSubjectiveScoresRes struct {
	SubjectiveScore float64 `json:"subjective_score" dc:"主观题总分"`
	TotalScore      float64 `json:"total_score" dc:"总分"`
}
