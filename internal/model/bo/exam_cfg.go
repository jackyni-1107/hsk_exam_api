package bo

// ExamCfg 对应 manifest/config 中 exam 节点（含 HLS 相关可选字段）。
type ExamCfg struct {
	DefaultDurationSeconds    int  `json:"defaultDurationSeconds"`
	MaxDurationSeconds        int  `json:"maxDurationSeconds"`
	SaveAnswersPerSecond      int  `json:"saveAnswersPerSecond"`
	EnableRandomAnswerHelper  bool `json:"enableRandomAnswerHelper"`
	AudioHlsTicketTTLSeconds  int  `json:"audioHlsTicketTTLSeconds"`
	AudioHlsPresignTTLSeconds int  `json:"audioHlsPresignTTLSeconds"`
}
