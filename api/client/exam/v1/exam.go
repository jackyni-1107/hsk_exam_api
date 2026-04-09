package v1

import "github.com/gogf/gf/v2/frame/g"

// --- 试卷（脱敏） ---

type PaperForExamReq struct {
	g.Meta  `path:"/exam/papers/{paperId}" method:"get" tags:"客户端-考试" summary:"考前初始化试卷结构（无题目/选项，无标准答案）"`
	PaperId int64 `json:"paperId" in:"path" v:"required|min:1" dc:"mock_examination_paper.id"`
}

type PaperForExamRes struct {
	Id                   int64                  `json:"id" dc:"mock_examination_paper.id"`
	Level                string                 `json:"level"`
	PaperId              string                 `json:"paper_id"`
	Title                string                 `json:"title"`
	SourceBaseUrl        string                 `json:"source_base_url"`
	ListeningAudioPrefix string                 `json:"listening_audio_prefix"`
	DurationSeconds      int                    `json:"duration_seconds"`
	Prepare              PaperForExamPrepare    `json:"prepare"`
	Items                []PaperForExamItemInit `json:"items"`
}

type PaperSectionForExamReq struct {
	g.Meta    `path:"/exam/papers/{paperId}/sections/{sectionId}" method:"get" tags:"客户端-考试" summary:"按大题拉取 topic JSON 结构（脱敏，与资源文件字段一致）"`
	PaperId   int64 `json:"paperId" in:"path" v:"required|min:1" dc:"mock_examination_paper.id"`
	SectionId int64 `json:"sectionId" in:"path" v:"required|min:1"`
}

type PaperForExamPrepare struct {
	Instruction string `json:"instruction"`
	AudioFile   string `json:"audio_file"`
	Title       string `json:"title"`
}

// PaperForExamItemInit 初始化用 item 概要（block 无 questions），字段命名与 index.json 对齐。
type PaperForExamItemInit struct {
	Id            int64                   `json:"id"`
	SortOrder     int                     `json:"sort_order"`
	TopicTitle    string                  `json:"topic_title"`
	TopicSubtitle string                  `json:"topic_subtitle"`
	TopicType     string                  `json:"topic_type"`
	PartCode      int                     `json:"part_code"`
	SegmentCode   string                  `json:"segment_code"`
	TopicItems    string                  `json:"topic_items"`
	TopicJson     string                  `json:"topic_json"`
	Blocks        []PaperForExamBlockInit `json:"blocks"`
}

type PaperForExamBlockInit struct {
	Id                      int64  `json:"id"`
	BlockOrder              int    `json:"block_order"`
	GroupIndex              int    `json:"group_index"`
	QuestionDescriptionJson string `json:"question_description_json"`
	QuestionCount           int    `json:"question_count"`
}

// --- 会话 ---
//
//type AttemptCreateReq struct {
//	g.Meta  `path:"/exam/papers/{paperId}/attempts" method:"post" tags:"客户端-考试" summary:"已废弃：请使用 POST /exam/batches/{batchId}/attempts"`
//	PaperId int64 `json:"paperId" in:"path" v:"required|min:1" dc:"mock_examination_paper.id"`
//}

type AttemptCreateRes struct {
	AttemptId int64 `json:"attempt_id"`
}

// AttemptCreateByBatchReq 按考试批次与报名等级创建答题会话（未开始）。
type AttemptCreateByBatchReq struct {
	g.Meta                 `path:"/exam/batches/{batchId}/attempts" method:"post" tags:"客户端-考试" summary:"按批次与 Mock 卷创建答题会话"`
	BatchId                int64 `json:"batchId" in:"path" v:"required|min:1" dc:"exam_batch.id"`
	MockExaminationPaperId int64 `json:"mock_examination_paper_id" v:"required|min:1" dc:"须与 exam_batch_member 绑定一致"`
}

type AttemptStartReq struct {
	g.Meta          `path:"/exam/attempts/{id}/start" method:"post" tags:"客户端-考试" summary:"开考"`
	Id              int64 `json:"id" in:"path" v:"required|min:1"`
	DurationSeconds int   `json:"duration_seconds" dc:"可选，覆盖默认时长，服务端会按 max 夹紧"`
}

type AttemptStartRes struct{}

type AttemptGetReq struct {
	g.Meta `path:"/exam/attempts/{id}" method:"get" tags:"客户端-考试" summary:"答题会话详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1"`
}

type AttemptGetRes struct {
	Id                 int64   `json:"id"`
	ExaminationPaperId int64   `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	Status             int     `json:"status"`
	DurationSeconds    int     `json:"duration_seconds"`
	StartedAt          *string `json:"started_at"`
	DeadlineAt         *string `json:"deadline_at"`
	SubmittedAt        *string `json:"submitted_at"`
	EndedAt            *string `json:"ended_at"`
	ObjectiveScore     float64 `json:"objective_score"`
	SubjectiveScore    float64 `json:"subjective_score"`
	TotalScore         float64 `json:"total_score"`
	HasSubjective      int     `json:"has_subjective"`
	ServerTime         string  `json:"server_time"`
	DeadlineReached    bool    `json:"deadline_reached"`
}

type AttemptSaveAnswersReq struct {
	g.Meta `path:"/exam/attempts/{id}/answers" method:"put" tags:"客户端-考试" summary:"保存答案（批量）"`
	Id     int64               `json:"id" in:"path" v:"required|min:1"`
	Items  []AttemptAnswerItem `json:"items"`
}

type AttemptAnswerItem struct {
	QuestionId int64 `json:"question_id" v:"required"`
	Answer     any   `json:"answer" dc:"客观题传选项ID；写作题传文本"`
}

type AttemptSaveAnswersRes struct{}

type AttemptSubmitReq struct {
	g.Meta `path:"/exam/attempts/{id}/submit" method:"post" tags:"客户端-考试" summary:"交卷"`
	Id     int64 `json:"id" in:"path" v:"required|min:1"`
}

type AttemptSubmitRes struct{}

// AttemptRandomAnswersReq 测试专用：按试卷全部小题随机填答并保存（需配置 exam.enableRandomAnswerHelper）。
type AttemptRandomAnswersReq struct {
	g.Meta    `path:"/exam/papers/{paperId}/attempts/{attemptId}/random-answers" method:"post" tags:"客户端-考试-测试" summary:"随机填答（仅测试环境）"`
	PaperId   int64 `json:"paperId" in:"path" v:"required|min:1" dc:"mock_examination_paper.id"`
	AttemptId int64 `json:"attemptId" in:"path" v:"required|min:1"`
}

type AttemptSaveAnswersBody struct {
	Items []AttemptAnswerItem `json:"items"`
}

type AttemptRandomAnswersRes struct {
	GeneratedCount int                    `json:"generated_count" dc:"已生成题目数"`
	SubmitJSON     AttemptSaveAnswersBody `json:"submit_json" dc:"可直接作为 /api/client/exam/attempts/{id}/answers 的请求体"`
}

// // AudioHlsPlayIssueReq 签发 HLS 播放用短期票据（GET m3u8 无需 Bearer）。
// type AudioHlsPlayIssueReq struct {
// 	g.Meta     `path:"/exam/attempts/{id}/questions/{questionId}/audio/play" method:"post" tags:"客户端-考试" summary:"签发 HLS 短期播放地址"`
// 	Id         int64 `json:"id" in:"path" v:"required|min:1" dc:"答题会话 id"`
// 	QuestionId int64 `json:"questionId" in:"path" v:"required|min:1" dc:"小题 id"`
// }

// type AudioHlsPlayIssueRes struct {
// 	PlayUrl   string `json:"play_url" dc:"相对路径，需拼接到 API 根；指向原始 m3u8"`
// 	ExpiresAt string `json:"expires_at" dc:"票据过期时间 RFC3339 UTC"`
// }
