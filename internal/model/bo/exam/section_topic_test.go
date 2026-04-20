package exam

import (
	"encoding/json"
	"testing"
)

func TestSectionTopic_UnknownFieldRoundTrip(t *testing.T) {
	raw := []byte(`{
		"unknown_root": 1,
		"items": [{
			"unknown_item": {"x":1},
			"questions": [{
				"unknown_question": "q",
				"answers": [{
					"flag":"A",
					"unknown_answer": true
				}]
			}]
		}]
	}`)

	var topic SectionTopic
	if err := json.Unmarshal(raw, &topic); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if _, ok := topic.Extra["unknown_root"]; !ok {
		t.Fatalf("missing unknown root field")
	}
	if _, ok := topic.Items[0].Extra["unknown_item"]; !ok {
		t.Fatalf("missing unknown item field")
	}
	if _, ok := topic.Items[0].Questions[0].Extra["unknown_question"]; !ok {
		t.Fatalf("missing unknown question field")
	}
	if _, ok := topic.Items[0].Questions[0].Answers[0].Extra["unknown_answer"]; !ok {
		t.Fatalf("missing unknown answer field")
	}

	b, err := json.Marshal(&topic)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal output failed: %v", err)
	}
	if _, ok := out["unknown_root"]; !ok {
		t.Fatalf("unknown root field should be preserved")
	}
}
