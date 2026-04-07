package exam

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
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

// MyExamBatches 当前学员在批次成员表中的考试批次（联 exam_batch、exam_paper）。
// 返回会员批次列表，并过滤已提交/已结束会话，不分页。
func (s *sExam) MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error) {
	//todo 待优化
	base := g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		InnerJoin("exam_batch eb", "eb.id = ebm.batch_id").
		InnerJoin("exam_paper ep", "ep.mock_examination_paper_id = eb.mock_examination_paper_id").
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
		MockLevelId            int64       `orm:"mock_level_id"`
	}
	var rows []joinRow
	err = base.Clone().
		Fields("eb.id AS batch_id, eb.title AS title, eb.mock_examination_paper_id AS mock_examination_paper_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at, ebm.mock_level_id AS mock_level_id").
		OrderDesc("eb.id").
		OrderDesc("ebm.mock_level_id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []bo.MyExamBatchItem{}, nil
	}

	// 批量加载批次可选等级
	ids := make([]int64, len(rows))
	for i := range rows {
		ids[i] = rows[i].BatchId
	}
	levelMap, err := s.mockLevelsByBatchIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	// 批量加载当前会员在这些批次/等级下的会话状态，过滤掉已提交/已结束。
	type attemptAgg struct {
		ExamBatchId int64 `orm:"exam_batch_id"`
		MockLevelId int64 `orm:"mock_level_id"`
		MaxStatus   int   `orm:"max_status"`
	}
	var attRows []attemptAgg
	if err := dao.ExamAttempt.Ctx(ctx).
		Fields("exam_batch_id, mock_level_id, MAX(status) AS max_status").
		Where("member_id", memberID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		WhereIn("exam_batch_id", ids).
		Group("exam_batch_id, mock_level_id").
		Scan(&attRows); err != nil {
		return nil, err
	}
	statusMap := make(map[string]int, len(attRows))
	for _, a := range attRows {
		k := fmt.Sprintf("%d:%d", a.ExamBatchId, a.MockLevelId)
		statusMap[k] = a.MaxStatus
	}

	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {
		// 保留窗口状态（upcoming/open/closed），由前端决定展示策略。
		winStatus := examWindowStatus(now, r.ExamStartAt, r.ExamEndAt)
		// 已有会话且状态为「已交卷/已结束」则不再展示。
		key := fmt.Sprintf("%d:%d", r.BatchId, r.MockLevelId)
		if st, ok := statusMap[key]; ok && (st == consts.ExamAttemptSubmitted || st == consts.ExamAttemptEnded) {
			continue
		}

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
			MockLevelId:            r.MockLevelId,
			MockLevelIds:           lids,
			WindowStatus:           winStatus,
		})
	}
	return list, nil
}
