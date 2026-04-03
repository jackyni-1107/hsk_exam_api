package examresult

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/do/exam"
	"exam/internal/model/entity"
)

// UpsertFromAttemptTx 将会话快照同步到 exam_result（交卷/计分后调用）。
func UpsertFromAttemptTx(ctx context.Context, tx gdb.TX, attemptID int64) error {
	var att entity.ExamAttempt
	if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&att); err != nil {
		return err
	}
	if att.Id == 0 {
		return nil
	}
	row := exam.ExamResult{
		AttemptId:              att.Id,
		ClientUserId:           att.ClientUserId,
		ExamPaperId:            att.ExamPaperId,
		MockExaminationPaperId: att.MockExaminationPaperId,
		Status:          att.Status,
		ObjectiveScore:  att.ObjectiveScore,
		SubjectiveScore: att.SubjectiveScore,
		TotalScore:      att.TotalScore,
		HasSubjective:   att.HasSubjective,
		StartedAt:       att.StartedAt,
		SubmittedAt:     att.SubmittedAt,
		EndedAt:         att.EndedAt,
		CreateTime:      att.CreateTime,
		UpdateTime:      gtime.Now(),
		DeleteFlag:      consts.DeleteFlagNotDeleted,
	}
	var exist entity.ExamResult
	_ = tx.Model(dao.ExamResult.Table()).Ctx(ctx).Where("attempt_id", attemptID).Scan(&exist)
	if exist.AttemptId == 0 {
		_, err := tx.Model(dao.ExamResult.Table()).Ctx(ctx).Insert(row)
		return err
	}
	_, err := tx.Model(dao.ExamResult.Table()).Ctx(ctx).Where("attempt_id", attemptID).Data(row).Update()
	return err
}
