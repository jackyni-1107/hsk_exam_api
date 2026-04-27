// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package paper

import (
	"context"
	"exam/internal/model/bo"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"
)

type (
	IPaper interface {
		// IssueAudioHlsPlay 签发短期播放票据，返回相对 play_url（以 / 开头）。
		IssueAudioHlsPlay(ctx context.Context, userID int64, attemptID int64, questionID int64) (playURL string, expiresAt string, err error)
		// IssuePaperHlsPlay 基于试卷级 HLS 配置签发短期播放票据，返回相对 play_url（以 / 开头）。
		// paperID 为 exam_paper.id（与客户端试卷路径参数一致）。
		IssuePaperHlsPlay(ctx context.Context, userID int64, paperID int64) (playURL string, expiresAt string, err error)
		// BuildHlsM3U8Playlist 校验 Redis 票据并生成 m3u8（内嵌 presigned URL）。
		BuildHlsM3U8Playlist(ctx context.Context, ticket string) ([]byte, error)
		InvalidatePaperForExamCache(ctx context.Context, examPaperId int64)
		InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64)
		// PaperDetail 返回试卷及嵌套的大题、题块、题目、选项。
		PaperDetail(ctx context.Context, examPaperId int64) (*exambo.PaperDetailTree, error)
		// PaperDetailSection 只加载单个大题的完整详情。
		PaperDetailSection(ctx context.Context, examPaperId int64, sectionId int64) (*exambo.SectionDetailView, error)
		// UpdatePaperSettings 修改试卷的 HLS 配置。
		UpdatePaperSettings(ctx context.Context, examPaperId int64, in exambo.PaperHlsExamAdminUpdate, updater string) error
		// UpdatePaperMeta 修改试卷元数据，不包含 HLS 与题目树。
		UpdatePaperMeta(ctx context.Context, examPaperId int64, in exambo.PaperMetaAdminUpdate, updater string) error
		// RandomFillAnswersForTest 仅返回随机答案草稿列表，不写库。若需生成并保存，使用 RandomFillAndSaveAnswers。
		RandomFillAnswersForTest(ctx context.Context, userID int64, examPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error)
		PaperSectionTopicForExam(ctx context.Context, examPaperID int64, sectionId int64) (*exambo.SectionTopic, error)
		PaperDetailForExamInit(ctx context.Context, examPaperID int64) (*exambo.PaperDetailForExamInitTree, error)
		PaperBootstrapForExam(ctx context.Context, examPaperID int64) (*exambo.PaperDetailForExamInitTree, []exambo.PaperPrepareSegment, *mockentity.MockExaminationPaper, error)
		PaperPrepareSegments(ctx context.Context, examPaperID int64) ([]exambo.PaperPrepareSegment, error)
		PaperSectionDetailForExam(ctx context.Context, examPaperID int64, sectionId int64) (*exambo.SectionDetailForExamView, error)
		// ImportFromIndex 根据 mock_examination_paper.resource_url 推导 index.json 并导入 exam_* 表。
		ImportFromIndex(ctx context.Context, p exambo.ImportParams) (*exambo.ImportResult, error)
		// PaperList 分页查询试卷列表。
		PaperList(ctx context.Context, page int, size int, level string, mockLevelId int64) (list []examentity.ExamPaper, total int, err error)
		// PaperPurgePhysical 从数据库永久删除 exam_paper 及其题目树（option/question/block/section）。
		// 调用方须已通过身份校验；删除前需校验 confirmText == "DELETE:{exam_paper_id}"。
		// 若试卷仍被考试批次引用或存在未逻辑删除的答题会话，则拒绝删除。
		PaperPurgePhysical(ctx context.Context, examPaperId int64, confirmText string) error
	}
)

var (
	localPaper IPaper
)

func Paper() IPaper {
	if localPaper == nil {
		panic("implement not found for interface IPaper, forgot register?")
	}
	return localPaper
}

func RegisterPaper(i IPaper) {
	localPaper = i
}
