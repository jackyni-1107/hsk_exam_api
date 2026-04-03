package clientexam

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/entity"
)

type randomAnswerPayload struct {
	SelectedOptionIDs []int64 `json:"selected_option_ids,omitempty"`
	Text              string  `json:"text,omitempty"`
}

// RandomFillAnswersForTest 按试卷小题随机生成答案并批量保存。仅当配置 exam.enableRandomAnswerHelper=true 时可用。
// paperID 为 mock_examination_paper.id，须与 attempt 所属卷一致；跳过例题；客观题从选项中随机选 1..N 项；主观题写入随机占位文本。
func RandomFillAnswersForTest(ctx context.Context, userID, paperID, attemptID int64) (filled int, err error) {
	cfg := LoadExamCfg(ctx)
	if !cfg.EnableRandomAnswerHelper {
		return 0, gerror.NewCode(consts.CodeExamTestHelperDisabled, "")
	}
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)

	var att entity.ExamAttempt
	err = dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("client_user_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att)
	if err != nil {
		return 0, err
	}
	if att.Id == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	if att.MockExaminationPaperId != paperID {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	switch att.Status {
	case consts.ExamAttemptNotStarted:
		return 0, gerror.NewCode(consts.CodeExamNotStarted, "")
	case consts.ExamAttemptInProgress:
		// ok
	default:
		return 0, gerror.NewCode(consts.CodeExamAlreadySubmitted, "")
	}
	now := gtime.Now()
	if att.DeadlineAt != nil && att.DeadlineAt.Before(now) {
		return 0, gerror.NewCode(consts.CodeExamTimeExpired, "")
	}

	var qs []entity.ExamQuestion
	if err := dao.ExamQuestion.Ctx(ctx).
		Where("exam_paper_id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&qs); err != nil {
		return 0, err
	}

	items := make([]SaveAnswerItem, 0, len(qs))
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
			b, err := json.Marshal(randomAnswerPayload{Text: randomSubjectiveText()})
			if err != nil {
				return 0, err
			}
			items = append(items, SaveAnswerItem{QuestionID: q.Id, AnswerJSON: string(b)})
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
			return 0, err
		}
		items = append(items, SaveAnswerItem{QuestionID: q.Id, AnswerJSON: string(b)})
	}

	if len(items) == 0 {
		return 0, nil
	}
	if err := SaveAnswers(ctx, userID, attemptID, items); err != nil {
		return 0, err
	}
	return len(items), nil
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

func randomSubjectiveText() string {
	return fmt.Sprintf("test-rand-%d-%016x", gtime.Now().TimestampMilli(), rand.Uint64())
}
