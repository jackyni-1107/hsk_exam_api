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
		IssuePaperHlsPlay(ctx context.Context, userID int64, paperID int64) (playURL string, expiresAt string, err error)
		// BuildHlsM3U8Playlist 校验 Redis 票据并生成 m3u8（内嵌 presigned URL）。
		BuildHlsM3U8Playlist(ctx context.Context, ticket string) ([]byte, error)
		// ImportFromIndex 根据 mock_examination_paper.resource_url 推导 index.json 地址并导入到 exam_* 表。
		ImportFromIndex(ctx context.Context, p exambo.ImportParams) (*exambo.ImportResult, error)
		// PaperList 分页试卷列表（管理端）
		PaperList(ctx context.Context, page int, size int, level string) (list []examentity.ExamPaper, total int, err error)
		InvalidatePaperForExamCache(ctx context.Context, examPaperId int64)
		InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64)
		// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
		PaperDetail(ctx context.Context, examPaperId int64) (*exambo.PaperDetailTree, error)
		// UpdatePaperSettings 修改试卷听力 HLS 配置（答题时长以 mock_examination_paper 为准）。
		UpdatePaperSettings(ctx context.Context, examPaperId int64, in exambo.PaperHlsExamAdminUpdate, updater string) error
		// UpdatePaperMeta 管理端修改试卷元数据（不含 HLS、题目树）。
		UpdatePaperMeta(ctx context.Context, examPaperId int64, in exambo.PaperMetaAdminUpdate, updater string) error
		// RandomFillAnswersForTest 仅返回随机答案草稿列表，不写库。若需生成并保存，使用 RandomFillAndSaveAnswers。
		RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error)
		PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*exambo.SectionTopic, error)
		PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, error)
		PaperBootstrapForExam(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, []exambo.PaperPrepareSegment, *mockentity.MockExaminationPaper, error)
		PaperPrepareSegments(ctx context.Context, mockPaperID int64) ([]exambo.PaperPrepareSegment, error)
		PaperSectionDetailForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*exambo.SectionDetailForExamView, error)
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
