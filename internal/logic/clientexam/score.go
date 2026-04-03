package clientexam

import (
	"sort"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
)

// QuestionScoreMeta 阅卷用题目元数据（由 DB 加载后填充）。
type QuestionScoreMeta struct {
	QuestionID    int64
	IsExample     int
	IsSubjective  int
	Score         float64
	CorrectOptIDs []int64
}

// AnswerPayload 客户端写入的 answer_json 结构。
type AnswerPayload struct {
	SelectedOptionIDs []int64 `json:"selected_option_ids"`
	Text              string  `json:"text"`
}

// ParseAnswerPayload 解析答题 JSON。
func ParseAnswerPayload(s string) AnswerPayload {
	if s == "" {
		return AnswerPayload{}
	}
	j := gjson.New(s)
	var ids []int64
	for _, v := range j.Get("selected_option_ids").Array() {
		ids = append(ids, gconv.Int64(v))
	}
	return AnswerPayload{
		SelectedOptionIDs: ids,
		Text:              j.Get("text").String(),
	}
}

// PaperHasSubjectiveNonExample 试卷是否含需人工的主观题（非例题）。
func PaperHasSubjectiveNonExample(questions []QuestionScoreMeta) bool {
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
func ScoreObjective(questions []QuestionScoreMeta, answers map[int64]AnswerPayload) (objective float64, paperHasSubjective bool) {
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
