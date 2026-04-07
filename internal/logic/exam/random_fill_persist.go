package exam

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sort"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
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

	var att examentity.ExamAttempt
	err = dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
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

	var qs []examentity.ExamQuestion
	if err := dao.ExamQuestion.Ctx(ctx).
		Where("exam_paper_id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&qs); err != nil {
		return 0, err
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
		return 0, err
	}

	items := make([]bo.SaveAnswerItem, 0, len(qs))
	for _, q := range qs {
		if q.IsExample != 0 {
			continue
		}
		opts := optsByQ[q.Id]

		if q.IsSubjective != 0 {
			b, err := json.Marshal(randomAnswerPayload{Text: randomSubjectiveText()})
			if err != nil {
				return 0, err
			}
			items = append(items, bo.SaveAnswerItem{QuestionID: q.Id, AnswerJSON: string(b)})
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
		items = append(items, bo.SaveAnswerItem{QuestionID: q.Id, AnswerJSON: string(b)})
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

// loadExamOptionsGrouped 按 question_id 批量加载选项，每组内按 sort_order 排序。
func loadExamOptionsGrouped(ctx context.Context, questionIDs []int64) (map[int64][]examentity.ExamOption, error) {
	out := make(map[int64][]examentity.ExamOption)
	if len(questionIDs) == 0 {
		return out, nil
	}
	ids := make([]interface{}, len(questionIDs))
	for i, id := range questionIDs {
		ids[i] = id
	}
	var all []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		WhereIn("question_id", ids).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&all); err != nil {
		return nil, err
	}
	for _, o := range all {
		qid := o.QuestionId
		out[qid] = append(out[qid], o)
	}
	for qid := range out {
		opts := out[qid]
		sort.Slice(opts, func(i, j int) bool { return opts[i].SortOrder < opts[j].SortOrder })
	}
	return out, nil
}
