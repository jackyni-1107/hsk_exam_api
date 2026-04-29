package examutil

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
	"exam/internal/dao"
)

func BuildAttemptAnswerDraftRows(attemptID int64, draftMap map[string]string) []g.Map {
	items := make([]g.Map, 0, len(draftMap))
	for _, val := range draftMap {
		itemMap := gconv.Map(val)
		qid := gconv.Int64(itemMap["q"])
		answerJSON := gconv.String(itemMap["a"])
		if qid <= 0 || answerJSON == "" {
			continue
		}
		saveAtTs := gconv.Int64(itemMap["t"])
		saveAt := gtime.Now()
		if saveAtTs > 0 {
			saveAt = gtime.NewFromTimeStamp(saveAtTs)
		}
		if saveAt == nil {
			saveAt = gtime.Now()
		}
		items = append(items, g.Map{
			"attempt_id":       attemptID,
			"exam_question_id": qid,
			"answer_json":      answerJSON,
			"version":          gconv.Int(itemMap["v"]),
			"create_time":      saveAt,
			"update_time":      saveAt,
			"delete_flag":      consts.DeleteFlagNotDeleted,
		})
	}
	return items
}

func UpsertAttemptAnswerDraftRowsTx(ctx context.Context, tx gdb.TX, items []g.Map) error {
	if len(items) == 0 {
		return nil
	}
	_, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
		Data(items).
		Batch(100).
		OnDuplicate(gdb.Raw(`
			answer_json = IF(VALUES(version) >= version, VALUES(answer_json), answer_json),
			update_time = IF(VALUES(version) >= version, VALUES(update_time), update_time),
			version     = IF(VALUES(version) >= version, VALUES(version), version)
		`)).Save()
	return err
}
