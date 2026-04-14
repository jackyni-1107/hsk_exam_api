package attempt

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
)

const (
	examAttemptDraftTTLSeconds int64 = 7200
)

func submitLockKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamSubmitLockKeyFmt, attemptID)
}

func saveRateKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamSaveRateKeyFmt, attemptID)
}

func attemptAnswersRedisKey(attemptID int64) string {
	return fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)
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
