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

func BuildAttemptAnswerDraftRows(attemptID int64, draftMap map[string]string, updater string) []g.Map {
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
			"creator":          updater,
			"create_time":      saveAt,
			"updater":          updater,
			"update_time":      saveAt,
			"delete_flag":      consts.DeleteFlagNotDeleted,
		})
	}
	return items
}

func UpsertAttemptAnswerDraftRowsTx(ctx context.Context, tx gdb.TX, items []g.Map) error {
	for _, item := range items {
		attemptID := gconv.Int64(item["attempt_id"])
		questionID := gconv.Int64(item["exam_question_id"])
		version := gconv.Int(item["version"])

		r, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
			Where("attempt_id", attemptID).
			Where("exam_question_id", questionID).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Where("version <= ?", version).
			Update(g.Map{
				"answer_json": item["answer_json"],
				"version":     item["version"],
				"updater":     item["updater"],
				"update_time": item["update_time"],
			})
		if err != nil {
			return err
		}
		if n, _ := r.RowsAffected(); n > 0 {
			continue
		}

		cnt, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
			Where("attempt_id", attemptID).
			Where("exam_question_id", questionID).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			continue
		}
		if _, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Insert(item); err != nil {
			return err
		}
	}
	return nil
}
