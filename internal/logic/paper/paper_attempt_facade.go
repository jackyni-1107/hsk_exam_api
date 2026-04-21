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
)

// RandomFillAnswersForTest 仅返回随机答案草稿列表，不写库。若需生成并保存，使用 RandomFillAndSaveAnswers。
func (s *sPaper) RandomFillAnswersForTest(ctx context.Context, userID int64, examPaperID int64, attemptID int64) ([]bo.RandomAnswerDraftItem, error) {
	cfg := LoadExamCfg(ctx)
	if !cfg.EnableRandomAnswerHelper {
		return nil, gerror.NewCode(consts.CodeExamTestHelperDisabled)
	}
	att, err := attempt.LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	if att.ExamPaperId != examPaperID {
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
				Text:       fmt.Sprintf("test-rand-%d-%016x", gtime.Now().TimestampMilli(), rand.Uint64()),
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
		out = append(out, bo.RandomAnswerDraftItem{QuestionID: q.Id, OptionID: picked})
	}
	return out, nil
}
