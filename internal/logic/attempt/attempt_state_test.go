package attempt

import (
	"testing"

	"exam/internal/consts"
	examentity "exam/internal/model/entity/exam"

	"github.com/gogf/gf/v2/os/gtime"
)

func TestAttemptStateTransitions(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		event      attemptEvent
		wantStatus int
		wantOK     bool
	}{
		{
			name:       "start not started",
			status:     consts.ExamAttemptNotStarted,
			event:      attemptEventStart,
			wantStatus: consts.ExamAttemptInProgress,
			wantOK:     true,
		},
		{
			name:       "save answers keeps in progress",
			status:     consts.ExamAttemptInProgress,
			event:      attemptEventSaveAnswers,
			wantStatus: consts.ExamAttemptInProgress,
			wantOK:     true,
		},
		{
			name:       "submit in progress",
			status:     consts.ExamAttemptInProgress,
			event:      attemptEventSubmit,
			wantStatus: consts.ExamAttemptSubmitted,
			wantOK:     true,
		},
		{
			name:       "timeout submits in progress",
			status:     consts.ExamAttemptInProgress,
			event:      attemptEventTimeout,
			wantStatus: consts.ExamAttemptSubmitted,
			wantOK:     true,
		},
		{
			name:       "batch expired submits in progress",
			status:     consts.ExamAttemptInProgress,
			event:      attemptEventBatchExpired,
			wantStatus: consts.ExamAttemptSubmitted,
			wantOK:     true,
		},
		{
			name:       "finalize submitted",
			status:     consts.ExamAttemptSubmitted,
			event:      attemptEventFinalize,
			wantStatus: consts.ExamAttemptEnded,
			wantOK:     true,
		},
		{
			name:       "grade scored attempt",
			status:     consts.ExamAttemptEnded,
			event:      attemptEventGrade,
			wantStatus: consts.ExamAttemptEnded,
			wantOK:     true,
		},
		{
			name:       "cannot submit not started",
			status:     consts.ExamAttemptNotStarted,
			event:      attemptEventSubmit,
			wantStatus: consts.ExamAttemptNotStarted,
			wantOK:     false,
		},
		{
			name:       "cannot save submitted",
			status:     consts.ExamAttemptSubmitted,
			event:      attemptEventSaveAnswers,
			wantStatus: consts.ExamAttemptSubmitted,
			wantOK:     false,
		},
		{
			name:       "cannot finalize scored twice",
			status:     consts.ExamAttemptEnded,
			event:      attemptEventFinalize,
			wantStatus: consts.ExamAttemptEnded,
			wantOK:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotStatus, gotOK := transitionAttemptStatus(tc.status, tc.event)
			if gotOK != tc.wantOK {
				t.Fatalf("ok=%v, want %v", gotOK, tc.wantOK)
			}
			if gotStatus != tc.wantStatus {
				t.Fatalf("status=%d, want %d", gotStatus, tc.wantStatus)
			}
		})
	}
}

func TestAttemptStateConvenienceGuardsUseTransitions(t *testing.T) {
	if !canStartAttempt(consts.ExamAttemptNotStarted) {
		t.Fatal("not started attempt should be startable")
	}
	if canStartAttempt(consts.ExamAttemptInProgress) {
		t.Fatal("in progress attempt should not be startable")
	}
	if !canSaveAttemptAnswers(consts.ExamAttemptInProgress) {
		t.Fatal("in progress attempt should allow answer saves")
	}
	if canSaveAttemptAnswers(consts.ExamAttemptSubmitted) {
		t.Fatal("submitted attempt should not allow answer saves")
	}
	if !canSubmitAttempt(consts.ExamAttemptInProgress) {
		t.Fatal("in progress attempt should be submittable")
	}
	if canSubmitAttempt(consts.ExamAttemptEnded) {
		t.Fatal("scored attempt should not be submittable")
	}
	if !canFinalizeAttempt(consts.ExamAttemptSubmitted) {
		t.Fatal("submitted attempt should be finalizable")
	}
	if canFinalizeAttempt(consts.ExamAttemptEnded) {
		t.Fatal("scored attempt should not be finalized again")
	}
	if !canGradeSubjectiveAttempt(examentity.ExamAttempt{Status: consts.ExamAttemptEnded, HasSubjective: 1}) {
		t.Fatal("scored subjective attempt should be gradable")
	}
	if canGradeSubjectiveAttempt(examentity.ExamAttempt{Status: consts.ExamAttemptSubmitted, HasSubjective: 1}) {
		t.Fatal("submitted attempt should not be gradable before final scoring")
	}
}

func TestAttemptDeadlineReached(t *testing.T) {
	now := gtime.NewFromStr("2026-04-23 12:00:00")
	past := gtime.NewFromStr("2026-04-23 11:59:59")
	exact := gtime.NewFromStr("2026-04-23 12:00:00")
	future := gtime.NewFromStr("2026-04-23 12:00:01")

	if !isAttemptDeadlineReached(examentity.ExamAttempt{Status: consts.ExamAttemptInProgress, DeadlineAt: past}, now) {
		t.Fatal("past deadline should be reached")
	}
	if !isAttemptDeadlineReached(examentity.ExamAttempt{Status: consts.ExamAttemptInProgress, DeadlineAt: exact}, now) {
		t.Fatal("exact deadline should be reached")
	}
	if isAttemptDeadlineReached(examentity.ExamAttempt{Status: consts.ExamAttemptInProgress, DeadlineAt: future}, now) {
		t.Fatal("future deadline should not be reached")
	}
	if isAttemptDeadlineReached(examentity.ExamAttempt{Status: consts.ExamAttemptSubmitted, DeadlineAt: past}, now) {
		t.Fatal("non-progress attempt should not be checked as overdue")
	}
}

func TestBatchWindowOpen(t *testing.T) {
	start := gtime.NewFromStr("2026-04-23 12:00:00")
	end := gtime.NewFromStr("2026-04-23 13:00:00")

	if !isBatchWindowOpen(start, start, end) {
		t.Fatal("start boundary should be open")
	}
	if !isBatchWindowOpen(end, start, end) {
		t.Fatal("end boundary should be open")
	}
	if isBatchWindowOpen(gtime.NewFromStr("2026-04-23 11:59:59"), start, end) {
		t.Fatal("before start should be closed")
	}
	if isBatchWindowOpen(gtime.NewFromStr("2026-04-23 13:00:01"), start, end) {
		t.Fatal("after end should be closed")
	}
}
