package tasks

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/tasks/handler"
)

var (
	cronEntries = make(map[int64]string)
	cronMu      sync.RWMutex
)

// StartScheduler 启动定时任务调度：加载 Cron 任务并启动延迟 Worker。
func StartScheduler(ctx context.Context) {
	logClusterTimeZoneHint(ctx)
	_ = triggerHandlerInit()
	loadCronTasks(ctx)
	go refreshCronLoop(ctx)
	go delayWorkerLoop(ctx)
}

func triggerHandlerInit() error {
	// 触发 handler 包加载，确保内置 Handler 已注册
	_ = handler.Get(handler.DemoHandlerName)
	_ = handler.Get(handler.ExamScoreFinalizeHandlerName)
	return nil
}

func loadCronTasks(ctx context.Context) {
	var list []sysentity.SysTask
	err := sysdao.SysTask.Ctx(ctx).
		Where("type", consts.TaskTypeCron).
		Where("status", consts.TaskStatusEnabled).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&list)
	if err != nil {
		g.Log().Error(ctx, "load cron tasks failed", err)
		return
	}
	// 移除库中已删除任务的 cron 注册
	cronMu.Lock()
	toRemove := make([]int64, 0)
	for tid := range cronEntries {
		found := false
		for _, t := range list {
			if t.Id == tid {
				found = true
				break
			}
		}
		if !found {
			toRemove = append(toRemove, tid)
		}
	}
	for _, tid := range toRemove {
		removeCronTaskLocked(ctx, tid)
	}
	// 注册新任务
	for _, t := range list {
		if _, ok := cronEntries[t.Id]; !ok {
			addCronTaskLocked(ctx, &t)
		}
	}
	cronMu.Unlock()
}

func addCronTaskLocked(ctx context.Context, t *sysentity.SysTask) {
	if t.CronExpr == "" {
		return
	}
	name := "task-" + strconv.FormatInt(t.Id, 10)
	taskID := t.Id
	_, err := gcron.AddSingleton(ctx, t.CronExpr, func(ctx context.Context) {
		runID := gconv.String(gtime.Now().TimestampMilli())
		Execute(ctx, ExecRequest{
			TaskID:      taskID,
			RunID:       runID,
			TriggerType: consts.TriggerTypeCron,
			RetryCount:  0,
		})
	}, name)
	if err != nil {
		g.Log().Errorf(ctx, "add cron task failed id=%d: %v", t.Id, err)
		return
	}
	cronEntries[t.Id] = name
	g.Log().Infof(ctx, "[Scheduler] cron registered task_id=%d expr=%s", t.Id, t.CronExpr)
}

func removeCronTaskLocked(ctx context.Context, taskID int64) {
	name, ok := cronEntries[taskID]
	if !ok {
		return
	}
	entry := gcron.Search(name)
	if entry != nil {
		entry.Stop()
	}
	delete(cronEntries, taskID)
}

func removeCronTask(ctx context.Context, taskID int64) {
	cronMu.Lock()
	defer cronMu.Unlock()
	removeCronTaskLocked(ctx, taskID)
}

func refreshCronLoop(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			loadCronTasks(ctx)
		}
	}
}
