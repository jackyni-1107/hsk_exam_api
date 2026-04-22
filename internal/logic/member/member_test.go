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
