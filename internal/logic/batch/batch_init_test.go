package batch

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	s := New()

	got, err := s.parseTime("")
	if err != nil {
		t.Fatalf("blank time should not error: %v", err)
	}
	if got != nil {
		t.Fatalf("blank time should return nil, got %v", got)
	}

	got, err = s.parseTime("2024-01-02T03:04:05+08:00")
	if err != nil {
		t.Fatalf("expected valid time, got error: %v", err)
	}
	wantWithOffset := time.Date(2024, 1, 1, 19, 4, 5, 0, time.UTC).Unix()
	if got == nil || got.Timestamp() != wantWithOffset {
		t.Fatalf("unexpected UTC timestamp for offset input: got=%v want=%d", got, wantWithOffset)
	}

	got, err = s.parseTime("2024-01-02T03:04:05Z")
	if err != nil {
		t.Fatalf("expected valid RFC3339, got error: %v", err)
	}
	wantUTC := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC).Unix()
	if got == nil || got.Timestamp() != wantUTC {
		t.Fatalf("unexpected UTC timestamp: got=%v want=%d", got, wantUTC)
	}

	if _, err = s.parseTime("2024-01-02 03:04:05"); err == nil {
		t.Fatalf("time without timezone should fail")
	}

	if _, err = s.parseTime("not-a-time"); err == nil {
		t.Fatalf("invalid time should fail")
	}
}
