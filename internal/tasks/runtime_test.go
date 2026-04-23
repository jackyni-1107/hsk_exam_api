package tasks

import (
	"testing"

	"exam/internal/consts"
	sysentity "exam/internal/model/entity/sys"
)

func TestBuildDispatchRequestForDelayTask(t *testing.T) {
	req, delaySec := buildDispatchRequest(&sysentity.SysTask{
		Id:           7,
		Type:         consts.TaskTypeDelay,
		DelaySeconds: 90,
	}, "run-delay")
	if req.TaskID != 7 || req.RunID != "run-delay" {
		t.Fatalf("unexpected dispatch request: %#v", req)
	}
	if req.TriggerType != consts.TriggerTypeDelay {
		t.Fatalf("expected delay trigger type, got %d", req.TriggerType)
	}
	if delaySec != 90 {
		t.Fatalf("expected delay seconds 90, got %d", delaySec)
	}
}

func TestBuildDispatchRequestForImmediateTask(t *testing.T) {
	req, delaySec := buildDispatchRequest(&sysentity.SysTask{
		Id:           8,
		Type:         consts.TaskTypeCron,
		DelaySeconds: 120,
	}, "run-now")
	if req.TaskID != 8 || req.RunID != "run-now" {
		t.Fatalf("unexpected dispatch request: %#v", req)
	}
	if req.TriggerType != consts.TriggerTypeManual {
		t.Fatalf("expected manual trigger type, got %d", req.TriggerType)
	}
	if delaySec != 0 {
		t.Fatalf("expected immediate dispatch, got delay %d", delaySec)
	}
}
