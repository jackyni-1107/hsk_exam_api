package batch

import (
	"context"
	"strings"

	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/service/batch"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type sBatch struct{}

func init() {
	batch.RegisterBatch(New())
}

func New() *sBatch {
	return &sBatch{}
}

// --- 公用工具函数 (Internal) ---

// GetExamWindowStatus 判定考试窗口状态
func (s *sBatch) GetExamWindowStatus(now, start, end *gtime.Time) string {
	if start == nil || end == nil {
		return "closed"
	}
	if now.Before(start) {
		return "upcoming"
	}
	if now.After(end) {
		return "closed"
	}
	return "open"
}

// GetBatchByID 内部公用查询
func (s *sBatch) GetBatchByID(ctx context.Context, id int64) (*examentity.ExamBatch, error) {
	var out *examentity.ExamBatch
	err := dao.ExamBatch.Ctx(ctx).Where("id", id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&out)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, gerror.NewCode(consts.CodeExamBatchNotFound)
	}
	return out, nil
}

// dedupIDs ID去重工具
func (s *sBatch) dedupIDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// parseTime 统一时间解析
func (s *sBatch) parseTime(str string) *gtime.Time {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	return gtime.NewFromStr(str)
}
