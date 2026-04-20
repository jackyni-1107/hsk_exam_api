package attempt

import (
	"context"
	"sort"

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
	_ = mockPaperID
	_ = userID
	return 0, gerror.NewCode(consts.CodeExamAttemptUseBatchApi)
}

// CreateAttemptForBatch 按批次与 Mock 卷创建会话（未开始）；每用户每批次每卷仅允许一条未删除记录。
func (s *sAttempt) CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64) (int64, error) {
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
	if !s.IsWindowOpen(now, batch.ExamStartAt, batch.ExamEndAt) {
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
func (s *sAttempt) StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error {
	cfg := s.loadExamSessionCfg(ctx)
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
	dur := resolveDurationSeconds(cfg, paper.DurationSeconds, clientDurationSeconds)
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
func (s *sAttempt) GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error) {
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)
	att, err := LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	now := gtime.Now()
	//deadlineReached := att.Status == consts.ExamAttemptInProgress && att.DeadlineAt != nil && att.DeadlineAt.Before(now)
	finalSegmentCode := RedisLatestSegmentCode(ctx, att.Id)
	rem := computeSegmentRemainingSeconds(ctx, att, finalSegmentCode)
	return &bo.AttemptView{
		Attempt:    att,
		ServerTime: utility.ToRFC3339UTCShift(now),
		//DeadlineReached:  deadlineReached,
		SegmentCode:      finalSegmentCode,
		RemainingSeconds: rem,
	}, nil
}

// GetAttemptAnswers 返回当前用户该会话下的答题明细：先读库再合并 Redis 草稿（与保存路径一致），仅包含非空答案。
func (s *sAttempt) GetAttemptAnswers(ctx context.Context, userID int64, attemptID int64) ([]bo.AttemptAnswerClientItem, error) {
	_, err := LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	byQ := make(map[int64]string)
	var rows []examentity.ExamAttemptAnswer
	if err := dao.ExamAttemptAnswer.Ctx(ctx).
		Where("attempt_id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&rows); err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.ExamQuestionId > 0 && r.AnswerJson != "" {
			byQ[r.ExamQuestionId] = r.AnswerJson
		}
	}
	if draftMap, err := RedisHGetAllAttemptDrafts(ctx, attemptID); err == nil && len(draftMap) > 0 {
		for _, val := range draftMap {
			itemMap := gconv.Map(val)
			q := gconv.Int64(itemMap["q"])
			if q == 0 {
				continue
			}
			if a, ok := itemMap["a"]; ok {
				byQ[q] = gconv.String(a)
			}
		}
	}
	out := make([]bo.AttemptAnswerClientItem, 0, len(byQ))
	for qid, jsonStr := range byQ {
		p := examutil.ParseAnswerPayload(jsonStr)
		if p.OptionID == 0 && p.Text == "" {
			continue
		}
		out = append(out, bo.AttemptAnswerClientItem{
			QuestionID: qid,
			OptionID:   p.OptionID,
			Text:       p.Text,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].QuestionID < out[j].QuestionID })
	return out, nil
}

// SaveAnswers 保存答案 redis -> db
func (s *sAttempt) SaveAnswers(ctx context.Context, userID int64, attemptID int64, segmentCode string, items []bo.SaveAnswerItem) error {
	cfg := s.loadExamSessionCfg(ctx)
	if err := RateLimitSaveAnswers(ctx, attemptID, cfg.SaveAnswersPerSecond); err != nil {
		return err
	}
	_, err := AssertAttemptInProgressByUser(ctx, attemptID, userID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	data := make(map[string]interface{})
	now := gtime.Now()
	for _, it := range items {
		answerJSON := examutil.MarshalAnswerPayload(bo.AnswerPayload{
			OptionID: it.OptionID,
			Text:     it.Text,
		})
		if answerJSON == "" {
			continue
		}
		val := g.Map{
			"q": it.QuestionID,
			"a": answerJSON,
			"v": it.ExpectedVersion,
			"t": now.Timestamp(),
		}
		data[gconv.String(it.QuestionID)] = val
	}
	if len(data) == 0 {
		return nil
	}
	if err := RedisSaveAnswerDrafts(ctx, attemptID, data); err != nil {
		return err
	}
	return RedisSaveSegmentSubmitTime(ctx, attemptID, segmentCode, now.Timestamp())
}

// SubmitAttempt 主动交卷：仅标记为已交卷（待算分）。客观分与 exam_result 由 sys_task（ExamScoreFinalizeHandler）统一算分写入。
func (s *sAttempt) SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error {
	_ = maybeAutoSubmitIfOverdue(ctx, userID, attemptID)
	att, err := LoadAttemptByUser(ctx, attemptID, userID)
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

// MarkSubmittedIfOverdue 供定时任务：超时未操作会话标记为已交卷（待算分，不校验用户）。算分由 ExamScoreFinalizeHandler 执行。
func (s *sAttempt) MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error {
	return markSubmitted(ctx, attemptID, true, updaterTask)
}

// FinalizeAttempt 对已交卷（待算分）会话计算客观分并置为已结束，写入 exam_result。仅应由 ExamScoreFinalizeHandler（sys_task）调用。
func (s *sAttempt) FinalizeAttempt(ctx context.Context, attemptID int64) error {
	return finalizeScoring(ctx, attemptID)
}
