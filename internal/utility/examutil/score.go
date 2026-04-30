package examutil

import (
	"math"
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
		if ObjectiveAnswerCorrect(q.CorrectOptIDs, ans.OptionID) {
			objective += q.Score
		}
	}
	return objective, paperHasSubjective
}

// ScoreBySegment 按 segment_code 统计客观+主观分，分段分与总分均四舍五入为整数。
func ScoreBySegment(questions []bo.QuestionScoreMeta, answers map[int64]bo.AnswerPayload, awardedByQuestion map[int64]float64) (segmentScores map[string]int, totalScore int, paperHasSubjective bool) {
	paperHasSubjective = PaperHasSubjectiveNonExample(questions)
	rawBySegment := make(map[string]float64)
	for _, q := range questions {
		if q.IsExample != 0 {
			continue
		}
		seg := q.SegmentCode
		if seg == "" {
			seg = "unknown"
		}
		if q.IsSubjective != 0 {
			if awarded, ok := awardedByQuestion[q.QuestionID]; ok {
				rawBySegment[seg] += awarded
			}
			continue
		}
		ans := answers[q.QuestionID]
		if ObjectiveAnswerCorrect(q.CorrectOptIDs, ans.OptionID) {
			rawBySegment[seg] += q.Score
		}
	}
	segmentScores = make(map[string]int, len(rawBySegment))
	for seg, raw := range rawBySegment {
		rounded := int(math.Round(raw))
		segmentScores[seg] = rounded
		totalScore += rounded
	}
	return segmentScores, totalScore, paperHasSubjective
}

// ObjectiveAnswerCorrect 单选客观题是否选对：正确选项集合唯一且与用户选项一致。
func ObjectiveAnswerCorrect(correctIDs []int64, optionID int64) bool {
	if optionID <= 0 || len(correctIDs) != 1 {
		return false
	}
	return correctIDs[0] == optionID
}

// EmptyAnswerRowsForPaper 根据试卷全部小题 ID 生成「空答题行」描述（供单测与客户端初始化占位）。
func EmptyAnswerRowsForPaper(questionIDs []int64) []int64 {
	out := append([]int64(nil), questionIDs...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
