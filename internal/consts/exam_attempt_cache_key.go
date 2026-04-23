package consts

// ExamAttemptCacheKeyPrefix 考试会话相关 Redis 键前缀（预留）
const ExamAttemptKeyFmt = "hskexam:exam:attempt:%d"

const ExamAttemptSyncQueueKey = "exam:sync:queue"

// ---------- Exam Redis 键格式（hskexam:exam:业务名:唯一标识） ----------

const (
	HlsPlayKeyPrefix        = "hskexam:exam:hls_play:"
	ExamSubmitLockKeyFmt    = "hskexam:exam:submit_lock:%d"
	ExamSaveRateKeyFmt      = "hskexam:exam:save_rate:%d"
	ExamSegmentSaveKeyFmt   = "hskexam:exam:segment_save:%d"
	ExamAttemptCreateKeyFmt = "hskexam:exam:attempt_create:%d:%d:%d"

	ExamSubmitLockTTL           = 45
	PaperForExamCacheTTLSeconds = 3600
	PaperForExamMaxStringBytes  = 256 * 1024
)

const (
	ExamPaperInitCacheKeyFmt     = "hskexam:exam:paper_init:%d"
	ExamPaperLegacyCacheKeyFmt   = "hskexam:exam:paper_legacy:%d"
	ExamPaperSectionCacheKeyFmt  = "hskexam:exam:paper_section:%d:%d"
	ExamPaperSectionCachePattern = "hskexam:exam:paper_section:%d:*"

	ExamPaperSectionTopicCacheKeyFmt  = "hskexam:exam:paper_section_topic:%d:%d"
	ExamPaperSectionTopicCachePattern = "hskexam:exam:paper_section_topic:%d:*"
	ExamPaperPrepareCacheKeyFmt       = "hskexam:exam:paper_prepare:%d"
)
