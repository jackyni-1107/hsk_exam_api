package examutil

import (
	"sort"

	"exam/internal/model/bo"
)

// PaperHasSubjectiveNonExample 试卷是否含需人工的主观题（非例题）。
func PaperHasSubjectiveNonExample(questions []bo.QuestionScoreMeta) bool {
	for _, q := range questions {
		if q.IsExample != 0 {
			continue
		}
		if q.IsSubjective != 0 {
			return true
		}
	}
	return false
}

// ScoreObjective 仅客观题自动分；例题与主观题不计分。返回答案侧客观题得分与试卷是否含主观题。
func ScoreObjective(questions []bo.QuestionScoreMeta, answers map[int64]bo.AnswerPayload) (objective float64, paperHasSubjective bool) {
	paperHasSubjective = PaperHasSubjectiveNonExample(questions)
	for _, q := range questions {
		if q.IsExample != 0 || q.IsSubjective != 0 {
			continue
		}
		ans := answers[q.QuestionID]
		if objectiveQuestionCorrect(q.CorrectOptIDs, ans.SelectedOptionIDs) {
			objective += q.Score
		}
	}
	return objective, paperHasSubjective
}

// ObjectiveAnswerCorrect 客观题是否选对（多选需与正确选项 id 集合完全一致，顺序无关）。
func ObjectiveAnswerCorrect(correctIDs, selected []int64) bool {
	return objectiveQuestionCorrect(correctIDs, selected)
}

func objectiveQuestionCorrect(correctIDs, selected []int64) bool {
	if len(correctIDs) == 0 {
		return false
	}
	a := append([]int64(nil), correctIDs...)
	b := append([]int64(nil), selected...)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// EmptyAnswerRowsForPaper 根据试卷全部小题 ID 生成「空答题行」描述（供单测与客户端初始化占位）。
func EmptyAnswerRowsForPaper(questionIDs []int64) []int64 {
	out := append([]int64(nil), questionIDs...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
