package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
	"exam/internal/dao"
	attemptsvc "exam/internal/service/attempt"
)

const (
	defaultExamScoreFinalizeBatchLimit = 200
	maxExamScoreFinalizeBatchLimit     = 2000
)

// examScoreFinalizeParams 任务 params JSON，字段均可选。
// 示例：{"batch_limit":100,"exam_batch_id":12}
type examScoreFinalizeParams struct {
	BatchLimit  int   `json:"batch_limit"`
	ExamBatchID int64 `json:"exam_batch_id"`
}

type examScoreFinalizeHandler struct{}

func (examScoreFinalizeHandler) Execute(ctx context.Context, taskID int64, params string) error {
	p := examScoreFinalizeParams{BatchLimit: defaultExamScoreFinalizeBatchLimit}
	if s := strings.TrimSpace(params); s != "" {
		if err := json.Unmarshal([]byte(s), &p); err != nil {
			return gerror.Wrap(err, "ExamScoreFinalizeHandler: invalid params JSON")
		}
	}
	limit := p.BatchLimit
	if limit <= 0 {
		limit = defaultExamScoreFinalizeBatchLimit
	}
	if limit > maxExamScoreFinalizeBatchLimit {
		limit = maxExamScoreFinalizeBatchLimit
	}

	model := dao.ExamAttempt.Ctx(ctx).
		Where("status", consts.ExamAttemptSubmitted).
		Where("delete_flag", consts.DeleteFlagNotDeleted)
	if p.ExamBatchID > 0 {
		model = model.Where("exam_batch_id", p.ExamBatchID)
	}

	var rows []struct {
		Id int64 `json:"id"`
	}
	if err := model.Limit(limit).OrderAsc("id").Scan(&rows); err != nil {
		return gerror.Wrap(err, "ExamScoreFinalizeHandler: list submitted attempts")
	}
	if len(rows) == 0 {
		g.Log().Infof(ctx, "[ExamScoreFinalizeHandler] task_id=%d no submitted attempts to finalize (limit=%d batch_id=%d)",
			taskID, limit, p.ExamBatchID)
		return nil
	}

	var firstErr error
	failN := 0
	for _, row := range rows {
		if err := attemptsvc.Attempt().FinalizeAttempt(ctx, row.Id); err != nil {
			failN++
			if firstErr == nil {
				firstErr = err
			}
			g.Log().Warningf(ctx, "[ExamScoreFinalizeHandler] task_id=%d attempt_id=%d finalize error: %v", taskID, row.Id, err)
		}
	}
	okN := len(rows) - failN
	g.Log().Infof(ctx, "[ExamScoreFinalizeHandler] task_id=%d done ok=%d fail=%d (scanned=%d limit=%d exam_batch_id=%d)",
		taskID, okN, failN, len(rows), limit, p.ExamBatchID)
	if firstErr != nil {
		return gerror.Wrapf(firstErr, "ExamScoreFinalizeHandler: %d of %d attempts failed", failN, len(rows))
	}
	return nil
}
