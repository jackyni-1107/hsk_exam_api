package paper

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/logic/attempt"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
	svcattempt "exam/internal/service/attempt"
	"exam/internal/utility/examutil"
)

func (s *sExam) CreateAttempt(ctx context.Context, userID int64, mockPaperID int64) (int64, error) {
	return attempt.CreateAttempt(ctx, userID, mockPaperID)
}

func (s *sExam) AttemptAdminList(ctx context.Context, page int, size int, level string, examinationPaperId int64, examBatchId int64, status int, username string) ([]bo.AttemptAdminListRow, int, error) {
	return svcattempt.Attempt().AttemptAdminList(ctx, page, size, level, examinationPaperId, examBatchId, status, username)
}

func (s *sExam) AttemptAdminDetail(ctx context.Context, attemptID int64) (*bo.AttemptAdminDetailView, error) {
	return svcattempt.Attempt().AttemptAdminDetail(ctx, attemptID)
}

func (s *sExam) AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []bo.SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error) {
	return svcattempt.Attempt().AttemptAdminSaveSubjectiveScores(ctx, attemptID, items)
}

func (s *sExam) CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64) (int64, error) {
	return svcattempt.Attempt().CreateAttemptForBatch(ctx, userID, batchID)
}

func (s *sExam) StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error {
	return svcattempt.Attempt().StartAttempt(ctx, userID, attemptID, clientDurationSeconds)
}

func (s *sExam) GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error) {
	return svcattempt.Attempt().GetAttempt(ctx, userID, attemptID)
}

func (s *sExam) GetAttemptAnswers(ctx context.Context, userID int64, attemptID int64) ([]bo.AttemptAnswerClientItem, error) {
	return svcattempt.Attempt().GetAttemptAnswers(ctx, userID, attemptID)
}

func (s *sExam) SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error {
	return svcattempt.Attempt().SaveAnswers(ctx, userID, attemptID, items)
}

func (s *sExam) SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error {
	return svcattempt.Attempt().SubmitAttempt(ctx, userID, attemptID)
}

func (s *sExam) MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error {
	return svcattempt.Attempt().MarkSubmittedIfOverdue(ctx, attemptID)
}

// FinalizeAttempt 仅供 sys_task（ExamScoreFinalizeHandler）算分；客户端交卷路径不应调用。
func (s *sExam) FinalizeAttempt(ctx context.Context, attemptID int64) error {
	return svcattempt.Attempt().FinalizeAttempt(ctx, attemptID)
}

func (s *sExam) TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error) {
	return attempt.TryAcquireSubmitLock(ctx, attemptID)
}

func (s *sExam) ReleaseSubmitLock(ctx context.Context, attemptID int64) {
	attempt.ReleaseSubmitLock(ctx, attemptID)
}

func (s *sExam) RateLimitSaveAnswers(ctx context.Context, attemptID int64, perSecond int) error {
	return attempt.RateLimitSaveAnswers(ctx, attemptID, perSecond)
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

// RandomFillAnswersForTest 仅返回随机答案草稿列表，不写库。若需生成并保存，使用 RandomFillAndSaveAnswers。
func (s *sExam) RandomFillAnswersForTest(ctx context.Context, userID int64, mockPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error) {
	cfg := LoadExamCfg(ctx)
	if !cfg.EnableRandomAnswerHelper {
		return nil, gerror.NewCode(consts.CodeExamTestHelperDisabled)
	}
	att, err := attempt.LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	if att.MockExaminationPaperId != mockPaperID {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
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
