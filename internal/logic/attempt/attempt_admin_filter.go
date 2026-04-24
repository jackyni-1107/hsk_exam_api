package attempt

import (
	"strings"

	"exam/internal/consts"
)

// AttemptAdminListQuery 管理端 exam_result 列表/统计共用的联表筛选（与 AttemptAdminList 原语义一致）。
type AttemptAdminListQuery struct {
	// Level 按 exam_paper.level 字符串筛（兼容旧参数）；若 MockLevelId>0 则优先按 r.mock_level_id
	Level string
	// MockLevelId 对应 mock_levels.id，与 exam_result.mock_level_id 一致
	MockLevelId        int64
	ExaminationPaperId int64
	ExamBatchId        int64
	Status             int // 0: 不限
	Username           string
	// SubjectivePending=1: 含主观题且 exam_result.status=4（已结束、待主观评阅，与 status=5 已完成算分区分）
	SubjectivePending int
}

// attemptAdminListJoin 返回与列表相同的 FROM 片段（至 WHERE 之前不含 WHERE 关键字）及 JOIN 侧绑定参数。
func attemptAdminListJoin() (from string, joinArgs []interface{}) {
	joinArgs = []interface{}{
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
	}
	from = ` FROM exam_result r
INNER JOIN exam_attempt a ON a.id = r.attempt_id AND a.delete_flag = ?
LEFT JOIN sys_member u ON u.id = r.member_id AND u.delete_flag = ?
LEFT JOIN exam_paper p ON p.id = r.exam_paper_id AND p.delete_flag = ?
LEFT JOIN mock_examination_paper m ON m.id = p.mock_examination_paper_id AND m.delete_flag = ?
LEFT JOIN mock_levels ml ON ml.id = r.mock_level_id AND ml.delete_flag = ?
WHERE `
	return from, joinArgs
}

// buildAttemptAdminWhere 生成 WHERE 子句（已含 "r.delete_flag = ?" 起头）与参数；不含 JOIN 侧参数。
func (q AttemptAdminListQuery) buildAttemptAdminWhere() (where string, args []interface{}) {
	var w strings.Builder
	w.WriteString("r.delete_flag = ?")
	args = []interface{}{consts.DeleteFlagNotDeleted}
	if q.MockLevelId > 0 {
		w.WriteString(" AND r.mock_level_id = ?")
		args = append(args, q.MockLevelId)
	} else if q.Level != "" {
		w.WriteString(" AND p.level = ?")
		args = append(args, q.Level)
	}
	if q.ExaminationPaperId > 0 {
		w.WriteString(" AND p.mock_examination_paper_id = ?")
		args = append(args, q.ExaminationPaperId)
	}
	if q.ExamBatchId > 0 {
		w.WriteString(" AND r.exam_batch_id = ?")
		args = append(args, q.ExamBatchId)
	}
	if q.Status > 0 {
		w.WriteString(" AND r.status = ?")
		args = append(args, q.Status)
	}
	if q.Username != "" {
		w.WriteString(" AND u.username LIKE ?")
		args = append(args, "%"+q.Username+"%")
	}
	if q.SubjectivePending == 1 {
		w.WriteString(" AND r.has_subjective = 1 AND r.status = ?")
		args = append(args, consts.ExamAttemptEnded)
	}
	return w.String(), args
}

// buildAttemptAdminListFrom 拼接 `FROM ... WHERE <cond>`，与 count/list SQL 组装一致。
func (q AttemptAdminListQuery) buildAttemptAdminListFrom() (fromSQL string, joinArgs []interface{}, whereArgs []interface{}) {
	from, joinArgs := attemptAdminListJoin()
	where, wArgs := q.buildAttemptAdminWhere()
	fromSQL = from + where
	return fromSQL, joinArgs, wArgs
}

// attemptAdminListCountArgs 将 JOIN 与 WHERE 参数按 Raw 顺序合并。
func attemptAdminListCountArgs(joinArgs, whereArgs []interface{}) []interface{} {
	return append(append([]interface{}{}, joinArgs...), whereArgs...)
}
