package clientexam

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// ExamCfg 考试时长与保存限流（manifest/config/config.yaml 中 exam 段）。
type ExamCfg struct {
	DefaultDurationSeconds   int
	MaxDurationSeconds       int
	SaveAnswersPerSecond     int
	EnableRandomAnswerHelper bool
}

// LoadExamCfg 读取配置，缺省为 3600 / 14400 / 20。
func LoadExamCfg(ctx context.Context) ExamCfg {
	c := g.Cfg()
	return ExamCfg{
		DefaultDurationSeconds:   c.MustGet(ctx, "exam.defaultDurationSeconds", 3600).Int(),
		MaxDurationSeconds:       c.MustGet(ctx, "exam.maxDurationSeconds", 14400).Int(),
		SaveAnswersPerSecond:     c.MustGet(ctx, "exam.saveAnswersPerSecond", 20).Int(),
		EnableRandomAnswerHelper: c.MustGet(ctx, "exam.enableRandomAnswerHelper", false).Bool(),
	}
}

// ResolveDurationSeconds 试卷时长优先，否则默认；客户端传入的时长会被 max 夹紧。
func ResolveDurationSeconds(cfg ExamCfg, paperDuration int, clientOverride int) int {
	d := paperDuration
	if d <= 0 {
		d = cfg.DefaultDurationSeconds
	}
	if clientOverride > 0 {
		d = clientOverride
	}
	if d > cfg.MaxDurationSeconds {
		d = cfg.MaxDurationSeconds
	}
	if d <= 0 {
		d = cfg.DefaultDurationSeconds
	}
	return d
}
