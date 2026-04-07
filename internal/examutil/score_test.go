package examutil

import (
	"testing"

	"exam/internal/model/bo"
)

func TestScoreObjective_AllCorrect(t *testing.T) {
	qs := []bo.QuestionScoreMeta{
		{QuestionID: 1, IsExample: 0, IsSubjective: 0, Score: 5, CorrectOptIDs: []int64{10}},
		{QuestionID: 2, IsExample: 0, IsSubjective: 0, Score: 3, CorrectOptIDs: []int64{20, 21}},
	}
	ans := map[int64]bo.AnswerPayload{
		1: {SelectedOptionIDs: []int64{10}},
		2: {SelectedOptionIDs: []int64{21, 20}},
	}
	obj, hasSubj := ScoreObjective(qs, ans)
	if hasSubj {
		t.Fatal("hasSubjective")
	}
	if obj != 8 {
		t.Fatalf("objective want 8 got %v", obj)
	}
}

func TestScoreObjective_SkipsExampleAndSubjective(t *testing.T) {
	qs := []bo.QuestionScoreMeta{
		{QuestionID: 1, IsExample: 1, IsSubjective: 0, Score: 99, CorrectOptIDs: []int64{1}},
		{QuestionID: 2, IsExample: 0, IsSubjective: 1, Score: 50, CorrectOptIDs: []int64{2}},
		{QuestionID: 3, IsExample: 0, IsSubjective: 0, Score: 2, CorrectOptIDs: []int64{3}},
	}
	ans := map[int64]bo.AnswerPayload{
		1: {SelectedOptionIDs: []int64{999}},
		2: {SelectedOptionIDs: []int64{2}},
		3: {SelectedOptionIDs: []int64{3}},
	}
	obj, hasSubj := ScoreObjective(qs, ans)
	if !hasSubj {
		t.Fatal("want paper has subjective")
	}
	if obj != 2 {
		t.Fatalf("objective want 2 got %v", obj)
	}
}

func TestPaperHasSubjectiveNonExample(t *testing.T) {
	if !PaperHasSubjectiveNonExample([]bo.QuestionScoreMeta{{IsExample: 0, IsSubjective: 1}}) {
		t.Fatal()
	}
	if PaperHasSubjectiveNonExample([]bo.QuestionScoreMeta{{IsExample: 1, IsSubjective: 1}}) {
		t.Fatal()
	}
}

func TestEmptyAnswerRowsForPaper(t *testing.T) {
	got := EmptyAnswerRowsForPaper([]int64{3, 1, 2})
	want := []int64{1, 2, 3}
	if len(got) != len(want) {
		t.Fatal()
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestParseAnswerPayload(t *testing.T) {
	p := ParseAnswerPayload(`{"selected_option_ids":[2,1],"text":"x"}`)
	if len(p.SelectedOptionIDs) != 2 || p.Text != "x" {
		t.Fatalf("%+v", p)
	}
}
