package batch

import (
	"context"
	"exam/internal/consts"
	"exam/internal/model/bo"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MyExamBatches 学员查询自己的批次
func (s *sBatch) MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error) {
	// 使用统一的 DB Model 构建
	m := g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		InnerJoin("exam_batch eb", "eb.id = ebm.batch_id").
		Where("ebm.member_id", memberID).
		Where("eb.delete_flag", consts.DeleteFlagNotDeleted)

	var rows []struct {
		bo.MyExamBatchItem
		ExamStartAt *gtime.Time
		ExamEndAt   *gtime.Time
	}
	if err = m.Scan(&rows); err != nil {
		return nil, err
	}

	now := gtime.Now()
	for _, r := range rows {
		// 复用 init.go 中的状态判定逻辑
		r.WindowStatus = s.GetExamWindowStatus(now, r.ExamStartAt, r.ExamEndAt)
		list = append(list, r.MyExamBatchItem)
	}

	return list, nil
}
