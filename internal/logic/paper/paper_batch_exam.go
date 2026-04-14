package paper

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/service/batch"
)

func examBatchParseTime(str string) *gtime.Time {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	return gtime.NewFromStr(str)
}

func dedupePositiveMockPaperIDs(ids []int64) []int64 {
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

func (s *sPaper) MyExamBatches(ctx context.Context, memberID int64) ([]bo.MyExamBatchItem, error) {
	return batch.Batch().MyExamBatches(ctx, memberID)
}

func (s *sPaper) ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error) {
	return batch.Batch().ExamBatchDetail(ctx, id)
}

func (s *sPaper) ExamBatchDelete(ctx context.Context, id int64) error {
	return batch.Batch().ExamBatchDelete(ctx, id, "")
}

func (s *sPaper) ExamBatchMembersRemove(ctx context.Context, batchID int64, mockExaminationPaperId int64, memberIDs []int64) (int, error) {
	return batch.Batch().ExamBatchMembersRemove(ctx, batchID, mockExaminationPaperId, memberIDs)
}

func (s *sPaper) ExamBatchMembersImport(ctx context.Context, batchID int64, mockExaminationPaperId int64, memberIDs []int64, creator string) (int, error) {
	if _, err := batch.Batch().GetBatchByID(ctx, batchID); err != nil {
		return 0, err
	}
	inserted := 0
	err := dao.ExamBatchMember.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, mid := range memberIDs {
			if mid <= 0 {
				continue
			}
			cnt, _ := dao.ExamBatchMember.Ctx(ctx).TX(tx).
				Where("batch_id", batchID).
				Where("mock_examination_paper_id", mockExaminationPaperId).
				Where("member_id", mid).
				Count()
			if cnt > 0 {
				continue
			}
			if _, err := dao.ExamBatchMember.Ctx(ctx).TX(tx).Insert(g.Map{
				"batch_id":                  batchID,
				"mock_examination_paper_id": mockExaminationPaperId,
				"member_id":                 mid,
				"creator":                   creator,
			}); err != nil {
				return err
			}
			inserted++
		}
		return nil
	})
	return inserted, err
}

func (s *sPaper) ExamBatchMemberList(ctx context.Context, batchID int64, page int, size int) ([]bo.ExamBatchMemberAdminRow, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := g.DB().Ctx(ctx).Model("exam_batch_member ebm").
		LeftJoin("sys_user u", "u.id = ebm.member_id").
		Fields("ebm.member_id, ebm.mock_examination_paper_id, u.username, u.nickname, ebm.create_time AS import_time").
		Where("ebm.batch_id", batchID)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []bo.ExamBatchMemberAdminRow{}, 0, nil
	}
	type row struct {
		MemberId               int64       `json:"member_id"`
		MockExaminationPaperId int64       `json:"mock_examination_paper_id"`
		Username               string      `json:"username"`
		Nickname               string      `json:"nickname"`
		ImportTime             *gtime.Time `json:"import_time"`
	}
	var raw []row
	if err := m.Page(page, size).OrderAsc("ebm.member_id").OrderAsc("ebm.mock_examination_paper_id").Scan(&raw); err != nil {
		return nil, 0, err
	}
	out := make([]bo.ExamBatchMemberAdminRow, len(raw))
	for i, r := range raw {
		out[i] = bo.ExamBatchMemberAdminRow{
			MemberId:               r.MemberId,
			MockExaminationPaperId: r.MockExaminationPaperId,
			Username:               r.Username,
			Nickname:               r.Nickname,
			ImportTime:             r.ImportTime,
		}
	}
	return out, total, nil
}

func (s *sPaper) ExamBatchList(ctx context.Context, mockExaminationPaperId int64, page int, size int) ([]bo.ExamBatchAdminItem, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	m := dao.ExamBatch.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if mockExaminationPaperId > 0 {
		m = m.Where(fmt.Sprintf("id IN (SELECT batch_id FROM %s WHERE mock_examination_paper_id=?)",
			dao.ExamBatchMockPaper.Table()), mockExaminationPaperId)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []bo.ExamBatchAdminItem{}, 0, nil
	}
	var list []examentity.ExamBatch
	if err := m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, err
	}
	batchIDs := make([]int64, len(list))
	for i, b := range list {
		batchIDs[i] = b.Id
	}
	cntMap, err := examBatchMemberCountByBatchIDs(ctx, batchIDs)
	if err != nil {
		return nil, 0, err
	}
	out := make([]bo.ExamBatchAdminItem, 0, len(list))
	for _, b := range list {
		pids, err := examBatchLoadMockPaperIDs(ctx, b.Id)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, bo.ExamBatchAdminItem{
			Id:                      b.Id,
			MockExaminationPaperIds: pids,
			Title:                   b.Title,
			ExamStartAt:             b.ExamStartAt,
			ExamEndAt:               b.ExamEndAt,
			MemberCount:             cntMap[b.Id],
			CreateTime:              b.CreateTime,
		})
	}
	return out, total, nil
}

func examBatchMemberCountByBatchIDs(ctx context.Context, batchIDs []int64) (map[int64]int, error) {
	if len(batchIDs) == 0 {
		return map[int64]int{}, nil
	}
	type cntRow struct {
		BatchId int64 `orm:"batch_id"`
		Cnt     int   `orm:"cnt"`
	}
	var rows []cntRow
	err := g.DB().Ctx(ctx).Model(dao.ExamBatchMember.Table()).
		Fields("batch_id, COUNT(*) AS cnt").
		WhereIn("batch_id", batchIDs).
		Group("batch_id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	m := make(map[int64]int, len(rows))
	for _, r := range rows {
		m[r.BatchId] = r.Cnt
	}
	return m, nil
}

func examBatchLoadMockPaperIDs(ctx context.Context, batchID int64) ([]int64, error) {
	var rows []examentity.ExamBatchMockPaper
	if err := dao.ExamBatchMockPaper.Ctx(ctx).Where("batch_id", batchID).OrderAsc("mock_examination_paper_id").Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.MockExaminationPaperId)
	}
	return out, nil
}

func (s *sPaper) ExamBatchCreate(ctx context.Context, title string, examStartAt string, examEndAt string, mockExaminationPaperIds []int64, creator string) (int64, error) {
	st := examBatchParseTime(examStartAt)
	et := examBatchParseTime(examEndAt)
	if st == nil || et == nil {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	mockIDs := dedupePositiveMockPaperIDs(mockExaminationPaperIds)
	var batchID int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.ExamBatch.Columns()
		r, err := tx.Model(dao.ExamBatch.Table()).Ctx(ctx).Insert(g.Map{
			cols.Title:       title,
			cols.ExamStartAt: st,
			cols.ExamEndAt:   et,
			cols.Creator:     creator,
			cols.DeleteFlag:  consts.DeleteFlagNotDeleted,
		})
		if err != nil {
			return err
		}
		id, err := r.LastInsertId()
		if err != nil {
			return err
		}
		batchID = id
		for _, mid := range mockIDs {
			if _, err := tx.Model(dao.ExamBatchMockPaper.Table()).Ctx(ctx).Insert(g.Map{
				"batch_id":                  batchID,
				"mock_examination_paper_id": mid,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	return batchID, err
}

func (s *sPaper) ExamBatchUpdate(ctx context.Context, id int64, title string, examStartAt string, examEndAt string, mockExaminationPaperIds []int64, updater string) error {
	if err := batch.Batch().ExamBatchUpdate(ctx, id, title, examStartAt, examEndAt, updater); err != nil {
		return err
	}
	mockIDs := dedupePositiveMockPaperIDs(mockExaminationPaperIds)
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.ExamBatchMockPaper.Table()).Ctx(ctx).Where("batch_id", id).Delete(); err != nil {
			return err
		}
		for _, mid := range mockIDs {
			if _, err := tx.Model(dao.ExamBatchMockPaper.Table()).Ctx(ctx).Insert(g.Map{
				"batch_id":                  id,
				"mock_examination_paper_id": mid,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
