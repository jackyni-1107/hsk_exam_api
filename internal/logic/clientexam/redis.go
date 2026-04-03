package clientexam

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
)

const submitLockTTL = 45

func submitLockKey(attemptID int64) string {
	return fmt.Sprintf("exam:attempt:%d:submit_lock", attemptID)
}

func saveRateKey(attemptID int64) string {
	return fmt.Sprintf("exam:attempt:%d:save_rate", attemptID)
}

// TryAcquireSubmitLock 交卷/超时自动交卷互斥，避免重复计分。
func TryAcquireSubmitLock(ctx context.Context, attemptID int64) (bool, error) {
	key := submitLockKey(attemptID)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", submitLockTTL)
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
		return gerror.NewCode(consts.CodeTooManyRequests, "")
	}
	return nil
}
