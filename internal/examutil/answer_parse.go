package examutil

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/model/bo"
)

// ParseAnswerPayload 解析答题 JSON。
func ParseAnswerPayload(s string) bo.AnswerPayload {
	if s == "" {
		return bo.AnswerPayload{}
	}
	j := gjson.New(s)
	var ids []int64
	for _, v := range j.Get("selected_option_ids").Array() {
		ids = append(ids, gconv.Int64(v))
	}
	if len(ids) == 0 {
		if oid := j.Get("option_id"); !oid.IsNil() {
			if id := gconv.Int64(oid); id > 0 {
				ids = []int64{id}
			}
		}
	}
	return bo.AnswerPayload{
		SelectedOptionIDs: ids,
		Text:              j.Get("text").String(),
	}
}
