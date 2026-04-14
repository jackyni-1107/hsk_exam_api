package paper

import (
	"context"
	"sort"

	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
)

// loadExamOptionsGrouped 按 question_id 批量加载选项，每组内按 sort_order 排序。
func loadExamOptionsGrouped(ctx context.Context, questionIDs []int64) (map[int64][]examentity.ExamOption, error) {
	out := make(map[int64][]examentity.ExamOption)
	if len(questionIDs) == 0 {
		return out, nil
	}
	ids := make([]interface{}, len(questionIDs))
	for i, id := range questionIDs {
		ids[i] = id
	}
	var all []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		WhereIn("question_id", ids).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&all); err != nil {
		return nil, err
	}
	for _, o := range all {
		qid := o.QuestionId
		out[qid] = append(out[qid], o)
	}
	for qid := range out {
		opts := out[qid]
		sort.Slice(opts, func(i, j int) bool { return opts[i].SortOrder < opts[j].SortOrder })
	}
	return out, nil
}
