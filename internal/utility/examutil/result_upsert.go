// Package examutil 含考试域通用逻辑。
//
// 时间字段约定：exam_attempt 为 started_at / submitted_at / ended_at（及 deadline_at）的唯一权威来源；
// exam_result 中对应列为查询优化用的冗余快照，仅允许通过 UpsertFromAttemptTx 在更新 attempt 的同一事务内写入，
// 禁止其它路径直接改 exam_result 的时间列。展示层在已 JOIN attempt 时应优先读 a.*（见 attempt_admin 列表 SQL）。
package examutil

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
)

// UpsertFromAttemptTx 将会话快照同步到 exam_result（交卷/计分后调用）。
func UpsertFromAttemptTx(ctx context.Context, tx gdb.TX, attemptID int64) error {
	var att examentity.ExamAttempt
	if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&att); err != nil {
		return err
	}
	if att.Id == 0 {
		return nil
	}
	row := examdo.ExamResult{
		AttemptId:              att.Id,
		MemberId:               att.MemberId,
		ExamPaperId:            att.ExamPaperId,
		MockExaminationPaperId: att.MockExaminationPaperId,
		ExamBatchId:            att.ExamBatchId,
		MockLevelId:            att.MockLevelId,
		Status:                 att.Status,
		ObjectiveScore:         att.ObjectiveScore,
		SubjectiveScore:        att.SubjectiveScore,
		TotalScore:             att.TotalScore,
		HasSubjective:          att.HasSubjective,
		StartedAt:              att.StartedAt,
		SubmittedAt:            att.SubmittedAt,
		EndedAt:                att.EndedAt,
		CreateTime:             att.CreateTime,
		UpdateTime:             gtime.Now(),
		DeleteFlag:             consts.DeleteFlagNotDeleted,
	}
	var exist examentity.ExamResult
	_ = tx.Model(dao.ExamResult.Table()).Ctx(ctx).Where("attempt_id", attemptID).Scan(&exist)
	if exist.AttemptId == 0 {
		_, err := tx.Model(dao.ExamResult.Table()).Ctx(ctx).Insert(row)
		return err
	}
	_, err := tx.Model(dao.ExamResult.Table()).Ctx(ctx).Where("attempt_id", attemptID).Data(row).Update()
	return err
}
