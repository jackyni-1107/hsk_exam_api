package bo

const (
	defaultExamDurationSeconds      = 3600
	defaultExamMaxDurationSeconds   = 14400
	defaultExamSaveAnswersPerSecond = 20
)

// ExamCfg 对应 manifest/config 中 exam 节点（含 HLS 相关可选字段）。
type ExamCfg struct {
	DefaultDurationSeconds    int  `json:"defaultDurationSeconds"`
	MaxDurationSeconds        int  `json:"maxDurationSeconds"`
	SaveAnswersPerSecond      int  `json:"saveAnswersPerSecond"`
	EnableRandomAnswerHelper  bool `json:"enableRandomAnswerHelper"`
	AudioHlsTicketTTLSeconds  int  `json:"audioHlsTicketTTLSeconds"`
	AudioHlsPresignTTLSeconds int  `json:"audioHlsPresignTTLSeconds"`
}

func (c ExamCfg) Normalize() ExamCfg {
	if c.DefaultDurationSeconds <= 0 {
		c.DefaultDurationSeconds = defaultExamDurationSeconds
	}
	if c.MaxDurationSeconds <= 0 {
		c.MaxDurationSeconds = defaultExamMaxDurationSeconds
	}
	if c.SaveAnswersPerSecond <= 0 {
		c.SaveAnswersPerSecond = defaultExamSaveAnswersPerSecond
	}
	return c
}

func (c ExamCfg) ResolveDurationSeconds(paperDuration int, clientOverride int) int {
	c = c.Normalize()

	duration := paperDuration
	if duration <= 0 {
		duration = c.DefaultDurationSeconds
	}
	if clientOverride > 0 {
		duration = clientOverride
	}
	if duration > c.MaxDurationSeconds {
		duration = c.MaxDurationSeconds
	}
	if duration <= 0 {
		duration = c.DefaultDurationSeconds
	}
	return duration
}
