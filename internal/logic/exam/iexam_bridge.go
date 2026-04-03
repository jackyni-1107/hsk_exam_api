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
	"exam/internal/logic/clientexam"
	"exam/internal/model/bo"
	"exam/internal/model/entity"
)

func (s *sExam) PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (map[string]interface{}, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	var sec entity.ExamSection
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
	return clientexam.CreateAttempt(ctx, userID, mockPaperID)
}

func (s *sExam) StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error {
	return clientexam.StartAttempt(ctx, userID, attemptID, clientDurationSeconds)
}

func (s *sExam) GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error) {
	v, err := clientexam.GetAttempt(ctx, userID, attemptID)
	if err != nil {
		return nil, err
	}
	return &bo.AttemptView{
		Attempt:         v.Attempt,
		ServerTime:      v.ServerTime,
		DeadlineReached: v.DeadlineReached,
	}, nil
}

func (s *sExam) SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error {
	x := make([]clientexam.SaveAnswerItem, len(items))
	for i := range items {
		x[i] = clientexam.SaveAnswerItem{
			QuestionID:      items[i].QuestionID,
			AnswerJSON:      items[i].AnswerJSON,
			ExpectedVersion: items[i].ExpectedVersion,
		}
	}
	return clientexam.SaveAnswers(ctx, userID, attemptID, x)
}

func (s *sExam) SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error {
	return clientexam.SubmitAttempt(ctx, userID, attemptID)
}

func (s *sExam) MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error {
	return clientexam.MarkSubmittedIfOverdue(ctx, attemptID)
}

func (s *sExam) FinalizeAttempt(ctx context.Context, attemptID int64) error {
	return clientexam.FinalizeAttempt(ctx, attemptID)
}

func (s *sExam) TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error) {
	return clientexam.TryAcquireSubmitLock(ctx, attemptID)
}

func (s *sExam) ReleaseSubmitLock(ctx context.Context, attemptID int64) {
	clientexam.ReleaseSubmitLock(ctx, attemptID)
}

func (s *sExam) RateLimitSaveAnswers(ctx context.Context, attemptID int64, perSecond int) error {
	return clientexam.RateLimitSaveAnswers(ctx, attemptID, perSecond)
}

func (s *sExam) ParseAnswerPayload(str string) bo.AnswerPayload {
	p := clientexam.ParseAnswerPayload(str)
	return bo.AnswerPayload{SelectedOptionIDs: p.SelectedOptionIDs, Text: p.Text}
}

func (s *sExam) PaperHasSubjectiveNonExample(questions []bo.QuestionScoreMeta) bool {
	q := make([]clientexam.QuestionScoreMeta, len(questions))
	for i := range questions {
		q[i] = clientexam.QuestionScoreMeta{
			QuestionID: questions[i].QuestionID, IsExample: questions[i].IsExample, IsSubjective: questions[i].IsSubjective,
			Score: questions[i].Score, CorrectOptIDs: questions[i].CorrectOptIDs,
		}
	}
	return clientexam.PaperHasSubjectiveNonExample(q)
}

func (s *sExam) ScoreObjective(questions []bo.QuestionScoreMeta, answers map[int64]bo.AnswerPayload) (float64, bool) {
	q := make([]clientexam.QuestionScoreMeta, len(questions))
	for i := range questions {
		q[i] = clientexam.QuestionScoreMeta{
			QuestionID: questions[i].QuestionID, IsExample: questions[i].IsExample, IsSubjective: questions[i].IsSubjective,
			Score: questions[i].Score, CorrectOptIDs: questions[i].CorrectOptIDs,
		}
	}
	m := make(map[int64]clientexam.AnswerPayload, len(answers))
	for k, v := range answers {
		m[k] = clientexam.AnswerPayload{SelectedOptionIDs: v.SelectedOptionIDs, Text: v.Text}
	}
	return clientexam.ScoreObjective(q, m)
}

func (s *sExam) ObjectiveAnswerCorrect(correctIDs []int64, selected []int64) bool {
	return clientexam.ObjectiveAnswerCorrect(correctIDs, selected)
}

func (s *sExam) EmptyAnswerRowsForPaper(questionIDs []int64) []int64 {
	return clientexam.EmptyAnswerRowsForPaper(questionIDs)
}

func (s *sExam) InvalidatePaperForExamCache(ctx context.Context, examPaperId int64) {
	InvalidatePaperForExamCache(ctx, examPaperId)
}

func (s *sExam) InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId int64, sectionId int64) {
	InvalidatePaperSectionForExamCache(ctx, examPaperId, sectionId)
}

func (s *sExam) PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*bo.PaperDetailForExamInitTree, error) {
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
	var out bo.PaperDetailForExamInitTree
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *sExam) PaperSectionDetailForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*bo.SectionDetailForExamView, error) {
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
	var out bo.SectionDetailForExamView
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type randomAnswerPayload struct {
	SelectedOptionIDs []int64 `json:"selected_option_ids,omitempty"`
	Text              string  `json:"text,omitempty"`
}

func (s *sExam) RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error) {
	cfg := clientexam.LoadExamCfg(ctx)
	if !cfg.EnableRandomAnswerHelper {
		return nil, gerror.NewCode(consts.CodeExamTestHelperDisabled, "")
	}
	var att entity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("client_user_id", userID).
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

	var qs []entity.ExamQuestion
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
		var opts []entity.ExamOption
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
		picked := randomPickOptionIDs(ids)
		b, err := json.Marshal(randomAnswerPayload{SelectedOptionIDs: picked})
		if err != nil {
			return nil, err
		}
		out = append(out, bo.RandomAnswerDraftItem{QuestionID: q.Id, Answer: string(b)})
	}
	return out, nil
}

func randomPickOptionIDs(optionIDs []int64) []int64 {
	if len(optionIDs) == 0 {
		return nil
	}
	x := append([]int64(nil), optionIDs...)
	rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
	n := 1 + rand.IntN(len(x))
	return x[:n]
}
