package examutil

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
)

const sqlSubjectiveAwardedCount = `SELECT COUNT(1) AS c FROM exam_attempt_answer eaa
INNER JOIN exam_question eq ON eq.id = eaa.exam_question_id
  AND eq.exam_paper_id = ? AND eq.is_subjective = 1 AND eq.is_example = 0
  AND eq.delete_flag = ?
WHERE eaa.attempt_id = ? AND eaa.delete_flag = ? AND eaa.awarded_score IS NOT NULL`

// HasSubjectiveAwarded 会话是否已有任意主观题写入 awarded_score（与管理端评阅语义一致）。
func HasSubjectiveAwarded(ctx context.Context, attemptID, paperID int64) (bool, error) {
	if attemptID <= 0 || paperID <= 0 {
		return false, nil
	}
	var row struct {
		C int `json:"c" orm:"c"`
	}
	if err := g.DB().Ctx(ctx).Raw(sqlSubjectiveAwardedCount,
		paperID, consts.DeleteFlagNotDeleted, attemptID, consts.DeleteFlagNotDeleted,
	).Scan(&row); err != nil {
		return false, err
	}
	return row.C > 0, nil
}

// HasSubjectiveAwardedTx 同 HasSubjectiveAwarded，在事务内执行。
func HasSubjectiveAwardedTx(ctx context.Context, tx gdb.TX, attemptID, paperID int64) (bool, error) {
	if attemptID <= 0 || paperID <= 0 {
		return false, nil
	}
	var row struct {
		C int `json:"c" orm:"c"`
	}
	if err := tx.Ctx(ctx).Raw(sqlSubjectiveAwardedCount,
		paperID, consts.DeleteFlagNotDeleted, attemptID, consts.DeleteFlagNotDeleted,
	).Scan(&row); err != nil {
		return false, err
	}
	return row.C > 0, nil
}
