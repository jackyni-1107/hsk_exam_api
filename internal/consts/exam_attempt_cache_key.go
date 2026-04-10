package consts

// ExamAttemptCacheKeyPrefix 考试会话相关 Redis 键前缀（预留）
const ExamAttemptCacheKeyPrefix = "hskexam:exam:attempt_cache:"

// ---------- Exam Redis 键格式（hskexam:exam:业务名:唯一标识） ----------

const (
	HlsPlayKeyPrefix            = "hskexam:exam:hls_play:"
	ExamSubmitLockKeyFmt        = "hskexam:exam:submit_lock:%d"
	ExamSaveRateKeyFmt          = "hskexam:exam:save_rate:%d"
	ExamSubmitLockTTL           = 45
	PaperForExamCacheTTLSeconds = 3600
	PaperForExamMaxStringBytes  = 256 * 1024
)

const (
	ExamPaperInitCacheKeyFmt     = "hskexam:exam:paper_init:%d"
	ExamPaperLegacyCacheKeyFmt   = "hskexam:exam:paper_legacy:%d"
	ExamPaperSectionCacheKeyFmt  = "hskexam:exam:paper_section:%d:%d"
	ExamPaperSectionCachePattern = "hskexam:exam:paper_section:%d:*"
)
