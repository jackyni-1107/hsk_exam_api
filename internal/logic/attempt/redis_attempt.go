package attempt

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
)

const (
	examAttemptDraftTTLSeconds int64 = 7200
	attemptLockRetryCount            = 10
	attemptLockRetryDelay            = 50 * time.Millisecond
)

func submitLockKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamSubmitLockKeyFmt, attemptID)
}

func attemptCreateLockKey(userID, batchID, paperID int64) string {
	return fmt.Sprintf(consts.ExamAttemptCreateKeyFmt, userID, batchID, paperID)
}

func saveRateKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamSaveRateKeyFmt, attemptID)
}

func attemptAnswersRedisKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)
}

func attemptSegmentSaveRedisKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamSegmentSaveKeyFmt, attemptID)
}

// TryAcquireSubmitLock 交卷/超时自动交卷互斥，避免重复计分。
func TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error) {
	key := submitLockKey(attemptID)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", consts.ExamSubmitLockTTL)
	if err != nil {
		return false, err
	}
	return v.String() == "OK", nil
}

// ReleaseSubmitLock 释放交卷锁。
func ReleaseSubmitLock(ctx context.Context, attemptID int64) {
	_, _ = g.Redis().Del(ctx, submitLockKey(attemptID))
}

func AcquireSubmitLockWithRetry(ctx context.Context, attemptID int64) (bool, error) {
	for i := 0; i < attemptLockRetryCount; i++ {
		ok, err := TryAcquireSubmitLock(ctx, attemptID)
		if err != nil || ok {
			return ok, err
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(attemptLockRetryDelay):
		}
	}
	return false, nil
}

func tryAcquireAttemptCreateLock(ctx context.Context, userID, batchID, paperID int64) (bool, error) {
	key := attemptCreateLockKey(userID, batchID, paperID)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", consts.ExamSubmitLockTTL)
	if err != nil {
		return false, err
	}
	return v.String() == "OK", nil
}

func releaseAttemptCreateLock(ctx context.Context, userID, batchID, paperID int64) {
	_, _ = g.Redis().Del(ctx, attemptCreateLockKey(userID, batchID, paperID))
}

// RateLimitSaveAnswers 单会话保存答案频率限制，超限返回 CodeTooManyRequests。
func RateLimitSaveAnswers(ctx context.Context, attemptID int64, perSecond int) error {
	if perSecond <= 0 {
		perSecond = 20
	}
	key := saveRateKey(attemptID)
	n, err := g.Redis().Incr(ctx, key)
	if err != nil {
		return err
	}
	if n == 1 {
		_, _ = g.Redis().Expire(ctx, key, 1)
	}
	if int(n) > perSecond {
		return gerror.NewCode(consts.CodeTooManyRequests)
	}
	return nil
}

// RedisHGetAllAttemptDrafts 读取答题草稿 Hash（field 为题目 id）。
func RedisHGetAllAttemptDrafts(ctx context.Context, attemptID int64) (map[string]string, error) {
	redisKey := attemptAnswersRedisKey(attemptID)
	res, err := g.Redis().HGetAll(ctx, redisKey)
	if err != nil {
		return nil, err
	}
	if res == nil || res.IsEmpty() {
		return nil, nil
	}
	raw := res.Map()
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = gconv.String(v)
	}
	return out, nil
}

// RedisMaxDraftSaveTime 草稿中各题的 t（Unix 秒）最大值。
func RedisMaxDraftSaveTime(ctx context.Context, attemptID int64) *gtime.Time {
	m, err := RedisHGetAllAttemptDrafts(ctx, attemptID)
	if err != nil || len(m) == 0 {
		return nil
	}
	var maxTs int64
	for _, s := range m {
		t := gjson.New(s).Get("t").Int64()
		if t > maxTs {
			maxTs = t
		}
	}
	if maxTs <= 0 {
		return nil
	}
	return gtime.NewFromTimeStamp(maxTs)
}

// RedisSaveSegmentSubmitTime 记录当前环节最近一次提交时间（Unix 秒）。
func RedisSaveSegmentSubmitTime(ctx context.Context, attemptID int64, segmentCode string, saveAtTs int64) error {
	if segmentCode == "" || saveAtTs <= 0 {
		return nil
	}
	redisKey := attemptSegmentSaveRedisKey(attemptID)
	if _, err := g.Redis().HSet(ctx, redisKey, map[string]any{segmentCode: saveAtTs}); err != nil {
		return err
	}
	_, _ = g.Redis().Expire(ctx, redisKey, examAttemptDraftTTLSeconds)
	return nil
}

// RedisGetSegmentLastSubmitTime 读取指定环节最近一次提交时间。
func RedisGetSegmentLastSubmitTime(ctx context.Context, attemptID int64, segmentCode string) *gtime.Time {
	if segmentCode == "" {
		return nil
	}
	redisKey := attemptSegmentSaveRedisKey(attemptID)
	v, err := g.Redis().HGet(ctx, redisKey, segmentCode)
	if err != nil || v.IsNil() {
		return nil
	}
	ts := v.Int64()
	if ts <= 0 {
		return nil
	}
	return gtime.NewFromTimeStamp(ts)
}

// RedisLatestSegmentCode 返回最近提交过答案的环节编码。
func RedisLatestSegmentCode(ctx context.Context, attemptID int64) string {
	redisKey := attemptSegmentSaveRedisKey(attemptID)
	res, err := g.Redis().HGetAll(ctx, redisKey)
	if err != nil || res == nil || res.IsEmpty() {
		return ""
	}
	items := res.MapStrVar()
	var (
		maxTs int64
		best  string
	)
	for code, tsVar := range items {
		ts := tsVar.Int64()
		if ts > maxTs {
			maxTs = ts
			best = code
		}
	}
	return best
}

// RedisSaveAnswerDrafts 写入草稿并设置过期、推同步队列。
func RedisSaveAnswerDrafts(ctx context.Context, attemptID int64, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	redisKey := attemptAnswersRedisKey(attemptID)
	if err := g.Redis().HMSet(ctx, redisKey, fields); err != nil {
		return err
	}
	_, _ = g.Redis().Expire(ctx, redisKey, examAttemptDraftTTLSeconds)
	if _, err := g.Redis().LPush(ctx, consts.ExamAttemptSyncQueueKey, attemptID); err != nil {
		g.Log().Error(ctx, "push exam attempt sync queue failed", err)
	}
	return nil
}
