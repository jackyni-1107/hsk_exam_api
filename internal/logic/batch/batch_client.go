package exam

import (
	"context"
	"exam/internal/consts"
	"exam/internal/model/bo"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func examWindowStatus(now, start, end *gtime.Time) string {
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

// MyExamBatches 当前学员在批次成员表中的考试批次（联 exam_batch、exam_paper）。
// 返回会员批次列表，并过滤已提交/已结束会话，不分页。
func (s *sExam) MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error) {
	base := g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		InnerJoin("exam_batch eb", "eb.id = ebm.batch_id").
		InnerJoin("exam_paper ep", "ep.mock_examination_paper_id = ebm.mock_examination_paper_id").
		Where("ebm.member_id", memberID).
		Where("eb.delete_flag", consts.DeleteFlagNotDeleted).
		Where("ep.delete_flag", consts.DeleteFlagNotDeleted)

	type joinRow struct {
		BatchId                int64       `orm:"batch_id"`
		Title                  string      `orm:"title"`
		MockExaminationPaperId int64       `orm:"mock_examination_paper_id"`
		PaperTitle             string      `orm:"paper_title"`
		ExamStartAt            *gtime.Time `orm:"exam_start_at"`
		ExamEndAt              *gtime.Time `orm:"exam_end_at"`
	}
	var rows []joinRow
	err = base.Clone().
		Fields("eb.id AS batch_id, eb.title AS title, ebm.mock_examination_paper_id AS mock_examination_paper_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at").
		OrderDesc("eb.id").
		OrderDesc("ebm.mock_examination_paper_id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []bo.MyExamBatchItem{}, nil
	}

	ids := make([]int64, len(rows))
	for i := range rows {
		ids[i] = rows[i].BatchId
	}
	paperMap, err := s.mockPapersByBatchIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	type attemptAgg struct {
		ExamBatchId            int64 `orm:"exam_batch_id"`
		MockExaminationPaperId int64 `orm:"mock_examination_paper_id"`
		MaxStatus              int   `orm:"max_status"`
	}
	var attRows []attemptAgg
	if err := dao.ExamAttempt.Ctx(ctx).
		Fields("exam_batch_id, mock_examination_paper_id, MAX(status) AS max_status").
		Where("member_id", memberID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		WhereIn("exam_batch_id", ids).
		Group("exam_batch_id, mock_examination_paper_id").
		Scan(&attRows); err != nil {
		return nil, err
	}
	statusMap := make(map[string]int, len(attRows))
	for _, a := range attRows {
		k := fmt.Sprintf("%d:%d", a.ExamBatchId, a.MockExaminationPaperId)
		statusMap[k] = a.MaxStatus
	}

	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {
		winStatus := examWindowStatus(now, r.ExamStartAt, r.ExamEndAt)
		key := fmt.Sprintf("%d:%d", r.BatchId, r.MockExaminationPaperId)
		if st, ok := statusMap[key]; ok && (st == consts.ExamAttemptSubmitted || st == consts.ExamAttemptEnded) {
			continue
		}

		pids := paperMap[r.BatchId]
		if pids == nil {
			pids = []int64{}
		}
		list = append(list, bo.MyExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			MockExaminationPaperId: r.MockExaminationPaperId,
			//MockExaminationPaperIds: pids,
			PaperTitle:   r.PaperTitle,
			ExamStartAt:  r.ExamStartAt,
			ExamEndAt:    r.ExamEndAt,
			WindowStatus: winStatus,
		})
	}
	return list, nil
}
