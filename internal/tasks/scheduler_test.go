package tasks

import (
	"reflect"
	"testing"

	sysentity "exam/internal/model/entity/sys"
)

func TestBuildCronReloadPlan(t *testing.T) {
	plan := buildCronReloadPlan(
		map[int64]cronEntry{
			1: {name: "task-1", expr: "0 */5 * * * *"},
			2: {name: "task-2", expr: "0 0 * * * *"},
			3: {name: "task-3", expr: "0 15 * * * *"},
		},
		[]sysentity.SysTask{
			{Id: 1, CronExpr: "0 */5 * * * *"},
			{Id: 2, CronExpr: "0 */10 * * * *"},
			{Id: 4, CronExpr: "0 30 * * * *"},
		},
	)

	if !reflect.DeepEqual(plan.removeIDs, []int64{2, 3}) {
		t.Fatalf("unexpected remove ids: %#v", plan.removeIDs)
	}

	gotUpserts := make([]int64, 0, len(plan.upsertTasks))
	for _, task := range plan.upsertTasks {
		gotUpserts = append(gotUpserts, task.Id)
	}
	if !reflect.DeepEqual(gotUpserts, []int64{2, 4}) {
		t.Fatalf("unexpected upsert ids: %#v", gotUpserts)
	}
}
