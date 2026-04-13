package cmd

import (
	"context"
	"exam/internal/consts"
	"exam/internal/dao"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func SyncAnswer(ctx context.Context) {
	g.Log().Info(ctx, "答案同步消费者已在后台启动...")
	queueKey := consts.ExamAttemptSyncQueueKey

	for {
		select {
		case <-ctx.Done():
			// 如果外部上下文取消（比如程序关闭），退出循环
			g.Log().Info(ctx, "收到退出信号，停止同步消费者")
			return
		default:
			// BLPop 建议设置较短的超时时间（例如 5 秒），以便能定期检查 ctx.Done()
			result, err := g.Redis().BLPop(ctx, 5, queueKey)
			if err != nil {
				// 忽略由于超时导致的错误
				continue
			}
			if len(result) < 2 {
				continue
			}

			attemptID := result[1].Int64()
			// 1. 执行同步
			if err := doDatabaseSync(ctx, attemptID); err != nil {
				g.Log().Errorf(ctx, "AttemptID %d 同步失败: %v", attemptID, err)
				// 失败了不要删除 Redis 数据，甚至可以考虑把 attemptID 重新 LPush 回去或加入死信队列
				continue
			}

			// 2. 同步成功后，安全删除缓存
			// 建议：如果系统并发极高，可以使用 Expire 设置 60 秒过期，防止缓存击穿
			// 如果追求内存立即释放，直接使用 Del
			redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)
			_, err = g.Redis().Expire(ctx, redisKey, 60)
			if err != nil {
				g.Log().Errorf(ctx, "清理 Redis 缓存失败 [%s]: %v", redisKey, err)
			} else {
				g.Log().Infof(ctx, "AttemptID %d 同步成功并已清理缓存", attemptID)
			}
		}
	}
}

// doDatabaseSync 执行具体的数据库同步逻辑
func doDatabaseSync(ctx context.Context, attemptID int64) error {
	redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)

	// 1. 获取所有暂存答案
	res, err := g.Redis().HGetAll(ctx, redisKey)
	if err != nil || res.IsEmpty() {
		return err
	}

	// 2. 批量解析数据
	// 预分配切片空间，提高效率
	items := make([]g.Map, 0, len(res.Map()))
	for _, val := range res.Map() {
		var item struct {
			Q int64  `json:"q"`
			A string `json:"a"`
			V int    `json:"v"` // 调用方传来的预期版本
			T int64  `json:"t"`
		}
		if err := gconv.Struct(val, &item); err != nil {
			continue
		}

		answeredTime := gtime.NewFromTimeStamp(item.T)
		items = append(items, g.Map{
			"attempt_id":       attemptID,
			"exam_question_id": item.Q,
			"answer_json":      item.A,
			"version":          item.V, // 直接使用前端/Redis传来的版本
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

	// 3. 批量 Upsert 优化
	// 使用 OnConflict + Upsert 配合版本校验
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
			Data(items).
			// 每次处理 50 条数据，防止单条 SQL 过大
			Batch(50).
			// 关键：只有当提交的版本 >= 数据库现有版本时才更新
			// 注意：具体的 SQL 语法可能需要根据你的 DB (MySQL/PostgreSQL) 微调
			OnConflict("attempt_id", "exam_question_id").
			Save()

		return err
	})
}
