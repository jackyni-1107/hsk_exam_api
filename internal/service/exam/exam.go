// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package exam

import (
	"context"
	"exam/internal/model/bo"
	"exam/internal/model/entity"
)

type (
	IExam interface {
		// AttemptAdminList 分页查询答题会话（联表学员、试卷）；examBatchId>0 时按 r.exam_batch_id 筛选。
		AttemptAdminList(ctx context.Context, page int, size int, level string, examinationPaperId int64, examBatchId int64, status int, username string) ([]bo.AttemptAdminListRow, int, error)
		// AttemptAdminDetail 按 id 加载会话、学员、试卷及答题明细（含客观题是否选对）。
		AttemptAdminDetail(ctx context.Context, attemptID int64) (*bo.AttemptAdminDetailView, error)
		// AttemptAdminSaveSubjectiveScores 写入主观题人工分并汇总 subjective_score、total_score（允许部分题目已评）。
		AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []bo.SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error)
		// CreateAttempt 已废弃，返回 err.exam_attempt_use_batch_api。
		CreateAttempt(ctx context.Context, userID int64, mockPaperID int64) (int64, error)
		// CreateAttemptForBatch 按批次与等级创建会话（未开始）。
		CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64, mockLevelID int64) (int64, error)
		// StartAttempt 开考：进入进行中并写入截止时间。
		StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error
		// GetAttempt 查询会话；若已超时仍进行中则自动标记为已交卷（算分由定时任务完成）。
		GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error)
		// SaveAnswers 批量保存答案（限流在调用方或此处）。
		SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error
		// SubmitAttempt 主动交卷：仅标记为已交卷，客观题算分由定时任务统一处理。
		SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error
		// MarkSubmittedIfOverdue 供定时任务：超时未操作会话标记为已交卷（不校验用户）。
		MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error
		// FinalizeAttempt 对已交卷（待算分）会话计算客观分并置为已结束，写入 exam_result。供定时任务调用。
		FinalizeAttempt(ctx context.Context, attemptID int64) error
		// RandomFillAnswersForTest 按试卷小题随机生成提交草稿（不入库）。仅当配置 exam.enableRandomAnswerHelper=true 时可用。
		// paperID 须与 attempt 所属试卷一致；跳过例题；客观题仅单选 1 个 option id；主观题写入随机占位文本。
		RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) (items []bo.RandomAnswerDraftItem, err error)
		// TryAcquireSubmitLock 交卷/超时自动交卷互斥，避免重复计分。
		TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error)
		// ReleaseSubmitLock 释放交卷锁。
		ReleaseSubmitLock(ctx context.Context, attemptID int64)
		// RateLimitSaveAnswers 单会话保存答案频率限制，超限返回 CodeTooManyRequests。
		RateLimitSaveAnswers(ctx context.Context, attemptID int64, perSecond int) error
		// ParseAnswerPayload 解析答题 JSON。
		ParseAnswerPayload(str string) bo.AnswerPayload
		// PaperHasSubjectiveNonExample 试卷是否含需人工的主观题（非例题）。
		PaperHasSubjectiveNonExample(questions []bo.QuestionScoreMeta) bool
		// ScoreObjective 仅客观题自动分；例题与主观题不计分。返回答案侧客观题得分与试卷是否含主观题。
		ScoreObjective(questions []bo.QuestionScoreMeta, answers map[int64]bo.AnswerPayload) (objective float64, paperHasSubjective bool)
		// ObjectiveAnswerCorrect 客观题是否选对（多选需与正确选项 id 集合完全一致，顺序无关）。
		ObjectiveAnswerCorrect(correctIDs []int64, selected []int64) bool
		// EmptyAnswerRowsForPaper 根据试卷全部小题 ID 生成「空答题行」描述（供单测与客户端初始化占位）。
		EmptyAnswerRowsForPaper(questionIDs []int64) []int64
		// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
		PaperDetail(ctx context.Context, examPaperId int64) (*bo.PaperDetailTree, error)
		// InvalidatePaperForExamCache 试卷树变更后删除考前相关缓存（初始化 + 各 section 详情 + 历史整卷 key）。
		InvalidatePaperForExamCache(ctx context.Context, examPaperId int64)
		// InvalidatePaperSectionForExamCache 删除单个 section 的考前详情缓存（精确 key）。
		InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64)
		// PaperDetailForExamInit 客户端考前初始化：仅试卷结构（paper + section 概要 + block 概要 + 题量），不含题目与选项。
		// 流程：Redis → singleflight → DB；TTL 1h。
		PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*bo.PaperDetailForExamInitTree, error)
		// PaperSectionDetailForExam 按 section 拉取完整题目树（blocks + questions + options），不含选项正误。
		PaperSectionDetailForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*bo.SectionDetailForExamView, error)
		// PaperSectionTopicForExam 按 section 返回与 topic JSON 根结构一致的 map（脱敏 + 注入 question_id）。
		PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (map[string]interface{}, error)
		// ImportFromIndex 拉取或解析 index.json，写入 exam_* 表。
		ImportFromIndex(ctx context.Context, p bo.ImportParams) (*bo.ImportResult, error)
		// UpdatePaperSettings 修改试卷听力 HLS 配置（答题时长以 mock 卷为准）。
		UpdatePaperSettings(ctx context.Context, examPaperId int64, in bo.PaperHlsExamAdminUpdate, updater string) error
		// PaperList 分页试卷列表（管理端）
		PaperList(ctx context.Context, page int, size int, level string) (list []entity.ExamPaper, total int, err error)
		// ExamBatch* 考试批次（时间窗、多选 mock_levels、批次学员）
		ExamBatchList(ctx context.Context, mockPaperID int64, page, size int) (list []bo.ExamBatchAdminItem, total int, err error)
		ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error)
		ExamBatchCreate(ctx context.Context, mockPaperID int64, title, examStartAt, examEndAt string, mockLevelIds []int64, creator string) (id int64, err error)
		ExamBatchUpdate(ctx context.Context, id int64, title, examStartAt, examEndAt string, mockLevelIds []int64, updater string) error
		ExamBatchDelete(ctx context.Context, id int64) error
		ExamBatchMembersImport(ctx context.Context, batchID int64, mockLevelID int64, memberIDs []int64, creator string) (inserted int, err error)
		ExamBatchMemberList(ctx context.Context, batchID int64, page, size int) (list []bo.ExamBatchMemberAdminRow, total int, err error)
		ExamBatchMembersRemove(ctx context.Context, batchID int64, mockLevelID int64, memberIDs []int64) (removed int, err error)
		// MyExamBatches 学员端：当前会员所在批次 + 试卷信息，分页。
		MyExamBatches(ctx context.Context, memberID int64, page, size int) (list []bo.MyExamBatchItem, total int, err error)
		// IssueAudioHlsPlay 校验会话与题目 HLS 配置后写入 Redis 票据，返回相对 play_url。
		IssueAudioHlsPlay(ctx context.Context, userID, attemptID, questionID int64) (playURL string, expiresAt string, err error)
		// BuildHlsM3U8Playlist 根据票据生成内嵌 presigned URL 的 m3u8 正文。
		BuildHlsM3U8Playlist(ctx context.Context, ticket string) ([]byte, error)
	}
)

var (
	localExam IExam
)

func Exam() IExam {
	if localExam == nil {
		panic("implement not found for interface IExam, forgot register?")
	}
	return localExam
}

func RegisterExam(i IExam) {
	localExam = i
}
