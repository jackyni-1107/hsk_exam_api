package member

import (
	"testing"

	"exam/internal/consts"
)

func TestNormalizeMemberUsername(t *testing.T) {
	if got := normalizeMemberUsername("  Admin  "); got != "Admin" {
		t.Fatalf("unexpected normalized username: %q", got)
	}
}

func TestNormalizeMemberProfile(t *testing.T) {
	nickname, email, mobile := normalizeMemberProfile("  Nick  ", "  a@example.com  ", "  13800138000  ")
	if nickname != "Nick" || email != "a@example.com" || mobile != "13800138000" {
		t.Fatalf("unexpected normalized profile: %q %q %q", nickname, email, mobile)
	}
}

func TestNormalizeMemberStatus(t *testing.T) {
	if got := normalizeMemberStatus(consts.StatusDisabled); got != consts.StatusDisabled {
		t.Fatalf("expected disabled status, got %d", got)
	}
	if got := normalizeMemberStatus(123); got != consts.StatusNormal {
		t.Fatalf("expected default normal status, got %d", got)
	}
}

func TestMemberImportPasswordFromEmail(t *testing.T) {
	pwd, err := memberImportPasswordFromEmail("demo@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if want := "dm@@hskmock"; pwd != want {
		t.Fatalf("got %q want %q", pwd, want)
	}
	if _, err := memberImportPasswordFromEmail("a@b"); err == nil {
		t.Fatal("expected error for short email")
	}
}

func TestMemberImportCol(t *testing.T) {
	idx := map[string]int{"password": 2}
	if got := memberImportCol(idx, "password"); got != 2 {
		t.Fatalf("got %d", got)
	}
	if got := memberImportCol(idx, "nickname"); got != -1 {
		t.Fatalf("missing key should be -1, got %d", got)
	}
}
