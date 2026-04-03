package exam

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/model/bo"
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

// MyExamBatches 当前学员在批次成员表中的考试批次（联 exam_batch、exam_paper），分页。
func (s *sExam) MyExamBatches(ctx context.Context, memberID int64, page, size int) (list []bo.MyExamBatchItem, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	base := g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		InnerJoin("exam_batch eb", "eb.id = ebm.batch_id").
		InnerJoin("exam_paper ep", "ep.mock_examination_paper_id = eb.mock_examination_paper_id").
		Where("ebm.member_id", memberID).
		Where("eb.delete_flag", consts.DeleteFlagNotDeleted).
		Where("ep.delete_flag", consts.DeleteFlagNotDeleted)

	total, err = base.Count()
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []bo.MyExamBatchItem{}, 0, nil
	}

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
		Fields("eb.id AS batch_id, eb.title AS title, eb.mock_examination_paper_id AS mock_examination_paper_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at").
		OrderDesc("eb.id").
		Page(page, size).
		Scan(&rows)
	if err != nil {
		return nil, 0, err
	}
	ids := make([]int64, len(rows))
	for i := range rows {
		ids[i] = rows[i].BatchId
	}
	levelMap, err := s.mockLevelsByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {
		lids := levelMap[r.BatchId]
		if lids == nil {
			lids = []int64{}
		}
		list = append(list, bo.MyExamBatchItem{
			BatchId:                r.BatchId,
			Title:                  r.Title,
			MockExaminationPaperId: r.MockExaminationPaperId,
			PaperTitle:             r.PaperTitle,
			ExamStartAt:            r.ExamStartAt,
			ExamEndAt:              r.ExamEndAt,
			MockLevelIds:           lids,
			WindowStatus:           examWindowStatus(now, r.ExamStartAt, r.ExamEndAt),
		})
	}
	return list, total, nil
}
