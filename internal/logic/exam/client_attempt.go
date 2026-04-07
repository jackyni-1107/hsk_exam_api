package exam

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/exampaper"
	"exam/internal/examutil"
	"exam/internal/model/bo"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/util"
)

// CreateAttempt 已废弃：请使用 CreateAttemptForBatch（POST /exam/batches/{batchId}/attempts）。
func CreateAttempt(ctx context.Context, userID int64, mockPaperID int64) (int64, error) {
	_ = ctx
	_ = userID
	_ = mockPaperID
	return 0, gerror.NewCode(consts.CodeExamAttemptUseBatchApi, "")
}

func batchExamWindowOpen(now *gtime.Time, start, end *gtime.Time) bool {
	if start == nil || end == nil {
		return false
	}
	if now.Before(start) || now.After(end) {
		return false
	}
	return true
}

// CreateAttemptForBatch 按批次与等级创建会话（未开始）；每用户每批次每等级仅允许一条未删除记录。
func CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64, mockLevelID int64) (int64, error) {
	if batchID <= 0 || mockLevelID <= 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	var batch examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).
		Where("id", batchID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&batch); err != nil {
		return 0, err
	}
	if batch.Id == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchNotFound, "")
	}
	now := gtime.Now()
	if !batchExamWindowOpen(now, batch.ExamStartAt, batch.ExamEndAt) {
		return 0, gerror.NewCode(consts.CodeExamBatchWindowNotOpen, "")
	}
	nLevel, err := dao.ExamBatchMockLevel.Ctx(ctx).
		Where("batch_id", batchID).
		Where("mock_level_id", mockLevelID).
		Count()
	if err != nil {
		return 0, err
	}
	if nLevel == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchLevelNotInBatch, "")
	}
	var link examentity.ExamBatchMember
	if err := dao.ExamBatchMember.Ctx(ctx).
		Where("batch_id", batchID).
		Where("member_id", userID).
		Where("mock_level_id", mockLevelID).
		Scan(&link); err != nil {
		return 0, err
	}
	if link.BatchId == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchMemberNotFound, "")
	}
	paper, err := exampaper.ByMockID(ctx, batch.MockExaminationPaperId)
	if err != nil {
		return 0, err
	}
	dup, err := dao.ExamAttempt.Ctx(ctx).
		Where("member_id", userID).
		Where("exam_batch_id", batchID).
		Where("mock_level_id", mockLevelID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, err
	}
	if dup > 0 {
		return 0, gerror.NewCode(consts.CodeExamAttemptExistsForBatch, "")
	}
	id, err := dao.ExamAttempt.Ctx(ctx).InsertAndGetId(examdo.ExamAttempt{
		MemberId:               userID,
		ExamPaperId:            paper.Id,
		MockExaminationPaperId: paper.MockExaminationPaperId,
		ExamBatchId:            batchID,
		MockLevelId:            mockLevelID,
		Status:                 consts.ExamAttemptNotStarted,
		DurationSeconds:        0,
		Creator:                "client",
		Updater:                "client",
		DeleteFlag:             consts.DeleteFlagNotDeleted,
		CreateTime:             gtime.Now(),
		UpdateTime:             gtime.Now(),
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// StartAttempt 开考：进入进行中并写入截止时间。
func StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error {
	cfg := LoadExamCfg(ctx)
	var att examentity.ExamAttempt
	err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att)
	if err != nil {
		return err
	}
	if att.Id == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	if att.Status != consts.ExamAttemptNotStarted {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	var paper examentity.ExamPaper
	if err := dao.ExamPaper.Ctx(ctx).Where("id", att.ExamPaperId).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&paper); err != nil {
		return err
	}
	if paper.Id == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_paper_not_found")
	}
	dur := ResolveDurationSeconds(cfg, paper.DurationSeconds, clientDurationSeconds)
	now := gtime.Now()
	deadline := gtime.NewFromTimeStamp(now.Timestamp() + int64(dur))

	_, err = dao.ExamAttempt.Ctx(ctx).Where("id", att.Id).Update(examdo.ExamAttempt{
		Status:          consts.ExamAttemptInProgress,
		DurationSeconds: dur,
		StartedAt:       now,
		DeadlineAt:      deadline,
		Updater:         "client",
		UpdateTime:      gtime.Now(),
	})
	return err
}

// GetAttempt 查询会话；若已超时仍进行中则自动交卷并计分。
func GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error) {
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)
	var att examentity.ExamAttempt
	err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att)
	if err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	now := gtime.Now()
	deadlineReached := att.Status == consts.ExamAttemptInProgress && att.DeadlineAt != nil && att.DeadlineAt.Before(now)
	return &bo.AttemptView{
		Attempt:         att,
		ServerTime:      util.ToRFC3339UTCShift(now),
		DeadlineReached: deadlineReached,
	}, nil
}

// SaveAnswers 批量保存答案（限流在调用方或此处）。
func SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error {
	cfg := LoadExamCfg(ctx)
	if err := RateLimitSaveAnswers(ctx, attemptID, cfg.SaveAnswersPerSecond); err != nil {
		return err
	}
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)

	var att examentity.ExamAttempt
	err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att)
	if err != nil {
		return err
	}
	if att.Id == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	if att.Status != consts.ExamAttemptInProgress {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_already_submitted")
	}
	now := gtime.Now()
	if att.DeadlineAt != nil && att.DeadlineAt.Before(now) {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_time_expired")
	}
	if len(items) == 0 {
		return nil
	}
	qids := make([]interface{}, 0, len(items))
	for _, it := range items {
		qids = append(qids, it.QuestionID)
	}
	cnt, err := dao.ExamQuestion.Ctx(ctx).
		Where("exam_paper_id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		WhereIn("id", qids).
		Count()
	if err != nil {
		return err
	}
	if cnt != len(items) {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, it := range items {
			var row examentity.ExamAttemptAnswer
			_ = tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
				Where("attempt_id", attemptID).
				Where("exam_question_id", it.QuestionID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&row)
			if row.Id == 0 {
				_, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Insert(examdo.ExamAttemptAnswer{
					AttemptId:      attemptID,
					ExamQuestionId: it.QuestionID,
					AnswerJson:     it.AnswerJSON,
					Version:        0,
					Creator:        "client",
					Updater:        "client",
					DeleteFlag:     consts.DeleteFlagNotDeleted,
					CreateTime:     gtime.Now(),
					UpdateTime:     gtime.Now(),
				})
				if err != nil {
					return err
				}
				continue
			}
			if it.ExpectedVersion != nil && *it.ExpectedVersion != row.Version {
				return gerror.NewCode(consts.CodeExamAnswerVersionConflict, "")
			}
			_, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Where("id", row.Id).Update(examdo.ExamAttemptAnswer{
				AnswerJson: it.AnswerJSON,
				Version:    row.Version + 1,
				Updater:    "client",
				UpdateTime: gtime.Now(),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// SubmitAttempt 主动交卷：标记已交卷后立即计算客观分并同步 exam_result（与超时自动交卷一致）。
func SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error {
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)
	var att examentity.ExamAttempt
	err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att)
	if err != nil {
		return err
	}
	if att.Id == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	if att.Status == consts.ExamAttemptEnded {
		return nil
	}
	if att.Status == consts.ExamAttemptSubmitted {
		_ = FinalizeAttempt(ctx, attemptID)
		return nil
	}
	if att.Status != consts.ExamAttemptInProgress {
		return gerror.NewCode(consts.CodeInvalidParams, "err.exam_already_submitted")
	}
	if err := markSubmitted(ctx, attemptID, false, "client"); err != nil {
		return err
	}
	_ = FinalizeAttempt(ctx, attemptID)
	return nil
}

// maybeAutoSubmitIfOverdue 考试时间到达且仍为进行中时自动交卷并立即计分（由客户端拉取/保存答案等触发）。
func maybeAutoSubmitIfOverdue(ctx context.Context, userID int64, attemptID int64) error {
	var att examentity.ExamAttempt
	err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att)
	if err != nil || att.Id == 0 {
		return err
	}
	if att.Status != consts.ExamAttemptInProgress || att.DeadlineAt == nil {
		_ = FinalizeAttempt(ctx, attemptID)
		return nil
	}
	now := gtime.Now()
	if !att.DeadlineAt.Before(now) {
		return nil
	}
	if err := markSubmitted(ctx, attemptID, true, "client"); err != nil {
		return err
	}
	_ = FinalizeAttempt(ctx, attemptID)
	return nil
}

// MarkSubmittedIfOverdue 供定时任务：超时未操作会话标记为已交卷并计分（不校验用户）。
func MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error {
	if err := markSubmitted(ctx, attemptID, true, "task"); err != nil {
		return err
	}
	_ = FinalizeAttempt(ctx, attemptID)
	return nil
}

// FinalizeAttempt 对已交卷（待算分）会话计算客观分并置为已结束，写入 exam_result。
func FinalizeAttempt(ctx context.Context, attemptID int64) error {
	return finalizeScoring(ctx, attemptID)
}

func markSubmitted(ctx context.Context, attemptID int64, onlyIfOverdue bool, updater string) error {
	ok, err := TryAcquireSubmitLock(ctx, attemptID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer ReleaseSubmitLock(ctx, attemptID)

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var att examentity.ExamAttempt
		if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&att); err != nil {
			return err
		}
		if att.Id == 0 || att.Status != consts.ExamAttemptInProgress {
			return nil
		}
		now := gtime.Now()
		if onlyIfOverdue {
			if att.DeadlineAt == nil || !att.DeadlineAt.Before(now) {
				return nil
			}
		}
		_, err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Update(examdo.ExamAttempt{
			Status:      consts.ExamAttemptSubmitted,
			SubmittedAt: now,
			Updater:     updater,
			UpdateTime:  gtime.Now(),
		})
		return err
	})
}

func finalizeScoring(ctx context.Context, attemptID int64) error {
	ok, err := TryAcquireSubmitLock(ctx, attemptID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer ReleaseSubmitLock(ctx, attemptID)

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var att examentity.ExamAttempt
		if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&att); err != nil {
			return err
		}
		if att.Id == 0 {
			return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
		}
		if att.Status == consts.ExamAttemptEnded {
			return nil
		}
		if att.Status != consts.ExamAttemptSubmitted {
			return nil
		}

		meta, err := loadQuestionScoreMetaTx(ctx, tx, att.ExamPaperId)
		if err != nil {
			return err
		}
		answers, err := loadAnswersMapTx(ctx, tx, attemptID)
		if err != nil {
			return err
		}
		objScore, hasSubj := examutil.ScoreObjective(meta, answers)
		now := gtime.Now()
		hasFlag := 0
		if hasSubj {
			hasFlag = 1
		}
		_, err = tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Update(examdo.ExamAttempt{
			Status:          consts.ExamAttemptEnded,
			EndedAt:         now,
			ObjectiveScore:  objScore,
			SubjectiveScore: 0,
			TotalScore:      objScore,
			HasSubjective:   hasFlag,
			Updater:         "task",
			UpdateTime:      gtime.Now(),
		})
		if err != nil {
			return err
		}
		return examutil.UpsertFromAttemptTx(ctx, tx, attemptID)
	})
}

func loadQuestionScoreMetaTx(ctx context.Context, tx gdb.TX, paperID int64) ([]bo.QuestionScoreMeta, error) {
	var qs []examentity.ExamQuestion
	if err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).
		Where("exam_paper_id", paperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&qs); err != nil {
		return nil, err
	}
	out := make([]bo.QuestionScoreMeta, 0, len(qs))
	for _, q := range qs {
		var opts []examentity.ExamOption
		_ = tx.Model(dao.ExamOption.Table()).Ctx(ctx).
			Where("question_id", q.Id).
			Where("is_correct", 1).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("sort_order").
			Scan(&opts)
		var correct []int64
		for _, o := range opts {
			correct = append(correct, o.Id)
		}
		out = append(out, bo.QuestionScoreMeta{
			QuestionID:    q.Id,
			IsExample:     q.IsExample,
			IsSubjective:  q.IsSubjective,
			Score:         q.Score,
			CorrectOptIDs: correct,
		})
	}
	return out, nil
}

func loadAnswersMapTx(ctx context.Context, tx gdb.TX, attemptID int64) (map[int64]bo.AnswerPayload, error) {
	var rows []examentity.ExamAttemptAnswer
	if err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
		Where("attempt_id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&rows); err != nil {
		return nil, err
	}
	m := make(map[int64]bo.AnswerPayload, len(rows))
	for _, r := range rows {
		m[r.ExamQuestionId] = examutil.ParseAnswerPayload(r.AnswerJson)
	}
	return m, nil
}
