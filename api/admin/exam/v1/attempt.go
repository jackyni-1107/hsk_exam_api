package v1

import "github.com/gogf/gf/v2/frame/g"

type AttemptListReq struct {
	g.Meta      `path:"/exam/attempt/list" method:"get" tags:"考试结果" summary:"答题会话列表"`
	Page        int    `json:"page" dc:"页码"`
	Size        int    `json:"size" dc:"每页条数"`
	Level       string `json:"level" dc:"试卷级别，如 hsk1"`
	ExaminationPaperId int64  `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	Status      int    `json:"status" dc:"会话状态 1-4，0 表示不限"`
	Username    string `json:"username" dc:"学员账号（模糊）"`
}

type AttemptListRes struct {
	List  []*AttemptListItem `json:"list"`
	Total int                `json:"total"`
}

type AttemptListItem struct {
	Id              int64   `json:"id"`
	ClientUserId    int64   `json:"client_user_id"`
	Username        string  `json:"username"`
	Nickname        string  `json:"nickname"`
	ExaminationPaperId int64 `json:"examination_paper_id"`
	PaperTitle      string  `json:"paper_title"`
	PaperLevel      string  `json:"paper_level"`
	RemotePaperId   string  `json:"remote_paper_id"`
	Status          int     `json:"status"`
	ObjectiveScore  float64 `json:"objective_score"`
	SubjectiveScore float64 `json:"subjective_score"`
	TotalScore      float64 `json:"total_score"`
	HasSubjective   int     `json:"has_subjective"`
	StartedAt       string  `json:"started_at"`
	SubmittedAt     string  `json:"submitted_at"`
	EndedAt         string  `json:"ended_at"`
	CreateTime      string  `json:"create_time"`
}

type AttemptDetailReq struct {
	g.Meta `path:"/exam/attempt/{id}" method:"get" tags:"考试结果" summary:"答题会话详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
}

type AttemptDetailRes struct {
	Attempt AttemptDetailAttempt  `json:"attempt"`
	User    AttemptDetailUser     `json:"user"`
	Paper   AttemptDetailPaper    `json:"paper"`
	Answers []AttemptDetailAnswer `json:"answers"`
}

type AttemptDetailOption struct {
	Id         int64  `json:"id"`
	Flag       string `json:"flag"`
	SortOrder  int    `json:"sort_order"`
	IsCorrect  int    `json:"is_correct"`
	OptionType string `json:"option_type"`
	Content    string `json:"content"`
}

type AttemptDetailAttempt struct {
	Id              int64   `json:"id"`
	ClientUserId    int64   `json:"client_user_id"`
	ExaminationPaperId int64 `json:"examination_paper_id"`
	Status          int     `json:"status"`
	DurationSeconds int     `json:"duration_seconds"`
	ObjectiveScore  float64 `json:"objective_score"`
	SubjectiveScore float64 `json:"subjective_score"`
	TotalScore      float64 `json:"total_score"`
	HasSubjective   int     `json:"has_subjective"`
	StartedAt       string  `json:"started_at"`
	DeadlineAt      string  `json:"deadline_at"`
	SubmittedAt     string  `json:"submitted_at"`
	EndedAt         string  `json:"ended_at"`
	CreateTime      string  `json:"create_time"`
}

type AttemptDetailUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type AttemptDetailPaper struct {
	Id      int64  `json:"id" dc:"mock_examination_paper.id"`
	Level   string `json:"level"`
	PaperId string `json:"paper_id"`
	Title   string `json:"title"`
}

type AttemptDetailAnswer struct {
	QuestionId       int64                 `json:"question_id"`
	QuestionNo       int                   `json:"question_no"`
	StemText         string                `json:"stem_text"`
	IsExample        int                   `json:"is_example"`
	IsSubjective     int                   `json:"is_subjective"`
	Score            float64               `json:"score"`
	AnswerJson       string                `json:"answer_json"`
	AwardedScore     *float64              `json:"awarded_score" dc:"主观题人工得分，未评为 null"`
	ObjectiveCorrect *bool                 `json:"objective_correct" dc:"客观非例题时是否选对，主观题/例题为 null"`
	SectionId        int64                 `json:"section_id" dc:"所属大题 ID"`
	SectionTitle     string                `json:"section_title" dc:"所属大题标题"`
	Options          []AttemptDetailOption `json:"options" dc:"题目选项列表（含是否正确标记）"`
}

type AttemptSubjectiveScoresReq struct {
	g.Meta `path:"/exam/attempt/{id}/subjective-scores" method:"put" tags:"考试结果" summary:"保存主观题得分"`
	Id     int64                        `json:"id" in:"path" v:"required|min:1#err.invalid_params"`
	Items  []AttemptSubjectiveScoreItem `json:"items" v:"required#err.invalid_params"`
}

type AttemptSubjectiveScoreItem struct {
	QuestionId int64   `json:"question_id" v:"required|min:1"`
	Score      float64 `json:"score" v:"min:0"`
}

type AttemptSubjectiveScoresRes struct {
	SubjectiveScore float64 `json:"subjective_score"`
	TotalScore      float64 `json:"total_score"`
}
