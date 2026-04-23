package attempt

import (
	"testing"

	"github.com/gogf/gf/v2/util/gconv"
)

func TestBuildAttemptAnswerDraftRows(t *testing.T) {
	rows := buildAttemptAnswerDraftRows(10, map[string]string{
		"1":   `{"q":1,"a":"{\"option_id\":100}","v":3,"t":1776926400}`,
		"2":   `{"q":2,"a":"","v":1,"t":1776926401}`,
		"bad": `{"q":0,"a":"{}","v":1,"t":1776926402}`,
	}, updaterClient)

	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	row := rows[0]
	if got := gconv.Int64(row["attempt_id"]); got != 10 {
		t.Fatalf("attempt_id=%d", got)
	}
	if got := gconv.Int64(row["exam_question_id"]); got != 1 {
		t.Fatalf("exam_question_id=%d", got)
	}
	if got := gconv.Int(row["version"]); got != 3 {
		t.Fatalf("version=%d", got)
	}
	if got := gconv.String(row["answer_json"]); got == "" {
		t.Fatal("answer_json should not be empty")
	}
	if got := gconv.String(row["updater"]); got != updaterClient {
		t.Fatalf("updater=%q", got)
	}
}
