package tasks

import (
	"context"
	"encoding/json"
	"time"

	"exam/internal/consts"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

const claimDueDelayedExecLua = `
local key = KEYS[1]
local now = ARGV[1]
local limit = tonumber(ARGV[2])
local members = redis.call('ZRANGEBYSCORE', key, '-inf', now, 'LIMIT', 0, limit)
if #members == 0 then
  return members
end
redis.call('ZREM', key, unpack(members))
return members
`

type delayedExecPayload struct {
	TaskID      int64  `json:"task_id"`
	RunID       string `json:"run_id"`
	TriggerType int    `json:"trigger_type"`
	RetryCount  int    `json:"retry_count"`
	Params      string `json:"params,omitempty"`
}

func delayQueueLoop(ctx context.Context) {
	dispatchDueDelayedExecs(ctx)
	ticker := time.NewTicker(time.Duration(consts.TaskDelayPollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dispatchDueDelayedExecs(ctx)
		}
	}
}

func dispatchDueDelayedExecs(ctx context.Context) {
	ok, err := tryAcquireDelayQueueScanLock(ctx)
	if err != nil {
		g.Log().Error(ctx, "acquire delay queue scan lock failed", err)
		return
	}
	if !ok {
		return
	}
	for {
		members, err := claimDueDelayedExecs(ctx, gtime.Now().TimestampMilli(), consts.TaskDelayBatchSize)
		if err != nil {
			g.Log().Error(ctx, "claim due delayed execs failed", err)
			return
		}
		if len(members) == 0 {
			return
		}
		for _, member := range members {
			req, err := decodeDelayedExec(member)
			if err != nil {
				g.Log().Error(ctx, "decode delayed exec failed", err)
				continue
			}
			executeAsync(req)
		}
		if len(members) < consts.TaskDelayBatchSize {
			return
		}
	}
}

func enqueueDelayedExec(ctx context.Context, req ExecRequest, delaySec int) error {
	member, err := encodeDelayedExec(req)
	if err != nil {
		return err
	}
	dueAt := gtime.Now().TimestampMilli() + int64(delaySec)*1000
	_, err = g.Redis().Do(ctx, "ZADD", consts.TaskDelayQueueKey, dueAt, member)
	return err
}

func tryAcquireDelayQueueScanLock(ctx context.Context) (bool, error) {
	v, err := g.Redis().Do(
		ctx,
		"SET",
		consts.TaskDelayScannerLockKey,
		"1",
		"NX",
		"PX",
		consts.TaskDelayScannerLockTTLMillis,
	)
	if err != nil {
		return false, err
	}
	return gconv.String(v.Val()) == "OK", nil
}

func scheduleAsyncFallback(ctx context.Context, req ExecRequest, delaySec int) {
	go func() {
		timer := time.NewTimer(time.Duration(delaySec) * time.Second)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			executeAsync(req)
		}
	}()
}

func claimDueDelayedExecs(ctx context.Context, nowMs int64, limit int) ([]string, error) {
	v, err := g.Redis().Do(ctx, "EVAL", claimDueDelayedExecLua, 1, consts.TaskDelayQueueKey, nowMs, limit)
	if err != nil {
		return nil, err
	}
	return gvar.New(v.Val()).Strings(), nil
}

func encodeDelayedExec(req ExecRequest) (string, error) {
	payload, err := json.Marshal(delayedExecPayload{
		TaskID:      req.TaskID,
		RunID:       req.RunID,
		TriggerType: req.TriggerType,
		RetryCount:  req.RetryCount,
		Params:      req.Params,
	})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func decodeDelayedExec(member string) (ExecRequest, error) {
	var payload delayedExecPayload
	if err := json.Unmarshal([]byte(member), &payload); err != nil {
		return ExecRequest{}, err
	}
	return ExecRequest{
		TaskID:      payload.TaskID,
		RunID:       payload.RunID,
		TriggerType: payload.TriggerType,
		RetryCount:  payload.RetryCount,
		Params:      payload.Params,
	}, nil
}
