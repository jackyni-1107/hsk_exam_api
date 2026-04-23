package systask

import (
	"testing"

	"exam/internal/consts"
)

func TestBuildTaskMutationDataNormalizesCronTask(t *testing.T) {
	data, err := buildTaskMutationData(taskMutationInput{
		Name:          "  Cron Task  ",
		Code:          "  cron.task  ",
		CronExpr:      "  */5 * * * *  ",
		Handler:       " DemoHandler ",
		Type:          consts.TaskTypeCron,
		DelaySeconds:  30,
		RetryTimes:    -1,
		RetryInterval: -2,
		Concurrency:   -3,
		AlertOnFail:   7,
		Status:        99,
	})
	if err != nil {
		t.Fatalf("build cron task mutation data failed: %v", err)
	}
	if data.Name != "Cron Task" || data.Code != "cron.task" || data.CronExpr != "*/5 * * * *" || data.Handler != "DemoHandler" {
		t.Fatalf("unexpected trimmed task data: %#v", data)
	}
	if data.DelaySeconds != 0 {
		t.Fatalf("expected cron task delay seconds to be cleared, got %v", data.DelaySeconds)
	}
	if data.RetryTimes != 0 || data.RetryInterval != 0 || data.Concurrency != 0 {
		t.Fatalf("expected negative numeric fields to normalize to zero, got %#v", data)
	}
	if data.AlertOnFail != 0 || data.Status != consts.TaskStatusEnabled {
		t.Fatalf("expected normalized flags/status, got %#v", data)
	}
}

func TestBuildTaskMutationDataNormalizesDelayTask(t *testing.T) {
	data, err := buildTaskMutationData(taskMutationInput{
		Name:         "Delay",
		Code:         "delay.task",
		CronExpr:     "*/5 * * * *",
		Handler:      "DemoHandler",
		Type:         consts.TaskTypeDelay,
		DelaySeconds: 60,
		AlertOnFail:  1,
		Status:       consts.TaskStatusDisabled,
	})
	if err != nil {
		t.Fatalf("build delay task mutation data failed: %v", err)
	}
	if data.CronExpr != "" {
		t.Fatalf("expected delay task cron expression to be cleared, got %v", data.CronExpr)
	}
	if data.DelaySeconds != 60 || data.AlertOnFail != 1 || data.Status != consts.TaskStatusDisabled {
		t.Fatalf("unexpected delay task data: %#v", data)
	}
}

func TestNormalizeTaskOperatorFallsBackToSystem(t *testing.T) {
	if got := normalizeTaskOperator("   "); got != defaultTaskOperator {
		t.Fatalf("expected default task operator, got %q", got)
	}
	if got := normalizeTaskOperator(" admin "); got != "admin" {
		t.Fatalf("expected trimmed operator, got %q", got)
	}
}
