package batch

import (
	"testing"

	"exam/internal/consts"
	"exam/internal/model/bo"
)

func TestNormalizeExamBatchPolicy(t *testing.T) {
	p, err := normalizeExamBatchPolicy(bo.ExamBatchPolicyInput{
		BatchKind:             consts.ExamBatchKindFormal,
		AllowMultipleAttempts: 2,
		MaxAttemptsPerMember:  5,
		SkipScoring:           1,
		AutoSubmitOnDeadline:  0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.AllowMultipleAttempts != 1 {
		t.Fatalf("allow multiple clamped to 1, got %d", p.AllowMultipleAttempts)
	}
	if p.MaxAttemptsPerMember != 5 {
		t.Fatalf("max attempts: %d", p.MaxAttemptsPerMember)
	}
	if p.SkipScoring != 1 || p.AutoSubmitOnDeadline != 0 {
		t.Fatalf("flags skip=%d auto=%d", p.SkipScoring, p.AutoSubmitOnDeadline)
	}

	p2, err := normalizeExamBatchPolicy(bo.ExamBatchPolicyInput{
		BatchKind:             consts.ExamBatchKindFormal,
		AllowMultipleAttempts: 0,
		MaxAttemptsPerMember:  9,
		SkipScoring:           0,
		AutoSubmitOnDeadline:  1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if p2.MaxAttemptsPerMember != 0 {
		t.Fatalf("when not allow multiple, max should reset to 0, got %d", p2.MaxAttemptsPerMember)
	}
}
