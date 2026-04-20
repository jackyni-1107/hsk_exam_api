package examutil

import (
	"testing"

	"exam/internal/model/bo"
)

func TestScoreObjective_AllCorrect(t *testing.T) {
	qs := []bo.QuestionScoreMeta{
		{QuestionID: 1, IsExample: 0, IsSubjective: 0, Score: 5, CorrectOptIDs: []int64{10}},
		{QuestionID: 2, IsExample: 0, IsSubjective: 0, Score: 3, CorrectOptIDs: []int64{20}},
	}
	ans := map[int64]bo.AnswerPayload{
		1: {OptionID: 10},
		2: {OptionID: 20},
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
		1: {OptionID: 999},
		2: {OptionID: 2},
		3: {OptionID: 3},
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
	p := ParseAnswerPayload(`{"o_id":42}`)
	if p.OptionID != 42 || p.Text != "" {
		t.Fatalf("%+v", p)
	}
	p2 := ParseAnswerPayload(`{"text":"x"}`)
	if p2.OptionID != 0 || p2.Text != "x" {
		t.Fatalf("%+v", p2)
	}
	p3 := ParseAnswerPayload("")
	if p3.OptionID != 0 || p3.Text != "" {
		t.Fatalf("%+v", p3)
	}
}

func TestMarshalAnswerPayload(t *testing.T) {
	if got := MarshalAnswerPayload(bo.AnswerPayload{OptionID: 559}); got != `{"o_id":559}` {
		t.Fatalf("option got %s", got)
	}
	if got := MarshalAnswerPayload(bo.AnswerPayload{Text: "hello"}); got != `{"text":"hello"}` {
		t.Fatalf("text got %s", got)
	}
	if got := MarshalAnswerPayload(bo.AnswerPayload{}); got != "" {
		t.Fatalf("empty got %s", got)
	}
}

func TestObjectiveAnswerCorrect(t *testing.T) {
	if !ObjectiveAnswerCorrect([]int64{10}, 10) {
		t.Fatal("should be correct")
	}
	if ObjectiveAnswerCorrect([]int64{10}, 11) {
		t.Fatal("should be wrong")
	}
	if ObjectiveAnswerCorrect([]int64{10}, 0) {
		t.Fatal("zero should be wrong")
	}
	if ObjectiveAnswerCorrect(nil, 10) {
		t.Fatal("empty correct ids should be wrong")
	}
}
