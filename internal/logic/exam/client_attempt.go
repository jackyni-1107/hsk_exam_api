package exam

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"
	"exam/internal/utility"
	"exam/internal/utility/exampaper"
	"exam/internal/utility/examutil"
)

// CreateAttempt 已废弃：请使用 CreateAttemptForBatch（POST /exam/batches/{batchId}/attempts）。
func CreateAttempt(ctx context.Context, userID int64, mockPaperID int64) (int64, error) {
	_ = ctx
	_ = userID
	_ = mockPaperID
	return 0, gerror.NewCode(consts.CodeExamAttemptUseBatchApi)
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

func loadAttemptByUser(ctx context.Context, attemptID, userID int64) (examentity.ExamAttempt, error) {
	var att examentity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("member_id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att); err != nil {
		return att, err
	}
	if att.Id == 0 {
		return att, gerror.NewCode(consts.CodeExamAttemptNotFound)
	}
	return att, nil
}

func assertAttemptInProgressByUser(ctx context.Context, attemptID, userID int64) (*examentity.ExamAttempt, error) {
	att, err := loadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	if att.Status != consts.ExamAttemptInProgress {
		return nil, gerror.NewCode(consts.CodeExamAlreadySubmitted)
	}
	now := gtime.Now()
	if att.DeadlineAt != nil && att.DeadlineAt.Before(now) {
		return nil, gerror.NewCode(consts.CodeExamTimeExpired)
	}
	return &att, nil
}

// CreateAttemptForBatch 按批次与 Mock 卷创建会话（未开始）；每用户每批次每卷仅允许一条未删除记录。
func CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64) (int64, error) {
	var batch examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).
		Where("id", batchID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&batch); err != nil {
		return 0, err
	}
	if batch.Id == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchNotFound)
	}
	now := gtime.Now()
	if !batchExamWindowOpen(now, batch.ExamStartAt, batch.ExamEndAt) {
		return 0, gerror.NewCode(consts.CodeExamBatchWindowNotOpen)
	}

	var link examentity.ExamBatchMember
	if err := dao.ExamBatchMember.Ctx(ctx).
		Where("batch_id", batchID).
		Where("member_id", userID).
		Limit(1).
		Scan(&link); err != nil {
		return 0, err
	}
	if link.BatchId == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchMemberNotFound)
	}
	paper, err := exampaper.ByMockID(ctx, link.MockExaminationPaperId)
	if err != nil {
		return 0, err
	}
	var mp mockentity.MockExaminationPaper
	_ = dao.MockExaminationPaper.Ctx(ctx).Where("id", link.MockExaminationPaperId).Scan(&mp)
	levelID := mp.LevelId
	attemptVar, err := dao.ExamAttempt.Ctx(ctx).
		Fields("id").
		Where("member_id", userID).
		Where("exam_batch_id", batchID).
		Where("mock_examination_paper_id", link.MockExaminationPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Value()
	if err != nil {
		return 0, err
	}
	attemptId := attemptVar.Int64()
	if attemptId > 0 {
		return attemptId, nil
	}
	id, err := dao.ExamAttempt.Ctx(ctx).InsertAndGetId(examdo.ExamAttempt{
		MemberId:               userID,
		ExamPaperId:            paper.Id,
		MockExaminationPaperId: paper.MockExaminationPaperId,
		ExamBatchId:            batchID,
		MockLevelId:            levelID,
		Status:                 consts.ExamAttemptNotStarted,
		DurationSeconds:        0,
		Creator:                updaterClient,
		Updater:                updaterClient,
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
		return gerror.NewCode(consts.CodeExamAttemptNotFound)
	}
	if att.Status != consts.ExamAttemptNotStarted {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	var paper examentity.ExamPaper
	if err := dao.ExamPaper.Ctx(ctx).Where("id", att.ExamPaperId).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&paper); err != nil {
		return err
	}
	if paper.Id == 0 {
		return gerror.NewCode(consts.CodeExamPaperNotFound)
	}
	dur := ResolveDurationSeconds(cfg, paper.DurationSeconds, clientDurationSeconds)
	now := gtime.Now()
	deadline := gtime.NewFromTimeStamp(now.Timestamp() + int64(dur))

	_, err = dao.ExamAttempt.Ctx(ctx).Where("id", att.Id).Update(examdo.ExamAttempt{
		Status:          consts.ExamAttemptInProgress,
		DurationSeconds: dur,
		StartedAt:       now,
		DeadlineAt:      deadline,
		Updater:         updaterClient,
		UpdateTime:      gtime.Now(),
	})
	return err
}

// GetAttempt 查询会话；若已超时仍进行中则自动交卷并计分。
func GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error) {
	//_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)
	att, err := loadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	now := gtime.Now()
	deadlineReached := att.Status == consts.ExamAttemptInProgress && att.DeadlineAt != nil && att.DeadlineAt.Before(now)
	rem := computeBatchExamRemainingSeconds(ctx, att)
	return &bo.AttemptView{
		Attempt:          att,
		ServerTime:       utility.ToRFC3339UTCShift(now),
		DeadlineReached:  deadlineReached,
		RemainingSeconds: rem,
	}, nil
}

// computeBatchExamRemainingSeconds 以批次 exam_end_at 为结束时刻，以最近一次保存答案时间（库表与 Redis 草稿取较新，均无则退回开考时间）为参考时刻，返回剩余秒数（仅进行中且批次有效时）。
func computeBatchExamRemainingSeconds(ctx context.Context, att examentity.ExamAttempt) *int {
	if att.Status != consts.ExamAttemptInProgress || att.ExamBatchId <= 0 {
		return nil
	}
	var batch examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).
		Where("id", att.ExamBatchId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&batch); err != nil || batch.Id == 0 || batch.ExamEndAt == nil {
		return nil
	}
	ref := lastAnswerReferenceTime(ctx, att.Id, att.StartedAt)
	if ref == nil {
		return nil
	}
	sec := batch.ExamEndAt.Timestamp() - ref.Timestamp()
	if sec < 0 {
		sec = 0
	}
	x := int(sec)
	return &x
}

func lastAnswerReferenceTime(ctx context.Context, attemptID int64, startedAt *gtime.Time) *gtime.Time {
	var row examentity.ExamAttemptAnswer
	_ = dao.ExamAttemptAnswer.Ctx(ctx).
		Fields("update_time").
		Where("attempt_id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderDesc("update_time").
		Limit(1).
		Scan(&row)
	var best *gtime.Time
	if row.UpdateTime != nil {
		best = row.UpdateTime
	}
	if rds := maxAnswerSaveTimeFromRedis(ctx, attemptID); rds != nil {
		if best == nil || rds.After(best) {
			best = rds
		}
	}
	if best == nil {
		best = startedAt
	}
	return best
}

func maxAnswerSaveTimeFromRedis(ctx context.Context, attemptID int64) *gtime.Time {
	redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)
	res, err := g.Redis().HGetAll(ctx, redisKey)
	if err != nil || res.IsEmpty() {
		return nil
	}
	var maxTs int64
	for _, s := range res.Map() {
		t := gjson.New(s).Get("t").Int64()
		if t > maxTs {
			maxTs = t
		}
	}
	if maxTs <= 0 {
		return nil
	}
	return gtime.NewFromTimeStamp(maxTs)
}

// SaveAnswers 保存答案 redis -> db
func SaveAnswers(ctx context.Context, userID int64, attemptID int64, items []bo.SaveAnswerItem) error {
	// 1. 基础校验 (限流、考试状态)
	cfg := LoadExamCfg(ctx)
	if err := RateLimitSaveAnswers(ctx, attemptID, cfg.SaveAnswersPerSecond); err != nil {
		return err
	}

	// 校验考试是否在进行中
	//_, err := assertAttemptInProgressByUser(ctx, attemptID, userID)
	//if err != nil {
	//	return err
	//}

	if len(items) == 0 {
		return nil
	}

	// 2. 准备 Redis 数据
	redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)
	data := make(map[string]interface{})
	now := gtime.Now()

	for _, it := range items {
		// 封装数据。注意：此处由于是异步，乐观锁版本校验可以先放在 Redis 层面或落库层面
		val := g.Map{
			"q": it.QuestionID,
			"a": it.AnswerJSON,
			"v": it.ExpectedVersion, // 记录调用方预期的版本
			"t": now.Timestamp(),
		}
		data[gconv.String(it.QuestionID)] = val
	}

	// 3. 写入 Redis Hash
	err := g.Redis().HMSet(ctx, redisKey, data)
	if err != nil {
		g.Log().Error(ctx, "Redis写入失败", err)
		return err
	}

	// 4. 设置过期时间（例如 2 小时）
	_, _ = g.Redis().Expire(ctx, redisKey, 7200)

	// 5. 推送到队列（通知异步落库）
	// 我们只需要把 attemptID 放进去，消费者就知道去哪个 Hash 拿数据
	if _, err := g.Redis().LPush(ctx, consts.ExamAttemptSyncQueueKey, attemptID); err != nil {
		g.Log().Error(ctx, "推送队列失败", err)
	}

	return nil
}

// SubmitAttempt 主动交卷：仅标记为已交卷（待算分）。客观分与 exam_result 由 sys_task（ExamScoreFinalizeHandler）统一算分写入。
func SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error {
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)
	att, err := loadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return err
	}
	if att.Status == consts.ExamAttemptEnded {
		return nil
	}
	if att.Status == consts.ExamAttemptSubmitted {
		return nil
	}
	if att.Status != consts.ExamAttemptInProgress {
		return gerror.NewCode(consts.CodeExamAlreadySubmitted)
	}
	return markSubmitted(ctx, attemptID, false, updaterClient)
}

// maybeAutoSubmitIfOverdue 考试时间到达且仍为进行中时自动标记已交卷（待算分）。算分仅由定时任务执行。
func maybeAutoSubmitIfOverdue(ctx context.Context, userID int64, attemptID int64) error {
	att, err := loadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil
	}
	if att.Status != consts.ExamAttemptInProgress || att.DeadlineAt == nil {
		return nil
	}
	now := gtime.Now()
	if !att.DeadlineAt.Before(now) {
		return nil
	}
	return markSubmitted(ctx, attemptID, true, updaterClient)
}

// MarkSubmittedIfOverdue 供定时任务：超时未操作会话标记为已交卷（待算分，不校验用户）。算分由 ExamScoreFinalizeHandler 执行。
func MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error {
	return markSubmitted(ctx, attemptID, true, updaterTask)
}

// FinalizeAttempt 对已交卷（待算分）会话计算客观分并置为已结束，写入 exam_result。仅应由 ExamScoreFinalizeHandler（sys_task）调用。
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
			return gerror.NewCode(consts.CodeExamAttemptNotFound)
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
			Updater:         updaterTask,
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
	if len(qs) == 0 {
		return nil, nil
	}
	qIDs := make([]interface{}, len(qs))
	for i, q := range qs {
		qIDs[i] = q.Id
	}
	var correctOpts []examentity.ExamOption
	if err := tx.Model(dao.ExamOption.Table()).Ctx(ctx).
		WhereIn("question_id", qIDs).
		Where("is_correct", 1).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&correctOpts); err != nil {
		return nil, err
	}
	correctByQ := make(map[int64][]int64, len(qs))
	for _, o := range correctOpts {
		correctByQ[o.QuestionId] = append(correctByQ[o.QuestionId], o.Id)
	}
	out := make([]bo.QuestionScoreMeta, 0, len(qs))
	for _, q := range qs {
		out = append(out, bo.QuestionScoreMeta{
			QuestionID:    q.Id,
			IsExample:     q.IsExample,
			IsSubjective:  q.IsSubjective,
			Score:         q.Score,
			CorrectOptIDs: correctByQ[q.Id],
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
