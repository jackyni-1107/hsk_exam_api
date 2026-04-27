package bo

import "testing"

func TestExamCfgNormalize(t *testing.T) {
	cfg := (ExamCfg{}).Normalize()
	if cfg.DefaultDurationSeconds != 3600 {
		t.Fatalf("DefaultDurationSeconds = %d", cfg.DefaultDurationSeconds)
	}
	if cfg.MaxDurationSeconds != 14400 {
		t.Fatalf("MaxDurationSeconds = %d", cfg.MaxDurationSeconds)
	}
	if cfg.SaveAnswersPerSecond != 20 {
		t.Fatalf("SaveAnswersPerSecond = %d", cfg.SaveAnswersPerSecond)
	}
}

func TestExamCfgResolveDurationSeconds(t *testing.T) {
	cfg := ExamCfg{
		DefaultDurationSeconds: 1800,
		MaxDurationSeconds:     2400,
	}

	if got := cfg.ResolveDurationSeconds(0, 0); got != 1800 {
		t.Fatalf("ResolveDurationSeconds default = %d", got)
	}
	if got := cfg.ResolveDurationSeconds(1500, 0); got != 1500 {
		t.Fatalf("ResolveDurationSeconds paper = %d", got)
	}
	if got := cfg.ResolveDurationSeconds(1500, 3000); got != 2400 {
		t.Fatalf("ResolveDurationSeconds max clamp = %d", got)
	}
}
