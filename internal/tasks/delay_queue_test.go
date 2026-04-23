package tasks

import (
	"testing"

	"exam/internal/consts"
)

func TestEncodeDecodeDelayedExecRoundTrip(t *testing.T) {
	original := ExecRequest{
		TaskID:      11,
		RunID:       "run-11",
		TriggerType: consts.TriggerTypeDelay,
		RetryCount:  2,
		Params:      `{"k":"v"}`,
	}
	member, err := encodeDelayedExec(original)
	if err != nil {
		t.Fatalf("encode delayed exec failed: %v", err)
	}
	got, err := decodeDelayedExec(member)
	if err != nil {
		t.Fatalf("decode delayed exec failed: %v", err)
	}
	if got != original {
		t.Fatalf("unexpected delayed exec round trip: %#v", got)
	}
}
