package attempt

import (
	"exam/internal/consts"
	examentity "exam/internal/model/entity/exam"

	"github.com/gogf/gf/v2/os/gtime"
)

type attemptState int

const (
	attemptStateUnknown    attemptState = 0
	attemptStateNotStarted attemptState = consts.ExamAttemptNotStarted
	attemptStateInProgress attemptState = consts.ExamAttemptInProgress
	attemptStateSubmitted  attemptState = consts.ExamAttemptSubmitted
	attemptStateScored     attemptState = consts.ExamAttemptEnded
)

type attemptEvent string

const (
	attemptEventStart        attemptEvent = "start"
	attemptEventSaveAnswers  attemptEvent = "save_answers"
	attemptEventSubmit       attemptEvent = "submit"
	attemptEventTimeout      attemptEvent = "timeout"
	attemptEventBatchExpired attemptEvent = "batch_expired"
	attemptEventFinalize     attemptEvent = "finalize"
	attemptEventGrade        attemptEvent = "grade"
)

type attemptTransition struct {
	From attemptState
	On   attemptEvent
	To   attemptState
}

var attemptTransitions = []attemptTransition{
	{From: attemptStateNotStarted, On: attemptEventStart, To: attemptStateInProgress},
	{From: attemptStateInProgress, On: attemptEventSaveAnswers, To: attemptStateInProgress},
	{From: attemptStateInProgress, On: attemptEventSubmit, To: attemptStateSubmitted},
	{From: attemptStateInProgress, On: attemptEventTimeout, To: attemptStateSubmitted},
	{From: attemptStateInProgress, On: attemptEventBatchExpired, To: attemptStateSubmitted},
	{From: attemptStateSubmitted, On: attemptEventFinalize, To: attemptStateScored},
	{From: attemptStateScored, On: attemptEventGrade, To: attemptStateScored},
}

func (state attemptState) valid() bool {
	switch state {
	case attemptStateNotStarted, attemptStateInProgress, attemptStateSubmitted, attemptStateScored:
		return true
	default:
		return false
	}
}

func nextAttemptState(status int, event attemptEvent) (attemptState, bool) {
	current := attemptState(status)
	if !current.valid() {
		return attemptStateUnknown, false
	}
	for _, t := range attemptTransitions {
		if t.From == current && t.On == event {
			return t.To, true
		}
	}
	return current, false
}

func canAttempt(status int, event attemptEvent) bool {
	_, ok := nextAttemptState(status, event)
	return ok
}

func transitionAttemptStatus(status int, event attemptEvent) (int, bool) {
	state, ok := nextAttemptState(status, event)
	if !ok {
		return status, false
	}
	return int(state), true
}

func isAttemptInProgress(status int) bool {
	return attemptState(status) == attemptStateInProgress
}

func isAttemptSubmitted(status int) bool {
	return attemptState(status) == attemptStateSubmitted
}

func isAttemptScored(status int) bool {
	return attemptState(status) == attemptStateScored
}

func isAttemptSubmittedOrScored(status int) bool {
	state := attemptState(status)
	return state == attemptStateSubmitted || state == attemptStateScored
}

func canStartAttempt(status int) bool {
	return canAttempt(status, attemptEventStart)
}

func canSaveAttemptAnswers(status int) bool {
	return canAttempt(status, attemptEventSaveAnswers)
}

func canSubmitAttempt(status int) bool {
	return canAttempt(status, attemptEventSubmit)
}

func canFinalizeAttempt(status int) bool {
	return canAttempt(status, attemptEventFinalize)
}

func canGradeSubjectiveAttempt(att examentity.ExamAttempt) bool {
	return canAttempt(att.Status, attemptEventGrade) && att.HasSubjective == 1
}

func isAttemptDeadlineReached(att examentity.ExamAttempt, now *gtime.Time) bool {
	if !canAttempt(att.Status, attemptEventTimeout) || att.DeadlineAt == nil || now == nil {
		return false
	}
	return !att.DeadlineAt.After(now)
}

func isBatchWindowOpen(now, start, end *gtime.Time) bool {
	if start == nil || end == nil || now == nil {
		return false
	}
	return !now.Before(start) && !now.After(end)
}
