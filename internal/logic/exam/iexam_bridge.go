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
	"exam/internal/exampaper"
	"exam/internal/examutil"
	"exam/internal/model/bo"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
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

func (s *sExam) CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64, mockLevelID int64) (int64, error) {
	return CreateAttemptForBatch(ctx, userID, batchID, mockLevelID)
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
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	var out exambo.PaperDetailForExamInitTree
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
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
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	var out exambo.SectionDetailForExamView
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *sExam) RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error) {
	cfg := LoadExamCfg(ctx)
	if !cfg.EnableRandomAnswerHelper {
		return nil, gerror.NewCode(consts.CodeExamTestHelperDisabled, "")
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
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	if att.MockExaminationPaperId != mockPaperID {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	if att.Status != consts.ExamAttemptInProgress {
		return nil, gerror.NewCode(consts.CodeExamAlreadySubmitted, "")
	}
	now := gtime.Now()
	if att.DeadlineAt != nil && att.DeadlineAt.Before(now) {
		return nil, gerror.NewCode(consts.CodeExamTimeExpired, "")
	}

	var qs []examentity.ExamQuestion
	if err := dao.ExamQuestion.Ctx(ctx).
		Where("exam_paper_id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&qs); err != nil {
		return nil, err
	}
	out := make([]bo.RandomAnswerDraftItem, 0, len(qs))
	for _, q := range qs {
		if q.IsExample != 0 {
			continue
		}
		var opts []examentity.ExamOption
		_ = dao.ExamOption.Ctx(ctx).
			Where("question_id", q.Id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("sort_order").
			Scan(&opts)
		if q.IsSubjective != 0 {
			b, err := json.Marshal(randomAnswerPayload{Text: fmt.Sprintf("test-rand-%d-%016x", gtime.Now().TimestampMilli(), rand.Uint64())})
			if err != nil {
				return nil, err
			}
			out = append(out, bo.RandomAnswerDraftItem{QuestionID: q.Id, Answer: string(b)})
			continue
		}
		if len(opts) == 0 {
			continue
		}
		ids := make([]int64, len(opts))
		for i, o := range opts {
			ids[i] = o.Id
		}
		picked := randomPickOptionIDsForFill(ids)
		b, err := json.Marshal(randomAnswerPayload{SelectedOptionIDs: picked})
		if err != nil {
			return nil, err
		}
		out = append(out, bo.RandomAnswerDraftItem{QuestionID: q.Id, Answer: string(b)})
	}
	return out, nil
}

func randomPickOptionIDsForFill(optionIDs []int64) []int64 {
	if len(optionIDs) == 0 {
		return nil
	}
	x := append([]int64(nil), optionIDs...)
	rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
	n := 1 + rand.IntN(len(x))
	return x[:n]
}
