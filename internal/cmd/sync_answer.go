package cmd

import (
	"context"
	"errors"
	"exam/internal/consts"
	"exam/internal/utility/examutil"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// SyncAnswer 启动消费者
func SyncAnswer(ctx context.Context) {
	g.Log().Info(ctx, "答案同步消费者启动...")

	queueKey := consts.ExamAttemptSyncQueueKey

	// worker 数量（根据机器配置调整）
	const workerNum = 5

	for i := 0; i < workerNum; i++ {
		go worker(ctx, queueKey, i)
	}

	// 阻塞等待退出
	<-ctx.Done()
	g.Log().Info(ctx, "SyncAnswer 已退出")
}

// worker 消费队列
func worker(ctx context.Context, queueKey string, workerID int) {
	g.Log().Infof(ctx, "worker-%d 启动", workerID)

	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "worker-%d 退出", workerID)
			return
		default:
			// 阻塞获取队列数据
			result, err := g.Redis().BLPop(ctx, 5, queueKey)
			if err != nil {
				// 超时正常
				if errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				g.Log().Errorf(ctx, "worker-%d Redis error: %v", workerID, err)
				time.Sleep(time.Second)
				continue
			}

			if len(result) < 2 {
				continue
			}

			attemptID := result[1].Int64()

			// 执行同步
			if err := doDatabaseSync(ctx, attemptID); err != nil {
				g.Log().Errorf(ctx, "worker-%d AttemptID %d 同步失败: %v", workerID, attemptID, err)

				// 👉 可选：失败重新入队（避免丢数据）
				_, _ = g.Redis().LPush(ctx, queueKey, attemptID)

				continue
			}

			g.Log().Infof(ctx, "worker-%d AttemptID %d 同步成功", workerID, attemptID)
		}
	}
}

func tryAcquireAttemptSyncLock(ctx context.Context, attemptID int64) (bool, error) {
	key := fmt.Sprintf(consts.ExamSubmitLockKeyFmt, attemptID)
	v, err := g.Redis().Do(ctx, "SET", key, "1", "NX", "EX", consts.ExamSubmitLockTTL)
	if err != nil {
		return false, err
	}
	return v.String() == "OK", nil
}

func releaseAttemptSyncLock(ctx context.Context, attemptID int64) {
	_, _ = g.Redis().Del(ctx, fmt.Sprintf(consts.ExamSubmitLockKeyFmt, attemptID))
}

// doDatabaseSync 执行数据库同步
func doDatabaseSync(ctx context.Context, attemptID int64) error {
	ok, err := tryAcquireAttemptSyncLock(ctx, attemptID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("attempt %d sync lock busy", attemptID)
	}
	defer releaseAttemptSyncLock(ctx, attemptID)

	redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)

	// 1. 获取 Redis 数据
	res, err := g.Redis().HGetAll(ctx, redisKey)
	if err != nil {
		return err
	}
	// 如果 Redis 数据已被其他协程清理或已过期，直接结束
	if res.IsEmpty() {
		g.Log().Debugf(ctx, "AttemptID %d Redis数据为空，跳过同步", attemptID)
		return nil
	}

	rawMap := res.Map()
	draftMap := make(map[string]string, len(rawMap))
	for k, val := range rawMap {
		draftMap[gconv.String(k)] = gconv.String(val)
	}
	items := examutil.BuildAttemptAnswerDraftRows(attemptID, draftMap, "system_async_worker")

	if len(items) == 0 {
		return nil
	}

	// 3. 同步写库（带版本控制）
	if err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return examutil.UpsertAttemptAnswerDraftRowsTx(ctx, tx, items)
	}); err != nil {
		return err
	}
	_, err = g.Redis().Expire(ctx, redisKey, 60)
	return err
}
