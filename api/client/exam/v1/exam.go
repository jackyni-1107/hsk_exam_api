package v1

import "github.com/gogf/gf/v2/frame/g"

// --- 试卷（脱敏） ---

type PaperForExamReq struct {
	g.Meta  `path:"/exam/papers/{paperId}" method:"get" tags:"客户端-考试" summary:"考前初始化试卷结构（无题目/选项，无标准答案）"`
	PaperId int64 `json:"paperId" in:"path" v:"required|min:1" dc:"exam_paper.id"`
}

type PaperForExamRes struct {
	Id                   int64                  `json:"id" dc:"exam_paper.id（主键）"`
	MockExaminationPaperId int64              `json:"mock_examination_paper_id,omitempty" dc:"冗余：mock_examination_paper.id"`
	Level                string                 `json:"level" dc:"试卷级别"`
	PaperId              string                 `json:"paper_id" dc:"远程试卷ID"`
	Title                string                 `json:"title" dc:"试卷标题"`
	SourceBaseUrl        string                 `json:"source_base_url" dc:"资源基础URL"`
	AudioUrl             string                 `json:"audio_url" dc:"m3u8链接"`
	DurationSeconds      int                    `json:"duration_seconds" dc:"考试时长(秒)"`
	ListenReviewDuration int                    `json:"listen_review_duration" dc:"听力结束后回顾时间(秒)"`
	Prepare              PaperForExamPrepare    `json:"prepare" dc:"准备阶段信息"`
	Items                []PaperForExamItemInit `json:"items" dc:"大题初始化列表"`
}

type PaperSectionForExamReq struct {
	g.Meta    `path:"/exam/papers/{paperId}/sections/{sectionId}" method:"get" tags:"客户端-考试" summary:"按大题拉取 topic JSON 结构（脱敏，与资源文件字段一致）"`
	PaperId   int64 `json:"paperId" in:"path" v:"required|min:1" dc:"exam_paper.id"`
	SectionId int64 `json:"sectionId" in:"path" v:"required|min:1" dc:"大题ID"`
}

type PaperForExamPrepare struct {
	Instruction string                `json:"instruction" dc:"准备阶段说明"`
	AudioFile   string                `json:"audio_file" dc:"准备阶段音频文件"`
	Title       string                `json:"title" dc:"准备阶段标题"`
	Name        string                `json:"name" dc:"试卷标题"`
	Segments    []PaperForExamSegment `json:"segments" dc:"考试环节信息"`
}

type PaperForExamSegment struct {
	SegmentCode   string                    `json:"segmentCode" dc:"环节编码"`
	TotalScore    int                       `json:"totalScore" dc:"环节总分"`
	QuestionCount int                       `json:"questionCount" dc:"题目数量"`
	Duration      int                       `json:"duration" dc:"环节时长(分钟)"`
	Seq           int                       `json:"seq" dc:"顺序"`
	SegmentDesc   string                    `json:"segmentDesc" dc:"环节描述"`
	Parts         []PaperForExamSegmentPart `json:"parts" dc:"环节下分项"`
}

type PaperForExamSegmentPart struct {
	PartCode                int     `json:"partCode" dc:"部分编号"`
	PartName                string  `json:"part_name" dc:"部分名称"`
	PartNameTrans           string  `json:"part_name_trans" dc:"部分名称国际化"`
	PartRate                float64 `json:"part_rate" dc:"部分分值占比"`
	PartScore               float64 `json:"part_score" dc:"部分分值"`
	QuestionCount           int     `json:"questionCount" dc:"题目数量"`
	ObjectiveQuestionCount  int     `json:"objectiveQuestionCount" dc:"客观题数量"`
	SubjectiveQuestionCount int     `json:"subjectiveQuestionCount" dc:"主观题数量"`
	PartAnswerTime          int     `json:"part_answer_time" dc:"部分答题时长"`
	ScoreTotal              float64 `json:"score_total" dc:"已得分"`
	CorrectCount            int     `json:"correct_count" dc:"答对数量"`
	CorrectRate             float64 `json:"correct_rate" dc:"正确率"`
	Practiced               bool    `json:"practiced" dc:"是否练习过"`
	QuestionType            string  `json:"question_type" dc:"题型"`
}

// PaperForExamItemInit 初始化用 item 概要（block 无 questions），字段命名与 index.json 对齐。
type PaperForExamItemInit struct {
	Id            int64  `json:"id" dc:"大题ID"`
	SortOrder     int    `json:"sort_order" dc:"排序"`
	TopicTitle    string `json:"topic_title" dc:"大题标题"`
	TopicSubtitle string `json:"topic_subtitle" dc:"大题副标题"`
	TopicType     string `json:"topic_type" dc:"题型"`
	PartCode      int    `json:"part_code" dc:"部分编号"`
	SegmentCode   string `json:"segment_code" dc:"段落编号"`
	TopicItems    string `json:"topic_items" dc:"题目文件名"`
	//TopicJson     string                  `json:"topic_json" dc:"大题JSON"`
	//Blocks        []PaperForExamBlockInit `json:"blocks" dc:"题块初始化列表"`
}

type PaperForExamBlockInit struct {
	Id                      int64  `json:"id" dc:"题块ID"`
	BlockOrder              int    `json:"block_order" dc:"题块排序"`
	GroupIndex              int    `json:"group_index" dc:"组索引"`
	QuestionDescriptionJson string `json:"question_description_json" dc:"题块描述(JSON)"`
	QuestionCount           int    `json:"question_count" dc:"题目数量"`
}

type AttemptCreateRes struct {
	AttemptId int64 `json:"attempt_id" dc:"答题会话ID"`
}

// AttemptCreateByBatchReq 按批次创建答题会话（未开始）；多卷批次须在 body 传 exam_paper_id。
type AttemptCreateByBatchReq struct {
	g.Meta        `path:"/exam/batches/{batchId}/attempts" method:"post" tags:"客户端-考试" summary:"按批次创建答题会话"`
	BatchId       int64 `json:"batchId" in:"path" v:"required|min:1" dc:"exam_batch.id"`
	ExamPaperId   int64 `json:"exam_paper_id" dc:"用户在批次中绑定的 exam_paper.id；该用户在批次下有多条 member 绑定（多卷）时必填"`
}

type AttemptStartReq struct {
	g.Meta          `path:"/exam/attempts/{id}/start" method:"post" tags:"客户端-考试" summary:"开考"`
	Id              int64 `json:"id" in:"path" v:"required|min:1" dc:"答题会话ID"`
	DurationSeconds int   `json:"duration_seconds" dc:"可选，覆盖默认时长，服务端会按 max 夹紧"`
}

type AttemptStartRes struct{}

type AttemptGetReq struct {
	g.Meta `path:"/exam/attempts/{id}" method:"get" tags:"客户端-考试" summary:"答题会话详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1" dc:"答题会话ID"`
}

type AttemptGetRes struct {
	Id                 int64 `json:"id" dc:"会话ID"`
	ExamPaperId        int64 `json:"exam_paper_id" dc:"exam_paper.id"`
	ExaminationPaperId int64 `json:"examination_paper_id" dc:"兼容字段：mock_examination_paper.id"`
	Status             int   `json:"status" dc:"会话状态"`
	//DurationSeconds    int     `json:"duration_seconds" dc:"考试时长(秒)"`
	StartedAt  string `json:"started_at" dc:"开考时间"`
	DeadlineAt string `json:"deadline_at" dc:"截止时间"`
	//SubmittedAt        string  `json:"submitted_at" dc:"交卷时间"`
	//EndedAt            string  `json:"ended_at" dc:"结束时间"`
	//ObjectiveScore     float64 `json:"objective_score" dc:"客观题得分"`
	//SubjectiveScore    float64 `json:"subjective_score" dc:"主观题得分"`
	//TotalScore         float64 `json:"total_score" dc:"总分"`
	//HasSubjective      int     `json:"has_subjective" dc:"是否含主观题：0否 1是"`
	ServerTime string `json:"server_time" dc:"服务端当前时间"`
	//DeadlineReached  bool   `json:"deadline_reached" dc:"是否已到截止时间"`
	SegmentCode      string `json:"segment_code,omitempty" dc:"当前环节编码"`
	RemainingSeconds *int   `json:"remaining_seconds,omitempty" dc:"当前环节剩余时间(秒)：环节时长减去该环节最近一次提交后流逝的时间"`
}

type AttemptAnswersGetReq struct {
	g.Meta `path:"/exam/attempts/{id}/answers" method:"get" tags:"客户端-考试" summary:"查询答题详情"`
	Id     int64 `json:"id" in:"path" v:"required|min:1" dc:"答题会话ID"`
}

type AttemptAnswersGetRes struct {
	Items []AttemptAnswerItem `json:"items" dc:"答案列表"`
}

type AttemptSaveAnswersReq struct {
	g.Meta      `path:"/exam/attempts/{id}/answers" method:"put" tags:"客户端-考试" summary:"保存答案（批量）"`
	Id          int64               `json:"id" in:"path" v:"required|min:1" dc:"答题会话ID"`
	SegmentCode string              `json:"segment_code" v:"required#err.invalid_params" dc:"当前作答环节编码（如 listen/read/write）"`
	Items       []AttemptAnswerItem `json:"items" dc:"答案列表"`
}

// AttemptAnswerItem 单题答案。客观题传 option_id，主观题/填空传 text，两者互斥；未作答时均为零值。
type AttemptAnswerItem struct {
	QuestionId int64  `json:"question_id" v:"required|min:1" dc:"题目ID"`
	OptionID   int64  `json:"option_id,omitempty" dc:"客观题：所选选项ID（与 text 互斥）"`
	Text       string `json:"text,omitempty" dc:"主观题/填空：文本（与 option_id 互斥）"`
}

type AttemptSaveAnswersRes struct{}

type AttemptSubmitReq struct {
	g.Meta `path:"/exam/attempts/{id}/submit" method:"post" tags:"客户端-考试" summary:"交卷"`
	Id     int64 `json:"id" in:"path" v:"required|min:1" dc:"答题会话ID"`
}

type AttemptSubmitRes struct{}

// AttemptRandomAnswersReq 测试专用：按试卷全部小题随机填答并保存（需配置 exam.enableRandomAnswerHelper）。
type AttemptRandomAnswersReq struct {
	g.Meta    `path:"/exam/papers/{paperId}/attempts/{attemptId}/random-answers" method:"post" tags:"客户端-考试-测试" summary:"随机填答（仅测试环境）"`
	PaperId   int64 `json:"paperId" in:"path" v:"required|min:1" dc:"exam_paper.id"`
	AttemptId int64 `json:"attemptId" in:"path" v:"required|min:1" dc:"答题会话ID"`
}

type AttemptSaveAnswersBody struct {
	Items []AttemptAnswerItem `json:"items" dc:"答案列表"`
}

type AttemptRandomAnswersRes struct {
	GeneratedCount int                    `json:"generated_count" dc:"已生成题目数"`
	SubmitJSON     AttemptSaveAnswersBody `json:"submit_json" dc:"可直接作为 /api/client/exam/attempts/{id}/answers 的请求体"`
}
