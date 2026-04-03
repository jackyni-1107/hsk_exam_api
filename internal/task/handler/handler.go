package handler

import (
	"context"
)

// Handler 任务处理器接口
type Handler interface {
	Execute(ctx context.Context, taskID int64, params string) error
}
