package attempt

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
)

const (
	dashboardSnapshotTable    = "exam_dashboard_stats_snapshot"
	dashboardSnapshotScopeKey = "global"
)

// AttemptAdminStats 无筛时读全量快照；有筛时即时计算。
func (s *sAttempt) AttemptAdminStats(ctx context.Context, level string, examinationPaperId, examBatchId, mockLevelId int64) (*bo.AttemptAdminStatsView, error) {
	empty := level == "" && examinationPaperId == 0 && examBatchId == 0 && mockLevelId == 0
	if empty {
		if v, ok, err := s.readGlobalDashboardSnapshot(ctx); err != nil {
			return nil, err
		} else if ok {
			return v, nil
		}
	}
	return s.computeAttemptAdminStatsView(ctx, AttemptAdminListQuery{
		Level:              level,
		MockLevelId:        mockLevelId,
		ExaminationPaperId: examinationPaperId,
		ExamBatchId:        examBatchId,
		BatchKind:          consts.ExamBatchKindFilterAll,
	})
}

// RefreshAttemptDashboardSnapshot 定时任务：全量重算并写入快照表。
func (s *sAttempt) RefreshAttemptDashboardSnapshot(ctx context.Context) error {
	v, err := s.computeAttemptAdminStatsView(ctx, AttemptAdminListQuery{BatchKind: consts.ExamBatchKindFilterAll})
	if err != nil {
		return err
	}
	pl := toSnapshotPayload(v)
	if err := s.upsertGlobalDashboardSnapshot(ctx, pl, gtime.Now()); err != nil {
		return err
	}
	return nil
}

func toSnapshotPayload(v *bo.AttemptAdminStatsView) *bo.AttemptAdminStatsSnapshotPayload {
	return &bo.AttemptAdminStatsSnapshotPayload{
		StatusNotStarted:  v.StatusNotStarted,
		StatusInProgress:  v.StatusInProgress,
		StatusSubmitted:   v.StatusSubmitted,
		StatusEnded:       v.StatusEnded,
		Total:             v.Total,
		FinishedCount:     v.FinishedCount,
		SubjectivePending: v.SubjectivePending,
		TodayNew:          v.TodayNew,
		CompletionRate:    v.CompletionRate,
		AvgObjective:      v.AvgObjective,
		AvgSubjective:     v.AvgSubjective,
		AvgTotal:          v.AvgTotal,
		Trend7d:           v.Trend7d,
		Buckets:           v.Buckets,
	}
}

func (s *sAttempt) readGlobalDashboardSnapshot(ctx context.Context) (*bo.AttemptAdminStatsView, bool, error) {
	var row struct {
		Payload    []byte      `orm:"payload"`
		UpdateTime *gtime.Time `orm:"update_time"`
	}
	// 勿在 SQL 中写 LIMIT 1：gdb.Raw().Scan 扫到 struct 时会自动追加 LIMIT 1，重复会触发 MySQL 1064
	err := g.DB().Ctx(ctx).Raw(
		`SELECT payload, update_time FROM `+dashboardSnapshotTable+` WHERE scope_key = ?`,
		dashboardSnapshotScopeKey,
	).Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if row.Payload == nil || len(row.Payload) == 0 {
		return nil, false, nil
	}
	var p bo.AttemptAdminStatsSnapshotPayload
	if err := json.Unmarshal(row.Payload, &p); err != nil {
		return nil, false, err
	}
	out := &bo.AttemptAdminStatsView{
		UpdatedAt:         row.UpdateTime,
		FromCache:         true,
		StatusNotStarted:  p.StatusNotStarted,
		StatusInProgress:  p.StatusInProgress,
		StatusSubmitted:   p.StatusSubmitted,
		StatusEnded:       p.StatusEnded,
		Total:             p.Total,
		FinishedCount:     p.FinishedCount,
		SubjectivePending: p.SubjectivePending,
		TodayNew:          p.TodayNew,
		CompletionRate:    p.CompletionRate,
		AvgObjective:      p.AvgObjective,
		AvgSubjective:     p.AvgSubjective,
		AvgTotal:          p.AvgTotal,
		Trend7d:           p.Trend7d,
		Buckets:           p.Buckets,
	}
	return out, true, nil
}

func (s *sAttempt) upsertGlobalDashboardSnapshot(ctx context.Context, p *bo.AttemptAdminStatsSnapshotPayload, at *gtime.Time) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = g.DB().Exec(ctx,
		`INSERT INTO `+dashboardSnapshotTable+` (scope_key, payload, update_time) VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE payload = VALUES(payload), update_time = VALUES(update_time)`,
		dashboardSnapshotScopeKey, string(b), at,
	)
	return err
}

// computeAttemptAdminStatsView 即时计算（不读快照）。不含 FromCache 语义，FromCache 由调用方设置。
func (s *sAttempt) computeAttemptAdminStatsView(ctx context.Context, q AttemptAdminListQuery) (*bo.AttemptAdminStatsView, error) {
	from, joinArgs, wArgs := q.buildAttemptAdminListFrom()
	bind := attemptAdminListCountArgs(joinArgs, wArgs)
	now := time.Now()
	loc := now.Location()
	day0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	day1 := day0.Add(24 * time.Hour)
	start7 := day0.AddDate(0, 0, -6)
	// 聚合
	var agg struct {
		S1    int     `orm:"s1"`
		S2    int     `orm:"s2"`
		S3    int     `orm:"s3"`
		S4    int     `orm:"s4"`
		Total int     `orm:"tot"`
		SF    int     `orm:"sf"`
		SP    int     `orm:"sp"`
		TN    int     `orm:"tn"`
		Avgo  float64 `orm:"avgo"`
		Avgs  float64 `orm:"avgs"`
		Avgt  float64 `orm:"avgt"`
	}
	sq1 := `SELECT
 COALESCE(SUM(CASE WHEN r.status = 1 THEN 1 ELSE 0 END),0) AS s1,
 COALESCE(SUM(CASE WHEN r.status = 2 THEN 1 ELSE 0 END),0) AS s2,
 COALESCE(SUM(CASE WHEN r.status = 3 THEN 1 ELSE 0 END),0) AS s3,
 COALESCE(SUM(CASE WHEN r.status IN (4,5) THEN 1 ELSE 0 END),0) AS s4,
 COALESCE(SUM(CASE WHEN r.status IN (3,4,5) THEN 1 ELSE 0 END),0) AS sf,
 COALESCE(SUM(CASE WHEN r.has_subjective = 1 AND r.status = 4 THEN 1 ELSE 0 END),0) AS sp,
 COALESCE(COUNT(1),0) AS tot,
 COALESCE(SUM(CASE WHEN r.create_time >= ? AND r.create_time < ? THEN 1 ELSE 0 END),0) AS tn,
 COALESCE(AVG(CASE WHEN r.status IN (3,4,5) THEN r.objective_score END),0) AS avgo,
 COALESCE(AVG(CASE WHEN r.status IN (3,4,5) THEN r.subjective_score END),0) AS avgs,
 COALESCE(AVG(CASE WHEN r.status IN (3,4,5) THEN r.total_score END),0) AS avgt` + from
	args1 := append([]interface{}{day0, day1}, bind...)
	if err := g.DB().Ctx(ctx).Raw(sq1, args1...).Scan(&agg); err != nil {
		return nil, err
	}
	// 趋势
	var trows []struct {
		D   string `orm:"d"`
		Cnt int    `orm:"c"`
	}
	sq2 := `SELECT DATE(r.create_time) AS d, COALESCE(COUNT(1),0) AS c` + from + ` AND r.create_time >= ? AND r.create_time < ?
GROUP BY DATE(r.create_time) ORDER BY d`
	args2 := append(append([]interface{}{}, bind...), start7, day1)
	if err := g.DB().Ctx(ctx).Raw(sq2, args2...).Scan(&trows); err != nil {
		return nil, err
	}
	byDate := make(map[string]int, len(trows))
	for _, r := range trows {
		byDate[r.D] = r.Cnt
	}
	trend7 := make([]bo.AttemptStatsDayPoint, 0, 7)
	for i := 0; i < 7; i++ {
		d := start7.AddDate(0, 0, i)
		ds := d.Format("2006-01-02")
		trend7 = append(trend7, bo.AttemptStatsDayPoint{Date: ds, Count: byDate[ds]})
	}
	// 分档
	var brows []struct {
		Bl  float64 `orm:"bl"`
		Cnt int     `orm:"c"`
	}
	sq3 := `SELECT FLOOR(r.total_score / 10) * 10 AS bl, COALESCE(COUNT(1),0) AS c` + from + ` AND r.status IN (3,4,5) AND (r.total_score IS NOT NULL)
GROUP BY bl ORDER BY bl`
	if err := g.DB().Ctx(ctx).Raw(sq3, bind...).Scan(&brows); err != nil {
		return nil, err
	}
	buckets := make([]bo.AttemptStatsScoreChunk, 0, len(brows))
	for _, r := range brows {
		if math.IsNaN(r.Bl) {
			continue
		}
		buckets = append(buckets, bo.AttemptStatsScoreChunk{BucketLow: r.Bl, Count: r.Cnt})
	}
	sort.Slice(buckets, func(i, j int) bool { return buckets[i].BucketLow < buckets[j].BucketLow })

	var comp float64
	if agg.Total > 0 {
		comp = float64(agg.SF) / float64(agg.Total)
	}
	out := &bo.AttemptAdminStatsView{
		UpdatedAt:         gtime.Now(),
		FromCache:         false,
		StatusNotStarted:  agg.S1,
		StatusInProgress:  agg.S2,
		StatusSubmitted:   agg.S3,
		StatusEnded:       agg.S4,
		Total:             agg.Total,
		FinishedCount:     agg.SF,
		SubjectivePending: agg.SP,
		TodayNew:          agg.TN,
		CompletionRate:    comp,
		AvgObjective:      agg.Avgo,
		AvgSubjective:     agg.Avgs,
		AvgTotal:          agg.Avgt,
		Trend7d:           trend7,
		Buckets:           buckets,
	}
	s.fillBatchMeta(ctx, out, q.ExamBatchId)
	return out, nil
}

func (s *sAttempt) fillBatchMeta(ctx context.Context, out *bo.AttemptAdminStatsView, examBatchId int64) {
	if examBatchId <= 0 {
		return
	}
	n, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", examBatchId).Count()
	if err != nil {
		return
	}
	out.BatchMemberCount = n
	if n > 0 {
		out.BatchCompletionRate = float64(out.FinishedCount) / float64(n)
	}
}
