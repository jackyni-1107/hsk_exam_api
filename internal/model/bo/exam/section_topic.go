package exam

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// SectionTopic 表示 section.topic_json 的强类型结构。
// 未建模字段会落在 Extra，确保历史字段不丢失。
type SectionTopic struct {
	Items     []TopicItem                `json:"items,omitempty"`
	IsExample json.RawMessage            `json:"is_example,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

type TopicItem struct {
	Questions []TopicQuestion            `json:"questions,omitempty"`
	Answers   []TopicAnswer              `json:"answers,omitempty"`
	Eqid      *int64                     `json:"eqid,omitempty"`
	IsExample json.RawMessage            `json:"is_example,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

type TopicQuestion struct {
	Answers   []TopicAnswer              `json:"answers,omitempty"`
	Eqid      *int64                     `json:"eqid,omitempty"`
	IsExample json.RawMessage            `json:"is_example,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

type TopicAnswer struct {
	Flag  string                     `json:"flag,omitempty"`
	Id    json.RawMessage            `json:"id,omitempty"`
	Index json.RawMessage            `json:"index,omitempty"`
	EAID  *int64                     `json:"eaid,omitempty"`
	Extra map[string]json.RawMessage `json:"-"`
}

func (s SectionTopic) MarshalJSON() ([]byte, error) {
	m := cloneRawMap(s.Extra)
	putRaw(m, "items", s.Items)
	putRawMessage(m, "is_example", s.IsExample)
	return json.Marshal(m)
}

func (s *SectionTopic) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["items"]; ok {
		if err := json.Unmarshal(v, &s.Items); err != nil {
			return err
		}
		delete(raw, "items")
	}
	if v, ok := raw["is_example"]; ok {
		s.IsExample = cloneRaw(v)
		delete(raw, "is_example")
	}
	s.Extra = raw
	return nil
}

func (s TopicItem) MarshalJSON() ([]byte, error) {
	m := cloneRawMap(s.Extra)
	putRaw(m, "questions", s.Questions)
	putRaw(m, "answers", s.Answers)
	if s.Eqid != nil {
		putRaw(m, "eqid", *s.Eqid)
	}
	putRawMessage(m, "is_example", s.IsExample)
	return json.Marshal(m)
}

func (s *TopicItem) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["questions"]; ok {
		if err := json.Unmarshal(v, &s.Questions); err != nil {
			return err
		}
		delete(raw, "questions")
	}
	if v, ok := raw["answers"]; ok {
		if err := json.Unmarshal(v, &s.Answers); err != nil {
			return err
		}
		delete(raw, "answers")
	}
	if v, ok := raw["eqid"]; ok {
		var eqid int64
		if err := json.Unmarshal(v, &eqid); err == nil {
			s.Eqid = &eqid
		}
		delete(raw, "eqid")
	}
	if v, ok := raw["is_example"]; ok {
		s.IsExample = cloneRaw(v)
		delete(raw, "is_example")
	}
	s.Extra = raw
	return nil
}

func (s TopicQuestion) MarshalJSON() ([]byte, error) {
	m := cloneRawMap(s.Extra)
	putRaw(m, "answers", s.Answers)
	if s.Eqid != nil {
		putRaw(m, "eqid", *s.Eqid)
	}
	putRawMessage(m, "is_example", s.IsExample)
	return json.Marshal(m)
}

func (s *TopicQuestion) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["answers"]; ok {
		if err := json.Unmarshal(v, &s.Answers); err != nil {
			return err
		}
		delete(raw, "answers")
	}
	if v, ok := raw["eqid"]; ok {
		var eqid int64
		if err := json.Unmarshal(v, &eqid); err == nil {
			s.Eqid = &eqid
		}
		delete(raw, "eqid")
	}
	if v, ok := raw["is_example"]; ok {
		s.IsExample = cloneRaw(v)
		delete(raw, "is_example")
	}
	s.Extra = raw
	return nil
}

func (s TopicAnswer) MarshalJSON() ([]byte, error) {
	m := cloneRawMap(s.Extra)
	if s.Flag != "" {
		putRaw(m, "flag", s.Flag)
	}
	putRawMessage(m, "id", s.Id)
	putRawMessage(m, "index", s.Index)
	if s.EAID != nil {
		putRaw(m, "eaid", *s.EAID)
	}
	return json.Marshal(m)
}

func (s *TopicAnswer) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["flag"]; ok {
		_ = json.Unmarshal(v, &s.Flag)
		delete(raw, "flag")
	}
	if v, ok := raw["id"]; ok {
		s.Id = cloneRaw(v)
		delete(raw, "id")
	}
	if v, ok := raw["index"]; ok {
		s.Index = cloneRaw(v)
		delete(raw, "index")
	}
	if v, ok := raw["eaid"]; ok {
		var eaid int64
		if err := json.Unmarshal(v, &eaid); err == nil {
			s.EAID = &eaid
		}
		delete(raw, "eaid")
	}
	s.Extra = raw
	return nil
}

func RawTruthy(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return false
	}
	var bv bool
	if err := json.Unmarshal(trimmed, &bv); err == nil {
		return bv
	}
	var iv int64
	if err := json.Unmarshal(trimmed, &iv); err == nil {
		return iv != 0
	}
	var fv float64
	if err := json.Unmarshal(trimmed, &fv); err == nil {
		return fv != 0
	}
	var sv string
	if err := json.Unmarshal(trimmed, &sv); err == nil {
		sv = strings.TrimSpace(strings.ToLower(sv))
		return sv == "1" || sv == "true"
	}
	return false
}

func RawString(raw json.RawMessage) (string, bool) {
	if len(raw) == 0 {
		return "", false
	}
	var sv string
	if err := json.Unmarshal(raw, &sv); err != nil {
		return "", false
	}
	return sv, true
}

func RawInt(raw json.RawMessage) (int, bool) {
	if len(raw) == 0 {
		return 0, false
	}
	var iv int
	if err := json.Unmarshal(raw, &iv); err == nil {
		return iv, true
	}
	var fv float64
	if err := json.Unmarshal(raw, &fv); err == nil {
		return int(fv), true
	}
	var sv string
	if err := json.Unmarshal(raw, &sv); err == nil {
		n, convErr := strconv.Atoi(strings.TrimSpace(sv))
		if convErr == nil {
			return n, true
		}
	}
	return 0, false
}

func cloneRaw(v json.RawMessage) json.RawMessage {
	if len(v) == 0 {
		return nil
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out
}

func cloneRawMap(in map[string]json.RawMessage) map[string]json.RawMessage {
	if len(in) == 0 {
		return make(map[string]json.RawMessage)
	}
	out := make(map[string]json.RawMessage, len(in))
	for k, v := range in {
		out[k] = cloneRaw(v)
	}
	return out
}

func putRaw(m map[string]json.RawMessage, key string, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil || string(b) == "null" {
		return
	}
	m[key] = b
}

func putRawMessage(m map[string]json.RawMessage, key string, value json.RawMessage) {
	if len(value) == 0 {
		return
	}
	m[key] = cloneRaw(value)
}
