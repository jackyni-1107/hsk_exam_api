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
		Where("eb.exam_start_at <= ?", gtime.Now()).
		Where("eb.exam_end_at >= ?", gtime.Now()).
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
		Id                     int64 `orm:"id"`
		ExamBatchId            int64 `orm:"exam_batch_id"`
		MockExaminationPaperId int64 `orm:"mock_examination_paper_id"`
		Status                 int   `orm:"status"`
	}
	err = dao.ExamAttempt.Ctx(ctx).
		Fields("id, exam_batch_id, mock_examination_paper_id, status").
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

	// 3. 将 attempt 记录映射到 Map 中
	// Key: BatchId_PaperId
	attemptMap := make(map[string]struct {
		Id     int64
		Status int
	}, len(attempts))

	for _, a := range attempts {
		key := fmt.Sprintf("%d_%d", a.ExamBatchId, a.MockExaminationPaperId)
		attemptMap[key] = struct {
			Id     int64
			Status int
		}{Id: a.Id, Status: a.Status}
	}

	// 4. 组装结果并过滤已完成的考试
	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {

		key := fmt.Sprintf("%d_%d", r.BatchId, r.MockExaminationPaperId)
		att, hasAttempt := attemptMap[key]

		// 过滤：已提交或系统强制结束的考试不显示在“待考”列表
		if hasAttempt && (att.Status == consts.ExamAttemptSubmitted || att.Status == consts.ExamAttemptEnded) {
			continue
		}

		winStatus := s.GetExamWindowStatus(now, r.ExamStartAt, r.ExamEndAt)
		// 过滤：过期的且没考过的卷子不显示
		if winStatus == "closed" && !hasAttempt {
			continue
		}
		var attemptId int64
		if hasAttempt && att.Status == consts.ExamAttemptInProgress {
			attemptId = att.Id
		}

		list = append(list, bo.MyExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			MockExaminationPaperId: r.MockExaminationPaperId,
			PaperTitle:             r.PaperTitle,
			ExamStartAt:            r.ExamStartAt,
			ExamEndAt:              r.ExamEndAt,
			AttemptId:              attemptId,
			WindowStatus:           winStatus,
		})

	}

	return list, nil
}
