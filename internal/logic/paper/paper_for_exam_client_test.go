package paper

import (
	"encoding/json"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
	"testing"
)

func TestStripSensitiveExamFields_KeepExampleAnswer(t *testing.T) {
	topic := &exambo.SectionTopic{
		Items: []exambo.TopicItem{
			{
				IsExample: json.RawMessage(`true`),
				Questions: []exambo.TopicQuestion{
					{
						Answers: []exambo.TopicAnswer{
							{
								Extra: map[string]json.RawMessage{
									"answer":         json.RawMessage(`"A"`),
									"correct":        json.RawMessage(`true`),
									"correct_answer": json.RawMessage(`"A"`),
									"score":          json.RawMessage(`2`),
								},
							},
						},
						Extra: map[string]json.RawMessage{
							"question_score": json.RawMessage(`3`),
						},
					},
				},
			},
			{
				Answers: []exambo.TopicAnswer{
					{
						Extra: map[string]json.RawMessage{
							"answer":         json.RawMessage(`"B"`),
							"correct_answer": json.RawMessage(`"B"`),
							"score_total":    json.RawMessage(`10`),
						},
					},
				},
			},
		},
	}

	stripSensitiveExamFields(topic)

	exampleExtra := topic.Items[0].Questions[0].Answers[0].Extra
	if _, ok := exampleExtra["answer"]; !ok {
		t.Fatalf("example answer should be kept")
	}
	if _, ok := exampleExtra["correct_answer"]; !ok {
		t.Fatalf("example correct_answer should be kept")
	}
	if _, ok := exampleExtra["score"]; ok {
		t.Fatalf("score must be removed")
	}

	normalExtra := topic.Items[1].Answers[0].Extra
	if _, ok := normalExtra["answer"]; ok {
		t.Fatalf("non-example answer must be removed")
	}
	if _, ok := normalExtra["correct_answer"]; ok {
		t.Fatalf("non-example correct_answer must be removed")
	}
	if _, ok := normalExtra["score_total"]; ok {
		t.Fatalf("score_total must be removed")
	}
}

func TestEnrichAnswersWithExamIDs(t *testing.T) {
	options := []examentity.ExamOption{
		{Id: 101, SortOrder: 0, Flag: "A"},
		{Id: 102, SortOrder: 1, Flag: "B"},
		{Id: 103, SortOrder: 2, Flag: "C"},
	}
	answers := []exambo.TopicAnswer{
		{Flag: " b "},
		{Id: json.RawMessage(`"103"`)},
		{Index: json.RawMessage(`0`)},
	}

	enrichAnswersWithExamIDs(answers, options)

	if answers[0].EAID == nil || *answers[0].EAID != 102 {
		t.Fatalf("flag match failed, got %+v", answers[0].EAID)
	}
	if answers[1].EAID == nil || *answers[1].EAID != 103 {
		t.Fatalf("id match failed, got %+v", answers[1].EAID)
	}
	if answers[2].EAID == nil || *answers[2].EAID != 101 {
		t.Fatalf("index match failed, got %+v", answers[2].EAID)
	}
}

func TestTopicHasStaleEaid(t *testing.T) {
	okID := int64(1)
	topic := &exambo.SectionTopic{
		Items: []exambo.TopicItem{
			{
				Questions: []exambo.TopicQuestion{
					{
						Answers: []exambo.TopicAnswer{{EAID: &okID}, {}},
					},
				},
			},
		},
	}
	if !topicHasStaleEaid(topic) {
		t.Fatalf("expected stale eaid to be detected")
	}
	topic.Items[0].Questions[0].Answers[1].EAID = &okID
	if topicHasStaleEaid(topic) {
		t.Fatalf("did not expect stale eaid")
	}
}

func TestStripYCTItemRenderFields(t *testing.T) {
	topic := &exambo.SectionTopic{
		Items: []exambo.TopicItem{
			{
				Extra: map[string]json.RawMessage{
					"_converter": json.RawMessage(`{}`),
					"_element":   json.RawMessage(`{}`),
					"keep":       json.RawMessage(`1`),
				},
			},
		},
	}
	stripYCTItemRenderFields(topic)
	if _, ok := topic.Items[0].Extra["_converter"]; ok {
		t.Fatalf("_converter should be removed")
	}
	if _, ok := topic.Items[0].Extra["_element"]; ok {
		t.Fatalf("_element should be removed")
	}
	if _, ok := topic.Items[0].Extra["keep"]; !ok {
		t.Fatalf("keep should remain")
	}
}
