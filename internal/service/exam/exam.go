// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package exam

import (
	"context"
	"exam/internal/model/bo"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
)

type (
	IExam interface {
		// AttemptAdminList 分页查询答题会话（联表学员、试卷）。
		AttemptAdminList(ctx context.Context, page int, size int, level string, examinationPaperId int64, examBatchId int64, status int, username string) ([]bo.AttemptAdminListRow, int, error)
		// AttemptAdminDetail 按 id 加载会话、学员、试卷及答题明细（含客观题是否选对）。
		AttemptAdminDetail(ctx context.Context, attemptID int64) (*bo.AttemptAdminDetailView, error)
		// AttemptAdminSaveSubjectiveScores 写入主观题人工分并汇总 subjective_score、total_score（允许部分题目已评）。
		AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []bo.SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error)
		// ExamBatchList 管理端批次分页；mockPaperID=0 时不按卷筛选。
		ExamBatchList(ctx context.Context, mockPaperID int64, page int, size int) (list []bo.ExamBatchAdminItem, total int, err error)
		// ExamBatchDetail 批次详情（含 Mock 卷 id 列表与学员数）。
		ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error)
		// ExamBatchCreate 创建批次并写入多选 Mock 卷（各卷须已导入 exam_paper）。
		ExamBatchCreate(ctx context.Context, title string, examStartAt string, examEndAt string, mockExaminationPaperIds []int64, creator string) (int64, error)
		// ExamBatchUpdate 更新批次时间与 Mock 卷多选（全量替换卷关联）。
		ExamBatchUpdate(ctx context.Context, id int64, title string, examStartAt string, examEndAt string, mockExaminationPaperIds []int64, updater string) error
		// ExamBatchDelete 逻辑删除批次（不删学员关联与等级行，便于审计；列表已过滤）。
		ExamBatchDelete(ctx context.Context, id int64) error
		// ExamBatchMembersImport 导入学员（指定批次内 Mock 卷）；已存在 (batch,member,paper) 则跳过。
		ExamBatchMembersImport(ctx context.Context, batchID int64, mockExaminationPaperId int64, memberIDs []int64, creator string) (inserted int, err error)
		// ExamBatchMemberList 批次内学员分页。
		ExamBatchMemberList(ctx context.Context, batchID int64, page int, size int) (list []bo.ExamBatchMemberAdminRow, total int, err error)
		// ExamBatchMembersRemove 从批次移除学员（指定 mock_examination_paper_id）。
		ExamBatchMembersRemove(ctx context.Context, batchID int64, mockExaminationPaperId int64, memberIDs []int64) (removed int, err error)
		// ExamBatchMemberDetail 获取指定 Mock 卷下的批次成员绑定行。
		ExamBatchMemberDetail(ctx context.Context, batchID int64, userID int64, mockExaminationPaperId int64) (*examentity.ExamBatchMember, error)
		// MyExamBatches 当前学员在批次成员表中的考试批次（联 exam_batch、exam_paper）。
		// 仅返回「考试时间窗口内」且「未提交/未结束」的考试资格列表，不分页。
		MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error)
		// IssueAudioHlsPlay 签发短期播放票据，返回相对 play_url（以 / 开头）。
		IssueAudioHlsPlay(ctx context.Context, userID int64, attemptID int64, questionID int64) (playURL string, expiresAt string, err error)
		// IssuePaperHlsPlay 基于试卷级 HLS 配置签发短期播放票据，返回相对 play_url（以 / 开头）。
		IssuePaperHlsPlay(ctx context.Context, userID int64, paperID int64) (playURL string, expiresAt string, err error)
		// BuildHlsM3U8Playlist 校验 Redis 票据并生成 m3u8（内嵌 presigned URL）。
		BuildHlsM3U8Playlist(ctx context.Context, ticket string) ([]byte, error)
		PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (map[string]interface{}, error)
		CreateAttempt(ctx context.Context, userID int64, mockPaperID int64) (int64, error)
		CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64, mockExaminationPaperId int64) (int64, error)
		StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error
		GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error)
		SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error
		SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error
		MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error
		FinalizeAttempt(ctx context.Context, attemptID int64) error
		TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error)
		ReleaseSubmitLock(ctx context.Context, attemptID int64)
		RateLimitSaveAnswers(ctx context.Context, attemptID int64, perSecond int) error
		ParseAnswerPayload(str string) bo.AnswerPayload
		PaperHasSubjectiveNonExample(questions []bo.QuestionScoreMeta) bool
		ScoreObjective(questions []bo.QuestionScoreMeta, answers map[int64]bo.AnswerPayload) (float64, bool)
		ObjectiveAnswerCorrect(correctIDs []int64, selected []int64) bool
		EmptyAnswerRowsForPaper(questionIDs []int64) []int64
		InvalidatePaperForExamCache(ctx context.Context, examPaperId int64)
		InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64)
		PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, error)
		PaperSectionDetailForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*exambo.SectionDetailForExamView, error)
		RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error)
		// ImportFromIndex 拉取或解析 index.json，写入 exam_* 表。
		ImportFromIndex(ctx context.Context, p exambo.ImportParams) (*exambo.ImportResult, error)
		// PaperList 分页试卷列表（管理端）
		PaperList(ctx context.Context, page int, size int, level string) (list []examentity.ExamPaper, total int, err error)
		// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
		PaperDetail(ctx context.Context, examPaperId int64) (*exambo.PaperDetailTree, error)
		// UpdatePaperSettings 修改试卷听力 HLS 配置（答题时长以 mock_examination_paper 为准）。
		UpdatePaperSettings(ctx context.Context, examPaperId int64, in exambo.PaperHlsExamAdminUpdate, updater string) error
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
