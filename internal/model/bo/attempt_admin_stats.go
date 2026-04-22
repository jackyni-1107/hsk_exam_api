package bo

import "github.com/gogf/gf/v2/os/gtime"

// AttemptAdminStatsView 管理端考试监控统计（与 API / 快照 JSON 一致）。
type AttemptAdminStatsView struct {
	UpdatedAt *gtime.Time `json:"updated_at"`
	// FromCache 无筛时为 true（读 exam_dashboard_stats_snapshot）
	FromCache bool `json:"from_cache"`

	StatusNotStarted  int     `json:"status_not_started"`
	StatusInProgress  int     `json:"status_in_progress"`
	StatusSubmitted   int     `json:"status_submitted"`
	StatusEnded       int     `json:"status_ended"`
	Total             int     `json:"total"`
	FinishedCount     int     `json:"finished_count"`
	SubjectivePending int     `json:"subjective_pending"`
	TodayNew          int     `json:"today_new"`
	CompletionRate    float64 `json:"completion_rate"` // finished/total

	AvgObjective  float64 `json:"avg_objective"`
	AvgSubjective float64 `json:"avg_subjective"`
	AvgTotal      float64 `json:"avg_total"`

	Trend7d []AttemptStatsDayPoint   `json:"trend_7d"`
	Buckets []AttemptStatsScoreChunk `json:"score_distribution"`

	BatchMemberCount    int     `json:"batch_member_count"`
	BatchCompletionRate float64 `json:"batch_completion_rate"` // finished_count / member_count
}

// AttemptStatsDayPoint 近 7 日（含今日）按自然日的新会话数。
type AttemptStatsDayPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// AttemptStatsScoreChunk 总分分档（floor(score/10)*10）。
type AttemptStatsScoreChunk struct {
	BucketLow float64 `json:"bucket_low"`
	Count     int     `json:"count"`
}

// AttemptAdminStatsSnapshotPayload 写入快照的 JSON 结构（不含 from_cache/updated_at，updated_at 用表字段）。
type AttemptAdminStatsSnapshotPayload struct {
	StatusNotStarted  int                      `json:"status_not_started"`
	StatusInProgress  int                      `json:"status_in_progress"`
	StatusSubmitted   int                      `json:"status_submitted"`
	StatusEnded       int                      `json:"status_ended"`
	Total             int                      `json:"total"`
	FinishedCount     int                      `json:"finished_count"`
	SubjectivePending int                      `json:"subjective_pending"`
	TodayNew          int                      `json:"today_new"`
	CompletionRate    float64                  `json:"completion_rate"`
	AvgObjective      float64                  `json:"avg_objective"`
	AvgSubjective     float64                  `json:"avg_subjective"`
	AvgTotal          float64                  `json:"avg_total"`
	Trend7d           []AttemptStatsDayPoint   `json:"trend_7d"`
	Buckets           []AttemptStatsScoreChunk `json:"score_distribution"`
}
