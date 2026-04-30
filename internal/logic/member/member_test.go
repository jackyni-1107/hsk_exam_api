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
	pwd, err := memberImportPasswordFromEmail("demo@example.com", "hskmock", []int{1, 3, 5})
	if err != nil {
		t.Fatal(err)
	}
	if want := "dm@@hskmock"; pwd != want {
		t.Fatalf("got %q want %q", pwd, want)
	}
	if _, err := memberImportPasswordFromEmail("a@b", "hskmock", []int{1, 3, 5}); err == nil {
		t.Fatal("expected error for short email")
	}
}

func TestMemberImportParsePickPositions(t *testing.T) {
	got, err := memberImportParsePickPositions("1, 3,5")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 3 || got[2] != 5 {
		t.Fatalf("unexpected parsed positions: %#v", got)
	}
	got, err = memberImportParsePickPositions("1，2；4")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 4 {
		t.Fatalf("unexpected parsed positions with cn punctuation: %#v", got)
	}
	if _, err := memberImportParsePickPositions("0,2"); err == nil {
		t.Fatal("expected invalid error for zero index")
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

func TestMemberImportParseUsernameSeq_digitWidth(t *testing.T) {
	r5 := &memberImportUsernameRule{country: "TH", year: "2026", digits: 5}
	if seq, ok := memberImportParseUsernameSeq("TH2026-00004", r5); !ok || seq != 4 {
		t.Fatalf("5-digit rule TH2026-00004: ok=%v seq=%d", ok, seq)
	}
	if _, ok := memberImportParseUsernameSeq("TH2026-000004", r5); ok {
		t.Fatal("6-char numeric suffix should not match 5-digit rule")
	}

	r6 := &memberImportUsernameRule{country: "TH", year: "2026", digits: 6}
	if seq, ok := memberImportParseUsernameSeq("TH2026-000004", r6); !ok || seq != 4 {
		t.Fatalf("6-digit rule TH2026-000004: ok=%v seq=%d", ok, seq)
	}
	if _, ok := memberImportParseUsernameSeq("TH2026-00004", r6); ok {
		t.Fatal("5-char numeric suffix should not match 6-digit rule")
	}
}
