package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	attemptsvc "exam/internal/service/attempt"
)

const (
	defaultExamBatchExpireBatchLimit = 200
	maxExamBatchExpireBatchLimit     = 2000
)

// examBatchExpireParams 任务 params JSON，字段均可选。
// 示例：{"batch_limit":100,"exam_batch_id":12}
type examBatchExpireParams struct {
	BatchLimit  int   `json:"batch_limit"`
	ExamBatchID int64 `json:"exam_batch_id"`
}

type examBatchExpireHandler struct{}

func (examBatchExpireHandler) Execute(ctx context.Context, taskID int64, params string) error {
	p := examBatchExpireParams{BatchLimit: defaultExamBatchExpireBatchLimit}
	if s := strings.TrimSpace(params); s != "" {
		if err := json.Unmarshal([]byte(s), &p); err != nil {
			return gerror.Wrap(err, "ExamBatchExpireHandler: invalid params JSON")
		}
	}
	limit := p.BatchLimit
	if limit <= 0 {
		limit = defaultExamBatchExpireBatchLimit
	}
	if limit > maxExamBatchExpireBatchLimit {
		limit = maxExamBatchExpireBatchLimit
	}

	model := dao.ExamAttempt.Ctx(ctx).
		Fields("exam_attempt.id").
		LeftJoin("exam_batch", "exam_batch.id=exam_attempt.exam_batch_id").
		Where("exam_attempt.status", consts.ExamAttemptInProgress).
		Where("exam_attempt.delete_flag", consts.DeleteFlagNotDeleted).
		Where("exam_batch.delete_flag", consts.DeleteFlagNotDeleted).
		WhereLT("exam_batch.exam_end_at", gtime.Now())
	if p.ExamBatchID > 0 {
		model = model.Where("exam_attempt.exam_batch_id", p.ExamBatchID)
	}

	var rows []struct {
		Id int64 `json:"id"`
	}
	if err := model.Limit(limit).OrderAsc("exam_attempt.id").Scan(&rows); err != nil {
		return gerror.Wrap(err, "ExamBatchExpireHandler: list expired attempts")
	}
	if len(rows) == 0 {
		g.Log().Infof(ctx, "[ExamBatchExpireHandler] task_id=%d no in-progress attempts to submit (limit=%d batch_id=%d)",
			taskID, limit, p.ExamBatchID)
		return nil
	}

	failN := 0
	for _, row := range rows {
		if err := attemptsvc.Attempt().MarkSubmittedByBatchExpired(ctx, row.Id); err != nil {
			failN++
			g.Log().Warningf(ctx, "[ExamBatchExpireHandler] task_id=%d attempt_id=%d mark submitted error: %v", taskID, row.Id, err)
		}
	}
	okN := len(rows) - failN
	g.Log().Infof(ctx, "[ExamBatchExpireHandler] task_id=%d done ok=%d fail=%d (scanned=%d limit=%d exam_batch_id=%d)",
		taskID, okN, failN, len(rows), limit, p.ExamBatchID)
	if failN > 0 {
		g.Log().Warningf(ctx, "[ExamBatchExpireHandler] task_id=%d partial failed=%d total=%d; skipped failed items and continue",
			taskID, failN, len(rows))
	}
	return nil
}
