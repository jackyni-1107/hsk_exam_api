package handler

import (
	"context"
	"sync"
)

var (
	reg = map[string]Handler{}
	mu  sync.RWMutex
)

// Register 注册任务处理器（进程启动时注册）。
func Register(name string, h Handler) {
	mu.Lock()
	defer mu.Unlock()
	reg[name] = h
}

// Get 按名称取处理器。
func Get(name string) Handler {
	mu.RLock()
	defer mu.RUnlock()
	return reg[name]
}

// Execute 执行已注册处理器。
func Execute(ctx context.Context, handlerName string, taskID int64, params string) error {
	h := Get(handlerName)
	if h == nil {
		return ErrHandlerNotFound
	}
	return h.Execute(ctx, taskID, params)
}

const (
	DemoHandlerName               = "DemoHandler"
	ExamScoreFinalizeHandlerName  = "ExamScoreFinalizeHandler"
)

func init() {
	Register(DemoHandlerName, demoHandler{})
	Register(ExamScoreFinalizeHandlerName, examScoreFinalizeHandler{})
}
