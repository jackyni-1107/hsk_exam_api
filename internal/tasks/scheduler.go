package tasks

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"

	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/tasks/handler"
)

type cronEntry struct {
	name string
	expr string
}

type cronReloadPlan struct {
	removeIDs   []int64
	upsertTasks []sysentity.SysTask
}

var (
	cronEntries = make(map[int64]cronEntry)
	cronMu      sync.RWMutex
)

func triggerHandlerInit() error {
	_ = handler.Get(handler.DemoHandlerName)
	_ = handler.Get(handler.ExamScoreFinalizeHandlerName)
	_ = handler.Get(handler.ExamBatchExpireHandlerName)
	_ = handler.Get(handler.ExamDashboardStatsHandlerName)
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

	cronMu.Lock()
	plan := buildCronReloadPlan(cronEntries, list)
	for _, taskID := range plan.removeIDs {
		removeCronTaskLocked(ctx, taskID)
	}
	for i := range plan.upsertTasks {
		addCronTaskLocked(ctx, &plan.upsertTasks[i])
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
		Execute(ctx, ExecRequest{
			TaskID:      taskID,
			RunID:       newRunID(),
			TriggerType: consts.TriggerTypeCron,
			RetryCount:  0,
		})
	}, name)
	if err != nil {
		g.Log().Errorf(ctx, "add cron task failed id=%d: %v", t.Id, err)
		return
	}
	cronEntries[t.Id] = cronEntry{name: name, expr: t.CronExpr}
	g.Log().Infof(ctx, "[Scheduler] cron registered task_id=%d expr=%s", t.Id, t.CronExpr)
}

func removeCronTaskLocked(ctx context.Context, taskID int64) {
	entryMeta, ok := cronEntries[taskID]
	if !ok {
		return
	}
	entry := gcron.Search(entryMeta.name)
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

func buildCronReloadPlan(existing map[int64]cronEntry, tasks []sysentity.SysTask) cronReloadPlan {
	activeTasks := make(map[int64]sysentity.SysTask, len(tasks))
	for _, task := range tasks {
		activeTasks[task.Id] = task
	}

	plan := cronReloadPlan{
		removeIDs:   make([]int64, 0),
		upsertTasks: make([]sysentity.SysTask, 0),
	}

	for taskID, entry := range existing {
		task, ok := activeTasks[taskID]
		if !ok || entry.expr != task.CronExpr {
			plan.removeIDs = append(plan.removeIDs, taskID)
		}
	}

	for _, task := range tasks {
		entry, ok := existing[task.Id]
		if !ok || entry.expr != task.CronExpr {
			plan.upsertTasks = append(plan.upsertTasks, task)
		}
	}

	sort.Slice(plan.removeIDs, func(i, j int) bool {
		return plan.removeIDs[i] < plan.removeIDs[j]
	})
	sort.Slice(plan.upsertTasks, func(i, j int) bool {
		return plan.upsertTasks[i].Id < plan.upsertTasks[j].Id
	})
	return plan
}
