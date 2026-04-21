package sysuser

import "testing"

func TestNormalizeUsername(t *testing.T) {
	if got := normalizeUsername("  Admin.User  "); got != "Admin.User" {
		t.Fatalf("normalizeUsername() = %q, want %q", got, "Admin.User")
	}
}

func TestNormalizeRoleIDs(t *testing.T) {
	got := normalizePositiveIDs([]int64{3, 0, 2, 3, -1, 2, 5})
	want := []int64{3, 2, 5}
	if len(got) != len(want) {
		t.Fatalf("len(normalizePositiveIDs()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("normalizePositiveIDs()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestNormalizeUserStatus(t *testing.T) {
	if got := normalizeUserStatus(999); got != 0 {
		t.Fatalf("normalizeUserStatus() = %d, want %d", got, 0)
	}
	if got := normalizeUserStatus(1); got != 1 {
		t.Fatalf("normalizeUserStatus() = %d, want %d", got, 1)
	}
}

func TestBuildUserRoleBatch(t *testing.T) {
	batch := buildUserRoleBatch(7, []int64{2, 5}, "admin")
	if len(batch) != 2 {
		t.Fatalf("len(buildUserRoleBatch()) = %d, want %d", len(batch), 2)
	}
	if batch[0].UserId != int64(7) || batch[0].RoleId != int64(2) || batch[0].Creator != "admin" {
		t.Fatalf("unexpected first batch item: %+v", batch[0])
	}
	if batch[1].UserId != int64(7) || batch[1].RoleId != int64(5) || batch[1].Updater != "admin" {
		t.Fatalf("unexpected second batch item: %+v", batch[1])
	}
}
