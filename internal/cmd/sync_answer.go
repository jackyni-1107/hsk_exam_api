package cmd

import (
	"context"
	"errors"
	"exam/internal/consts"
	"exam/internal/dao"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
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

			// 同步成功 → 删除缓存
			redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)

			_, err = g.Redis().Expire(ctx, redisKey, 60)
			if err != nil {
				g.Log().Errorf(ctx, "worker-%d 删除缓存失败 [%s]: %v", workerID, redisKey, err)
			} else {
				g.Log().Infof(ctx, "worker-%d AttemptID %d 同步成功", workerID, attemptID)
			}
		}
	}
}

// doDatabaseSync 执行数据库同步
func doDatabaseSync(ctx context.Context, attemptID int64) error {
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

	dataMap := res.Map()

	// 2. 构建批量数据（避免 struct 反射，提高性能）
	items := make([]g.Map, 0, len(dataMap))

	for _, val := range dataMap {
		itemMap := gconv.Map(val)
		q := gconv.Int64(itemMap["q"])
		if q == 0 {
			continue
		}

		answeredTime := gtime.NewFromTimeStamp(gconv.Int64(itemMap["t"]))
		items = append(items, g.Map{
			"attempt_id":       attemptID,
			"exam_question_id": q,
			"answer_json":      itemMap["a"],
			"version":          gconv.Int(itemMap["v"]),
			"updater":          "system_async_worker",
			"update_time":      answeredTime,
			"delete_flag":      consts.DeleteFlagNotDeleted,
			"create_time":      answeredTime,
			"creator":          "system_async_worker",
		})
	}

	if len(items) == 0 {
		return nil
	}

	// 3. 批量 Upsert（带版本控制）
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.ExamAttemptAnswer.Table()).
			Ctx(ctx).
			Data(items).
			Batch(100).
			OnDuplicate(gdb.Raw(`
			answer_json = IF(VALUES(version) >= version, VALUES(answer_json), answer_json),
			update_time = IF(VALUES(version) >= version, VALUES(update_time), update_time),
			version     = IF(VALUES(version) >= version, VALUES(version), version)
		`)).Save()

		return err
	})
}
