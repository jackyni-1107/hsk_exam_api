package sysrole

import "testing"

func TestNormalizePositiveIDs(t *testing.T) {
	got := normalizePositiveIDs([]int64{3, 0, 2, 3, -1, 2, 1})
	want := []int64{3, 2, 1}
	if len(got) != len(want) {
		t.Fatalf("unexpected length: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected value at %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestNormalizeRoleStatus(t *testing.T) {
	if got := normalizeRoleStatus(99); got != 0 {
		t.Fatalf("unexpected normalized status: got=%d want=0", got)
	}
	if got := normalizeRoleStatus(1); got != 1 {
		t.Fatalf("unexpected normalized status: got=%d want=1", got)
	}
}

func TestBuildRoleMenuBatch(t *testing.T) {
	batch := buildRoleMenuBatch(12, []int64{5, 8}, "admin")
	if len(batch) != 2 {
		t.Fatalf("unexpected batch length: %d", len(batch))
	}
	if batch[0].RoleId != int64(12) || batch[0].MenuId != int64(5) || batch[0].Creator != "admin" || batch[0].Updater != "admin" {
		t.Fatalf("unexpected first batch row: %+v", batch[0])
	}
	if batch[1].RoleId != int64(12) || batch[1].MenuId != int64(8) || batch[1].Creator != "admin" || batch[1].Updater != "admin" {
		t.Fatalf("unexpected second batch row: %+v", batch[1])
	}
}

func TestNormalizeRoleInput(t *testing.T) {
	name, code, remark := normalizeRoleInput(" admin ", " role:manage ", " memo ")
	if name != "admin" || code != "role:manage" || remark != "memo" {
		t.Fatalf("unexpected normalized input: %q %q %q", name, code, remark)
	}
}
