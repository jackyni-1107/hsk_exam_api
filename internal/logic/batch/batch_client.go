package batch

import (
	"context"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MyExamBatches 学员查询自己的批次
func (s *sBatch) MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error) {
	// 1. 获取学员报名的批次及关联的试卷信息
	// 使用 struct 局部变量接收，避免污染外部 model
	var rows []struct {
		BatchId                int64       `orm:"batch_id"`
		Title                  string      `orm:"title"`
		MockExaminationPaperId int64       `orm:"mock_examination_paper_id"`
		PaperTitle             string      `orm:"paper_title"`
		ExamStartAt            *gtime.Time `orm:"exam_start_at"`
		ExamEndAt              *gtime.Time `orm:"exam_end_at"`
	}

	err = g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		LeftJoin("exam_batch eb", "eb.id = ebm.batch_id").
		LeftJoin("exam_paper ep", "ep.mock_examination_paper_id = ebm.mock_examination_paper_id").
		Fields("eb.id AS batch_id, eb.title, ebm.mock_examination_paper_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at").
		Where("ebm.member_id = ?", memberID).
		Where("eb.exam_start_at > ?", gtime.Now()).
		Where(g.Map{
			"eb.delete_flag": consts.DeleteFlagNotDeleted,
			"ep.delete_flag": consts.DeleteFlagNotDeleted,
		}).
		OrderDesc("eb.id").
		Scan(&rows)

	if err != nil || len(rows) == 0 {
		return []bo.MyExamBatchItem{}, err
	}

	// 2. 批量获取该学员在这些批次下的答题记录状态
	// 既然有唯一索引 (member_id, exam_batch_id, mock_examination_paper_id)，直接查出 status 即可
	batchIDs := gdb.ListItemValuesUnique(rows, "BatchId")
	var attempts []struct {
		ExamBatchId            int64 `orm:"exam_batch_id"`
		MockExaminationPaperId int64 `orm:"mock_examination_paper_id"`
		Status                 int   `orm:"status"`
	}
	err = dao.ExamAttempt.Ctx(ctx).
		Fields("exam_batch_id, mock_examination_paper_id, status").
		Where(g.Map{
			"member_id":   memberID,
			"delete_flag": consts.DeleteFlagNotDeleted,
		}).
		WhereIn("exam_batch_id", batchIDs).
		WhereIn("status", []int{consts.ExamAttemptSubmitted, consts.ExamAttemptEnded}).
		Scan(&attempts)
	if err != nil {
		return nil, err
	}

	// 3. 过滤掉已完成的batchId
	// 使用空结构体 map[string]struct{} 性能最优，不占用内存
	finishedSet := make(map[string]struct{}, len(attempts))
	for _, a := range attempts {
		// 建议 Key 格式：批次ID_试卷ID
		key := fmt.Sprintf("%d_%d", a.ExamBatchId, a.MockExaminationPaperId)
		finishedSet[key] = struct{}{}
	}

	// 4. 组装结果并过滤已完成的考试
	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {

		// 生成同样的 Key 进行比对
		key := fmt.Sprintf("%d_%d", r.BatchId, r.MockExaminationPaperId)
		if _, exists := finishedSet[key]; exists {
			// 如果在黑名单（已提交或已结束）中，直接跳过，不返回给前端
			continue
		}

		// 5. [补充建议] 考试窗口状态判断
		winStatus := s.GetExamWindowStatus(now, r.ExamStartAt, r.ExamEndAt)

		// 如果业务要求“过期未考”也不显示，可以加上：
		if winStatus == "closed" {
			continue
		}

		list = append(list, bo.MyExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			MockExaminationPaperId: r.MockExaminationPaperId,
			PaperTitle:             r.PaperTitle,
			ExamStartAt:            r.ExamStartAt,
			ExamEndAt:              r.ExamEndAt,
			WindowStatus:           winStatus,
		})
	}

	return list, nil
}
