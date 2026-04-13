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
			if err := doDatabaseSync(ctx, attemptID); err != nil {
				g.Log().Errorf(ctx, "同步失败: %v", err)
				// 失败重回队列逻辑...
			}
		}
	}
}

// doDatabaseSync 执行具体的数据库同步逻辑
func doDatabaseSync(ctx context.Context, attemptID int64) error {
	redisKey := fmt.Sprintf(consts.ExamAttemptKeyFmt, attemptID)

	// 1. 从 Redis Hash 中获取所有暂存的答案
	// HGetAll 返回的是 Map[string]string，其中 Key 是 QuestionID
	res, err := g.Redis().HGetAll(ctx, redisKey)
	if err != nil {
		return fmt.Errorf("failed to fetch from redis: %v", err)
	}
	if res.IsEmpty() {
		return nil
	}

	// 2. 开启数据库事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, val := range res.Map() {
			// 解析 Redis 中的 JSON 数据
			var item struct {
				Q int64  `json:"q"` // QuestionID
				A string `json:"a"` // AnswerJSON
				V *int   `json:"v"` // ExpectedVersion
				T int64  `json:"t"`
			}
			if err := gconv.Struct(val, &item); err != nil {
				g.Log().Errorf(ctx, "解析Redis数据失败: %v", err)
				continue
			}

			answeredTime := gtime.NewFromTimeStamp(item.T)

			// 3. 执行 Upsert (Save) 操作
			// 注意：这里需要数据库表满足 attempt_id + exam_question_id 的唯一索引
			// 如果冲突，则更新 answer_json 和 update_time
			_, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
				Data(g.Map{
					"attempt_id":       attemptID,
					"exam_question_id": item.Q,
					"answer_json":      item.A,
					"version":          gdb.Raw("version + 1"), // 简单逻辑：落库时版本自增
					"updater":          "system_async_worker",
					"update_time":      answeredTime,
					"delete_flag":      consts.DeleteFlagNotDeleted,
					"create_time":      answeredTime, // 仅在 Insert 时生效
					"creator":          "system_async_worker",
				}).
				OnConflict("attempt_id", "exam_question_id").
				Save()

			if err != nil {
				return fmt.Errorf("db save error for qid %d: %v", item.Q, err)
			}
		}
		return nil
	})
}
