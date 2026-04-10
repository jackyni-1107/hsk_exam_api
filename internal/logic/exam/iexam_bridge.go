package exam

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/utility/exampaper"
	"exam/internal/utility/examutil"
)

func (s *sExam) PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (map[string]interface{}, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	var sec examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("id", sectionId).
		Where("exam_paper_id", paper.Id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&sec); err != nil {
		return nil, err
	}
	if sec.Id == 0 {
		return nil, gerror.NewCode(consts.CodeExamSectionNotFound)
	}
	if sec.TopicJson == "" {
		return nil, gerror.NewCode(consts.CodeExamSectionTopicEmpty)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(sec.TopicJson), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *sExam) CreateAttempt(ctx context.Context, userID int64, mockPaperID int64) (int64, error) {
	return CreateAttempt(ctx, userID, mockPaperID)
}

func (s *sExam) CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64) (int64, error) {
	return CreateAttemptForBatch(ctx, userID, batchID)
}

func (s *sExam) StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error {
	return StartAttempt(ctx, userID, attemptID, clientDurationSeconds)
}

func (s *sExam) GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error) {
	return GetAttempt(ctx, userID, attemptID)
}

func (s *sExam) SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error {
	return SaveAnswers(ctx, userID, attemptID, items)
}

func (s *sExam) SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error {
	return SubmitAttempt(ctx, userID, attemptID)
}

func (s *sExam) MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error {
	return MarkSubmittedIfOverdue(ctx, attemptID)
}

// FinalizeAttempt 仅供 sys_task（ExamScoreFinalizeHandler）算分；客户端交卷路径不应调用。
func (s *sExam) FinalizeAttempt(ctx context.Context, attemptID int64) error {
	return FinalizeAttempt(ctx, attemptID)
}

func (s *sExam) TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error) {
	return TryAcquireSubmitLock(ctx, attemptID)
}

func (s *sExam) ReleaseSubmitLock(ctx context.Context, attemptID int64) {
	ReleaseSubmitLock(ctx, attemptID)
}

func (s *sExam) RateLimitSaveAnswers(ctx context.Context, attemptID int64, perSecond int) error {
	return RateLimitSaveAnswers(ctx, attemptID, perSecond)
}

func (s *sExam) ParseAnswerPayload(str string) bo.AnswerPayload {
	return examutil.ParseAnswerPayload(str)
}

func (s *sExam) PaperHasSubjectiveNonExample(questions []bo.QuestionScoreMeta) bool {
	return examutil.PaperHasSubjectiveNonExample(questions)
}

func (s *sExam) ScoreObjective(questions []bo.QuestionScoreMeta, answers map[int64]bo.AnswerPayload) (float64, bool) {
	return examutil.ScoreObjective(questions, answers)
}

func (s *sExam) ObjectiveAnswerCorrect(correctIDs []int64, selected []int64) bool {
	return examutil.ObjectiveAnswerCorrect(correctIDs, selected)
}

func (s *sExam) EmptyAnswerRowsForPaper(questionIDs []int64) []int64 {
	return examutil.EmptyAnswerRowsForPaper(questionIDs)
}

func (s *sExam) InvalidatePaperForExamCache(ctx context.Context, examPaperId int64) {
	InvalidatePaperForExamCache(ctx, examPaperId)
}

func (s *sExam) InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64) {
	InvalidatePaperSectionForExamCache(ctx, examPaperId, sectionId)
}

func (s *sExam) PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	t, err := PaperDetailForExamInit(ctx, paper.Id)
	if err != nil {
		return nil, err
	}
	out := paperDetailForExamInitTreeToBO(t)
	return &out, nil
}

func (s *sExam) PaperSectionDetailForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*exambo.SectionDetailForExamView, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	t, err := PaperSectionDetailForExam(ctx, paper.Id, sectionId)
	if err != nil {
		return nil, err
	}
	out := sectionDetailForExamViewToBO(t)
	return &out, nil
}

func (s *sExam) RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error) {
	cfg := LoadExamCfg(ctx)
	if !cfg.EnableRandomAnswerHelper {
		return nil, gerror.NewCode(consts.CodeExamTestHelperDisabled)
	}
	var att examentity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att); err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.NewCode(consts.CodeExamAttemptNotFound)
	}
	if att.MockExaminationPaperId != mockPaperID {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	if att.Status != consts.ExamAttemptInProgress {
		return nil, gerror.NewCode(consts.CodeExamAlreadySubmitted)
	}
	now := gtime.Now()
	if att.DeadlineAt != nil && att.DeadlineAt.Before(now) {
		return nil, gerror.NewCode(consts.CodeExamTimeExpired)
	}

	var qs []examentity.ExamQuestion
	if err := dao.ExamQuestion.Ctx(ctx).
		Where("exam_paper_id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&qs); err != nil {
		return nil, err
	}
	qIDs := make([]int64, 0, len(qs))
	for _, q := range qs {
		if q.IsExample != 0 {
			continue
		}
		qIDs = append(qIDs, q.Id)
	}
	optsByQ, err := loadExamOptionsGrouped(ctx, qIDs)
	if err != nil {
		return nil, err
	}
	out := make([]bo.RandomAnswerDraftItem, 0, len(qs))
	for _, q := range qs {
		if q.IsExample != 0 {
			continue
		}
		opts := optsByQ[q.Id]
		if q.IsSubjective != 0 {
			out = append(out, bo.RandomAnswerDraftItem{
				QuestionID: q.Id,
				Answer:     fmt.Sprintf("test-rand-%d-%016x", gtime.Now().TimestampMilli(), rand.Uint64()),
			})
			continue
		}
		if len(opts) == 0 {
			continue
		}
		ids := make([]int64, len(opts))
		for i, o := range opts {
			ids[i] = o.Id
		}
		picked := ids[rand.IntN(len(ids))]
		out = append(out, bo.RandomAnswerDraftItem{QuestionID: q.Id, Answer: picked})
	}
	return out, nil
}

func paperDetailForExamInitTreeToBO(t *PaperDetailForExamInitTree) exambo.PaperDetailForExamInitTree {
	if t == nil {
		return exambo.PaperDetailForExamInitTree{}
	}
	out := exambo.PaperDetailForExamInitTree{
		Paper: exambo.PaperHeadForExamView{
			Id:                 t.Paper.Id,
			Level:              t.Paper.Level,
			PaperId:            t.Paper.PaperId,
			Title:              t.Paper.Title,
			PrepareInstruction: t.Paper.PrepareInstruction,
			PrepareAudioFile:   t.Paper.PrepareAudioFile,
			SourceBaseUrl:      t.Paper.SourceBaseUrl,
			IndexJson:          t.Paper.IndexJson,
			DurationSeconds:    t.Paper.DurationSeconds,
			CreateTime:         t.Paper.CreateTime,
		},
		Sections: make([]exambo.SectionOutlineForExamView, len(t.Sections)),
	}
	for i, s := range t.Sections {
		blocks := make([]exambo.BlockOutlineForExamView, len(s.Blocks))
		for j, b := range s.Blocks {
			blocks[j] = exambo.BlockOutlineForExamView{
				Id:                      b.Id,
				BlockOrder:              b.BlockOrder,
				GroupIndex:              b.GroupIndex,
				QuestionDescriptionJson: b.QuestionDescriptionJson,
				QuestionCount:           b.QuestionCount,
			}
		}
		out.Sections[i] = exambo.SectionOutlineForExamView{
			Id:             s.Id,
			SortOrder:      s.SortOrder,
			TopicTitle:     s.TopicTitle,
			TopicSubtitle:  s.TopicSubtitle,
			TopicType:      s.TopicType,
			PartCode:       s.PartCode,
			SegmentCode:    s.SegmentCode,
			TopicItemsFile: s.TopicItemsFile,
			TopicJson:      s.TopicJson,
			Blocks:         blocks,
		}
	}
	return out
}

func sectionDetailForExamViewToBO(v *SectionDetailForExamView) exambo.SectionDetailForExamView {
	if v == nil {
		return exambo.SectionDetailForExamView{}
	}
	out := exambo.SectionDetailForExamView{
		Id:             v.Id,
		SortOrder:      v.SortOrder,
		TopicTitle:     v.TopicTitle,
		TopicSubtitle:  v.TopicSubtitle,
		TopicType:      v.TopicType,
		PartCode:       v.PartCode,
		SegmentCode:    v.SegmentCode,
		TopicItemsFile: v.TopicItemsFile,
		TopicJson:      v.TopicJson,
		Blocks:         make([]exambo.BlockDetailForExamView, len(v.Blocks)),
	}
	for i, b := range v.Blocks {
		out.Blocks[i] = blockDetailForExamViewToBO(b)
	}
	return out
}

func blockDetailForExamViewToBO(b BlockDetailForExamView) exambo.BlockDetailForExamView {
	qs := make([]exambo.QuestionDetailForExamView, len(b.Questions))
	for i, q := range b.Questions {
		qs[i] = questionDetailForExamViewToBO(q)
	}
	return exambo.BlockDetailForExamView{
		Id:                      b.Id,
		BlockOrder:              b.BlockOrder,
		GroupIndex:              b.GroupIndex,
		QuestionDescriptionJson: b.QuestionDescriptionJson,
		Questions:               qs,
	}
}

func questionDetailForExamViewToBO(q QuestionDetailForExamView) exambo.QuestionDetailForExamView {
	opts := make([]exambo.OptionDetailForExamView, len(q.Options))
	for i, o := range q.Options {
		opts[i] = exambo.OptionDetailForExamView{
			Id:         o.Id,
			Flag:       o.Flag,
			SortOrder:  o.SortOrder,
			OptionType: o.OptionType,
			Content:    o.Content,
		}
	}
	return exambo.QuestionDetailForExamView{
		Id:                      q.Id,
		SortInBlock:             q.SortInBlock,
		QuestionNo:              q.QuestionNo,
		Score:                   q.Score,
		IsExample:               q.IsExample,
		IsSubjective:            q.IsSubjective,
		ContentType:             q.ContentType,
		AudioFile:               q.AudioFile,
		StemText:                q.StemText,
		ScreenTextJson:          q.ScreenTextJson,
		AnalysisJson:            q.AnalysisJson,
		QuestionDescriptionJson: q.QuestionDescriptionJson,
		RawJson:                 q.RawJson,
		Options:                 opts,
	}
}
