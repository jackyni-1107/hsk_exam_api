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
	var rows []struct {
		BatchId                int64       `orm:"batch_id"`
		Title                  string      `orm:"title"`
		ExamPaperId            int64       `orm:"exam_paper_id"`
		MockExaminationPaperId int64       `orm:"mock_src_id"`
		PaperTitle             string      `orm:"paper_title"`
		ExamStartAt            *gtime.Time `orm:"exam_start_at"`
		ExamEndAt              *gtime.Time `orm:"exam_end_at"`
	}

	err = g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		LeftJoin("exam_batch eb", "eb.id = ebm.batch_id").
		LeftJoin("exam_paper ep", "ep.id = ebm.exam_paper_id").
		Fields("eb.id AS batch_id, eb.title, ebm.exam_paper_id, ep.mock_examination_paper_id AS mock_src_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at").
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

	batchIDs := gdb.ListItemValuesUnique(rows, "BatchId")
	var attempts []struct {
		Id          int64 `orm:"id"`
		ExamBatchId int64 `orm:"exam_batch_id"`
		ExamPaperId int64 `orm:"exam_paper_id"`
		Status      int   `orm:"status"`
	}
	err = dao.ExamAttempt.Ctx(ctx).
		Fields("id, exam_batch_id, exam_paper_id, status").
		Where(g.Map{
			"member_id":   memberID,
			"delete_flag": consts.DeleteFlagNotDeleted,
		}).
		WhereIn("exam_batch_id", batchIDs).
		Scan(&attempts)
	if err != nil {
		return nil, err
	}

	attemptMap := make(map[string]struct {
		Id     int64
		Status int
	}, len(attempts))

	for _, a := range attempts {
		key := fmt.Sprintf("%d_%d", a.ExamBatchId, a.ExamPaperId)
		attemptMap[key] = struct {
			Id     int64
			Status int
		}{Id: a.Id, Status: a.Status}
	}

	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {

		key := fmt.Sprintf("%d_%d", r.BatchId, r.ExamPaperId)
		att, hasAttempt := attemptMap[key]

		if hasAttempt && (att.Status == consts.ExamAttemptSubmitted || att.Status == consts.ExamAttemptEnded) {
			continue
		}

		winStatus := s.GetExamWindowStatus(now, r.ExamStartAt, r.ExamEndAt)
		if winStatus == "closed" && !hasAttempt {
			continue
		}
		list = append(list, bo.MyExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			ExamPaperId:            r.ExamPaperId,
			MockExaminationPaperId: r.MockExaminationPaperId,
			PaperTitle:             r.PaperTitle,
			ExamStartAt:            r.ExamStartAt,
			ExamEndAt:              r.ExamEndAt,
			AttemptId:              att.Id,
			WindowStatus:           winStatus,
		})

	}

	return list, nil
}
