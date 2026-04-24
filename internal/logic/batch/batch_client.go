package batch

import (
	"context"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	"fmt"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MyExamBatches 学员查询自己的批次
func (s *sBatch) MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error) {
	type rowT struct {
		BatchId                int64       `orm:"batch_id"`
		Title                  string      `orm:"title"`
		ExamPaperId            int64       `orm:"exam_paper_id"`
		MockExaminationPaperId int64       `orm:"mock_src_id"`
		PaperTitle             string      `orm:"paper_title"`
		ExamStartAt            *gtime.Time `orm:"exam_start_at"`
		ExamEndAt              *gtime.Time `orm:"exam_end_at"`
		BatchKind              int         `orm:"batch_kind"`
		AllowMultipleAttempts  int         `orm:"allow_multiple_attempts"`
	}

	var practiceRows []rowT
	if err = g.DB().Ctx(ctx).Model("exam_batch_paper ebp").
		InnerJoin("exam_batch eb", "eb.id = ebp.batch_id").
		InnerJoin("exam_paper ep", "ep.id = ebp.exam_paper_id").
		Fields("eb.id AS batch_id, eb.title, ebp.exam_paper_id, ep.mock_examination_paper_id AS mock_src_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at, eb.batch_kind, eb.allow_multiple_attempts").
		Where("eb.batch_kind", consts.ExamBatchKindPractice).
		Where("eb.exam_start_at <= ?", gtime.Now()).
		Where("eb.exam_end_at >= ?", gtime.Now()).
		Where(g.Map{
			"eb.delete_flag": consts.DeleteFlagNotDeleted,
			"ep.delete_flag": consts.DeleteFlagNotDeleted,
		}).
		OrderDesc("eb.id").
		Scan(&practiceRows); err != nil {
		return nil, err
	}

	var memberRows []rowT
	if err = g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		LeftJoin("exam_batch eb", "eb.id = ebm.batch_id").
		LeftJoin("exam_paper ep", "ep.id = ebm.exam_paper_id").
		Fields("eb.id AS batch_id, eb.title, ebm.exam_paper_id, ep.mock_examination_paper_id AS mock_src_id, ep.title AS paper_title, eb.exam_start_at, eb.exam_end_at, eb.batch_kind, eb.allow_multiple_attempts").
		Where("ebm.member_id = ?", memberID).
		Where("eb.batch_kind", consts.ExamBatchKindFormal).
		Where("eb.exam_start_at <= ?", gtime.Now()).
		Where("eb.exam_end_at >= ?", gtime.Now()).
		Where(g.Map{
			"eb.delete_flag": consts.DeleteFlagNotDeleted,
			"ep.delete_flag": consts.DeleteFlagNotDeleted,
		}).
		OrderDesc("eb.id").
		Scan(&memberRows); err != nil {
		return nil, err
	}

	rows := make([]rowT, 0, len(practiceRows)+len(memberRows))
	rows = append(rows, practiceRows...)
	rows = append(rows, memberRows...)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].BatchId != rows[j].BatchId {
			return rows[i].BatchId > rows[j].BatchId
		}
		return rows[i].ExamPaperId < rows[j].ExamPaperId
	})

	if len(rows) == 0 {
		return []bo.MyExamBatchItem{}, nil
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
		OrderDesc("id").
		Scan(&attempts)
	if err != nil {
		return nil, err
	}

	type attInfo struct {
		Id     int64
		Status int
	}
	attemptsByKey := make(map[string][]attInfo)
	for _, a := range attempts {
		key := fmt.Sprintf("%d_%d", a.ExamBatchId, a.ExamPaperId)
		attemptsByKey[key] = append(attemptsByKey[key], attInfo{Id: a.Id, Status: a.Status})
	}

	pickAttemptID := func(key string, allowMulti bool) (int64, bool) {
		arr := attemptsByKey[key]
		if len(arr) == 0 {
			return 0, false
		}
		if !allowMulti {
			return arr[0].Id, true
		}
		for _, a := range arr {
			if a.Status == consts.ExamAttemptNotStarted || a.Status == consts.ExamAttemptInProgress {
				return a.Id, true
			}
		}
		return 0, false
	}

	allSubmittedOrEnded := func(arr []attInfo) bool {
		if len(arr) == 0 {
			return false
		}
		for _, a := range arr {
			if a.Status != consts.ExamAttemptSubmitted && a.Status != consts.ExamAttemptEnded {
				return false
			}
		}
		return true
	}

	now := gtime.Now()
	list = make([]bo.MyExamBatchItem, 0, len(rows))
	for _, r := range rows {
		key := fmt.Sprintf("%d_%d", r.BatchId, r.ExamPaperId)
		arr := attemptsByKey[key]
		allowMulti := r.AllowMultipleAttempts != 0

		if !allowMulti {
			if len(arr) > 0 && allSubmittedOrEnded(arr) {
				continue
			}
		}

		attID, hasAttempt := pickAttemptID(key, allowMulti)

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
			AttemptId:              attID,
			WindowStatus:           winStatus,
			BatchKind:              r.BatchKind,
			AllowMultipleAttempts:  r.AllowMultipleAttempts,
		})
	}

	return list, nil
}
