package examutil

import (
	"encoding/json"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/model/bo"
)

// ParseAnswerPayload 解析答题 JSON。存储格式：{"o_id":X} 或 {"text":"..."}；空串视为未作答。
func ParseAnswerPayload(s string) bo.AnswerPayload {
	if s == "" {
		return bo.AnswerPayload{}
	}
	j := gjson.New(s)
	return bo.AnswerPayload{
		OptionID: gconv.Int64(j.Get("o_id")),
		Text:     j.Get("text").String(),
	}
}

// MarshalAnswerPayload 将 AnswerPayload 序列化为存储 JSON；两字段均零值时返回空串（表示未作答）。
func MarshalAnswerPayload(p bo.AnswerPayload) string {
	if p.OptionID == 0 && p.Text == "" {
		return ""
	}
	raw, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(raw)
}
