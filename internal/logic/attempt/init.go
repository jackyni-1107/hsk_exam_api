package attempt

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/service/attempt"
	"exam/internal/utility/examutil"
)

type sAttempt struct{}

func init() {
	attempt.RegisterAttempt(New())
}

func New() *sAttempt {
	return &sAttempt{}
}

// --- 公用方法 ---

// GetAttemptByID 获取答题会话详情
func (s *sAttempt) GetAttemptByID(ctx context.Context, id int64) (*examentity.ExamAttempt, error) {
	var out *examentity.ExamAttempt
	err := dao.ExamAttempt.Ctx(ctx).Where("id", id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&out)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, gerror.NewCode(consts.CodeDataNotFound, "答题记录不存在")
	}
	return out, nil
}

func (s *sAttempt) LoadAttemptByUser(ctx context.Context, attemptID int64, userID int64) (*examentity.ExamAttempt, error) {
	att, err := LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (s *sAttempt) AssertAttemptInProgressByUser(ctx context.Context, attemptID int64, userID int64) (*examentity.ExamAttempt, error) {
	return AssertAttemptInProgressByUser(ctx, attemptID, userID)
}

// IsWindowOpen 判断考试窗口是否开启（含起止边界时刻）。
func (s *sAttempt) IsWindowOpen(now, start, end *gtime.Time) bool {
	return isBatchWindowOpen(now, start, end)
}

// getPageSize 统一分页工具
func (s *sAttempt) getPageSize(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return page, size
}

func boolPtr(b bool) *bool {
	return &b
}

// LoadAttemptByUser 按用户加载未删除的答题会话（供本包与 exam 桥接使用）。
func LoadAttemptByUser(ctx context.Context, attemptID, userID int64) (examentity.ExamAttempt, error) {
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

// AssertAttemptInProgressByUser 校验会话归属且为进行中。
func AssertAttemptInProgressByUser(ctx context.Context, attemptID, userID int64) (*examentity.ExamAttempt, error) {
	att, err := LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil, err
	}
	if !canSaveAttemptAnswers(att.Status) {
		return nil, gerror.NewCode(consts.CodeExamAlreadySubmitted)
	}
	return &att, nil
}

func loadCorrectOptionIDsByQuestion(ctx context.Context, qIDs []interface{}) map[int64][]int64 {
	out := make(map[int64][]int64)
	if len(qIDs) == 0 {
		return out
	}
	var opts []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		WhereIn("question_id", qIDs).
		Where("is_correct", 1).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&opts); err != nil {
		g.Log().Warningf(ctx, "loadCorrectOptionIDsByQuestion: %v", err)
		return out
	}
	for _, o := range opts {
		out[o.QuestionId] = append(out[o.QuestionId], o.Id)
	}
	return out
}

func loadBlocksByID(ctx context.Context, blockIDs []interface{}) map[int64]examentity.ExamQuestionBlock {
	out := make(map[int64]examentity.ExamQuestionBlock)
	if len(blockIDs) == 0 {
		return out
	}
	var blocks []examentity.ExamQuestionBlock
	if err := dao.ExamQuestionBlock.Ctx(ctx).
		WhereIn("id", blockIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&blocks); err != nil {
		g.Log().Warningf(ctx, "loadBlocksByID: %v", err)
		return out
	}
	for _, b := range blocks {
		out[b.Id] = b
	}
	return out
}

func loadSectionsByID(ctx context.Context, examPaperId int64, blocks map[int64]examentity.ExamQuestionBlock) map[int64]examentity.ExamSection {
	out := make(map[int64]examentity.ExamSection)
	if len(blocks) == 0 {
		return out
	}
	sectionIDs := make([]interface{}, 0, len(blocks))
	seen := make(map[int64]struct{})
	for _, b := range blocks {
		if _, ok := seen[b.SectionId]; ok {
			continue
		}
		seen[b.SectionId] = struct{}{}
		sectionIDs = append(sectionIDs, b.SectionId)
	}
	if len(sectionIDs) == 0 {
		return out
	}
	var sections []examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		WhereIn("id", sectionIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&sections); err != nil {
		g.Log().Warningf(ctx, "loadSectionsByID: %v", err)
		return out
	}
	for _, sec := range sections {
		out[sec.Id] = sec
	}
	return out
}

func loadOptionsByQuestion(ctx context.Context, qIDs []interface{}) map[int64][]examentity.ExamOption {
	out := make(map[int64][]examentity.ExamOption)
	if len(qIDs) == 0 {
		return out
	}
	var opts []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		WhereIn("question_id", qIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&opts); err != nil {
		g.Log().Warningf(ctx, "loadOptionsByQuestion: %v", err)
		return out
	}
	for _, o := range opts {
		qid := o.QuestionId
		out[qid] = append(out[qid], o)
	}
	return out
}

func loadSegmentDurationSeconds(ctx context.Context, levelID int64, segmentCode string) int {
	if levelID <= 0 || segmentCode == "" {
		return 0
	}
	var seg struct {
		Duration int `json:"duration"`
	}
	if err := dao.MockExaminationSegment.Ctx(ctx).
		Fields("duration").
		Where("level_id", levelID).
		Where("segment_code", segmentCode).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Limit(1).
		Scan(&seg); err != nil {
		g.Log().Warningf(ctx, "loadSegmentDurationSeconds failed: %v", err)
		return 0
	}
	if seg.Duration <= 0 {
		return 0
	}
	return seg.Duration * 60
}

func computeSegmentRemainingSeconds(ctx context.Context, att examentity.ExamAttempt, segmentCode string) *int {
	if !isAttemptInProgress(att.Status) || segmentCode == "" {
		return nil
	}
	segmentDurationSeconds := loadSegmentDurationSeconds(ctx, att.MockLevelId, segmentCode)
	if segmentDurationSeconds <= 0 {
		return nil
	}
	ref := RedisGetSegmentLastSubmitTime(ctx, att.Id, segmentCode)
	if ref == nil {
		ref = att.StartedAt
	}
	if ref == nil {
		return nil
	}
	nowTs := gtime.Now().Timestamp()
	used := nowTs - ref.Timestamp()
	if used < 0 {
		used = 0
	}
	sec := int64(segmentDurationSeconds) - used
	if sec < 0 {
		sec = 0
	}
	x := int(sec)
	return &x
}

func maybeAutoSubmitIfOverdue(ctx context.Context, userID int64, attemptID int64) error {
	att, err := LoadAttemptByUser(ctx, attemptID, userID)
	if err != nil {
		return nil
	}
	if !isAttemptInProgress(att.Status) || att.DeadlineAt == nil {
		return nil
	}
	flags, err := loadExamBatchFlags(ctx, att.ExamBatchId)
	if err != nil {
		return nil
	}
	if !flags.AutoSubmitOnDeadline {
		return nil
	}
	now := gtime.Now()
	if !isAttemptDeadlineReached(att, now) {
		return nil
	}
	return markSubmitted(ctx, attemptID, attemptEventTimeout)
}

func isExamBatchExpired(ctx context.Context, batchID int64, now *gtime.Time) (bool, error) {
	if batchID <= 0 {
		return false, nil
	}
	var row struct {
		ExamEndAt *gtime.Time `orm:"exam_end_at"`
	}
	if err := dao.ExamBatch.Ctx(ctx).
		Fields("exam_end_at").
		Where("id", batchID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Limit(1).
		Scan(&row); err != nil {
		return false, err
	}
	if row.ExamEndAt == nil {
		return false, nil
	}
	return row.ExamEndAt.Before(now), nil
}

func isExamBatchWindowOpen(ctx context.Context, batchID int64, now *gtime.Time) (bool, error) {
	if batchID <= 0 {
		return true, nil
	}
	var row struct {
		ExamStartAt *gtime.Time `orm:"exam_start_at"`
		ExamEndAt   *gtime.Time `orm:"exam_end_at"`
	}
	if err := dao.ExamBatch.Ctx(ctx).
		Fields("exam_start_at", "exam_end_at").
		Where("id", batchID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Limit(1).
		Scan(&row); err != nil {
		return false, err
	}
	return isBatchWindowOpen(now, row.ExamStartAt, row.ExamEndAt), nil
}

func applyAttemptTransitionTx(ctx context.Context, tx gdb.TX, att examentity.ExamAttempt, event attemptEvent, patch examdo.ExamAttempt) (bool, error) {
	nextStatus, ok := transitionAttemptStatus(att.Status, event)
	if att.Id == 0 || !ok {
		return false, nil
	}
	patch.Status = nextStatus
	if patch.UpdateTime == nil {
		patch.UpdateTime = gtime.Now()
	}
	r, err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).
		Where("id", att.Id).
		Where("status", att.Status).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Update(patch)
	if err != nil {
		return false, err
	}
	n, _ := r.RowsAffected()
	return n > 0, nil
}

func markSubmitted(ctx context.Context, attemptID int64, event attemptEvent) error {
	ok, err := AcquireSubmitLockWithRetry(ctx, attemptID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer ReleaseSubmitLock(ctx, attemptID)

	return markSubmittedLocked(ctx, attemptID, event)
}

func markSubmittedLocked(ctx context.Context, attemptID int64, event attemptEvent) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var att examentity.ExamAttempt
		if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&att); err != nil {
			return err
		}
		if att.Id == 0 || !canAttempt(att.Status, event) {
			return nil
		}
		now := gtime.Now()
		if event == attemptEventTimeout {
			if !isAttemptDeadlineReached(att, now) {
				return nil
			}
		}
		if err := syncAttemptDraftsToDBTx(ctx, tx, attemptID); err != nil {
			return err
		}
		_, err := applyAttemptTransitionTx(ctx, tx, att, event, examdo.ExamAttempt{
			SubmittedAt: now,
			UpdateTime:  gtime.Now(),
		})
		return err
	})
}

func finalizeScoring(ctx context.Context, attemptID int64) error {
	ok, err := AcquireSubmitLockWithRetry(ctx, attemptID)
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
		if isAttemptScored(att.Status) {
			return nil
		}
		if !canAttempt(att.Status, attemptEventFinalize) {
			return nil
		}
		if err := syncAttemptDraftsToDBTx(ctx, tx, attemptID); err != nil {
			return err
		}

		meta, err := loadQuestionScoreMetaTx(ctx, tx, att.ExamPaperId)
		if err != nil {
			return err
		}

		skipScoring := false
		if att.ExamBatchId > 0 {
			var brow struct {
				SkipScoring int `orm:"skip_scoring"`
			}
			if err := tx.Model(dao.ExamBatch.Table()).Ctx(ctx).
				Fields("skip_scoring").
				Where("id", att.ExamBatchId).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&brow); err != nil {
				return err
			}
			skipScoring = brow.SkipScoring != 0
		}

		now := gtime.Now()
		hasFlag := 0
		for _, m := range meta {
			if m.IsSubjective != 0 && m.IsExample == 0 {
				hasFlag = 1
				break
			}
		}

		if skipScoring {
			applied, err := applyAttemptTransitionTx(ctx, tx, att, attemptEventFinalize, examdo.ExamAttempt{
				EndedAt:          now,
				ObjectiveScore:   0,
				SubjectiveScore:  0,
				TotalScore:       0,
				SegmentScoreJson: "{}",
				HasSubjective:    hasFlag,
				UpdateTime:       gtime.Now(),
			})
			if err != nil {
				return err
			}
			if !applied {
				return nil
			}
			return examutil.UpsertFromAttemptTx(ctx, tx, attemptID)
		}

		answers, err := loadAnswersMapTx(ctx, tx, attemptID)
		if err != nil {
			return err
		}
		awardedByQ, err := loadAwardedScoresMapTx(ctx, tx, attemptID)
		if err != nil {
			return err
		}
		segmentScores, totalScore, hasSubj := examutil.ScoreBySegment(meta, answers, awardedByQ)
		segmentScoreJSON := marshalSegmentScores(segmentScores)
		objScore := float64(totalScore)
		if hasSubj {
			hasFlag = 1
		}
		applied, err := applyAttemptTransitionTx(ctx, tx, att, attemptEventFinalize, examdo.ExamAttempt{
			EndedAt:          now,
			ObjectiveScore:   objScore,
			SubjectiveScore:  0,
			TotalScore:       objScore,
			SegmentScoreJson: segmentScoreJSON,
			HasSubjective:    hasFlag,
			UpdateTime:       gtime.Now(),
		})
		if err != nil {
			return err
		}
		if !applied {
			return nil
		}
		return examutil.UpsertFromAttemptTx(ctx, tx, attemptID)
	})
}

func syncAttemptDraftsToDBTx(ctx context.Context, tx gdb.TX, attemptID int64) error {
	draftMap, err := RedisHGetAllAttemptDrafts(ctx, attemptID)
	if err != nil {
		return err
	}
	if len(draftMap) == 0 {
		return nil
	}
	items := examutil.BuildAttemptAnswerDraftRows(attemptID, draftMap)
	if len(items) == 0 {
		return nil
	}
	return examutil.UpsertAttemptAnswerDraftRowsTx(ctx, tx, items)
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
	blockIDs := make([]interface{}, 0, len(qs))
	seenBlock := make(map[int64]struct{}, len(qs))
	for _, q := range qs {
		if _, ok := seenBlock[q.BlockId]; ok {
			continue
		}
		seenBlock[q.BlockId] = struct{}{}
		blockIDs = append(blockIDs, q.BlockId)
	}
	segmentByBlock, err := loadSegmentCodeByBlockTx(ctx, tx, blockIDs)
	if err != nil {
		return nil, err
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
			SegmentCode:   segmentByBlock[q.BlockId],
			IsExample:     q.IsExample,
			IsSubjective:  q.IsSubjective,
			Score:         q.Score,
			CorrectOptIDs: correctByQ[q.Id],
		})
	}
	return out, nil
}

func loadSegmentCodeByBlockTx(ctx context.Context, tx gdb.TX, blockIDs []interface{}) (map[int64]string, error) {
	out := make(map[int64]string)
	if len(blockIDs) == 0 {
		return out, nil
	}
	var blocks []struct {
		Id        int64 `orm:"id"`
		SectionId int64 `orm:"section_id"`
	}
	if err := tx.Model(dao.ExamQuestionBlock.Table()).Ctx(ctx).
		Fields("id", "section_id").
		WhereIn("id", blockIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&blocks); err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		return out, nil
	}
	sectionIDs := make([]interface{}, 0, len(blocks))
	seenSection := make(map[int64]struct{}, len(blocks))
	for _, b := range blocks {
		if _, ok := seenSection[b.SectionId]; ok {
			continue
		}
		seenSection[b.SectionId] = struct{}{}
		sectionIDs = append(sectionIDs, b.SectionId)
	}
	segmentBySection := make(map[int64]string, len(sectionIDs))
	var sections []struct {
		Id          int64  `orm:"id"`
		SegmentCode string `orm:"segment_code"`
	}
	if err := tx.Model(dao.ExamSection.Table()).Ctx(ctx).
		Fields("id", "segment_code").
		WhereIn("id", sectionIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&sections); err != nil {
		return nil, err
	}
	for _, sec := range sections {
		segmentBySection[sec.Id] = sec.SegmentCode
	}
	for _, b := range blocks {
		out[b.Id] = segmentBySection[b.SectionId]
	}
	return out, nil
}

func loadAnswersMapTx(ctx context.Context, tx gdb.TX, attemptID int64) (map[int64]bo.AnswerPayload, error) {
	var rows []examentity.ExamAttemptAnswer
	if err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
		Where("attempt_id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("exam_question_id").
		OrderAsc("version").
		OrderAsc("update_time").
		OrderAsc("id").
		Scan(&rows); err != nil {
		return nil, err
	}
	m := make(map[int64]bo.AnswerPayload, len(rows))
	for _, r := range rows {
		m[r.ExamQuestionId] = examutil.ParseAnswerPayload(r.AnswerJson)
	}
	return m, nil
}

func loadAwardedScoresMapTx(ctx context.Context, tx gdb.TX, attemptID int64) (map[int64]float64, error) {
	out := make(map[int64]float64)
	var rows []struct {
		ExamQuestionId int64    `orm:"exam_question_id"`
		AwardedScore   *float64 `orm:"awarded_score"`
	}
	if err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
		Fields("exam_question_id", "awarded_score").
		Where("attempt_id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&rows); err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r.AwardedScore == nil {
			continue
		}
		out[r.ExamQuestionId] = *r.AwardedScore
	}
	return out, nil
}

func marshalSegmentScores(segmentScores map[string]int) string {
	if len(segmentScores) == 0 {
		return "{}"
	}
	b, err := json.Marshal(segmentScores)
	if err != nil {
		return "{}"
	}
	return string(b)
}
