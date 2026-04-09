package exam

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/exampaper"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"
)

func dedupPositiveInt64(ids []int64) []int64 {
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

func parseBatchTime(s string) (*gtime.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	t := gtime.NewFromStr(s)
	if t == nil {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	return t, nil
}

func (s *sExam) ensureMockPaperInBatch(ctx context.Context, batchID, mockPaperID int64) error {
	if mockPaperID <= 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	c, err := dao.ExamBatchMockPaper.Ctx(ctx).Where("batch_id", batchID).Where("mock_examination_paper_id", mockPaperID).Count()
	if err != nil {
		return err
	}
	if c == 0 {
		return gerror.NewCode(consts.CodeExamBatchPaperNotInBatch, "")
	}
	return nil
}

func (s *sExam) ensureMockPapersImported(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if _, err := exampaper.ByMockID(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *sExam) ensureMembersExist(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	c, err := dao.SysMember.Ctx(ctx).WhereIn("id", ids).Where("delete_flag", consts.DeleteFlagNotDeleted).Count()
	if err != nil {
		return err
	}
	if c != len(ids) {
		return gerror.NewCode(consts.CodeUserNotFound, "")
	}
	return nil
}

func (s *sExam) examBatchByID(ctx context.Context, id int64) (examentity.ExamBatch, error) {
	var b examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&b); err != nil {
		return b, err
	}
	if b.Id == 0 {
		return b, gerror.NewCode(consts.CodeExamBatchNotFound, "")
	}
	return b, nil
}

func (s *sExam) loadMockPaperIDsForBatch(ctx context.Context, batchID int64) ([]int64, error) {
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

func (s *sExam) memberCountByBatchIDs(ctx context.Context, batchIDs []int64) (map[int64]int, error) {
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

func (s *sExam) mockPapersByBatchIDs(ctx context.Context, batchIDs []int64) (map[int64][]int64, error) {
	if len(batchIDs) == 0 {
		return map[int64][]int64{}, nil
	}
	var rows []examentity.ExamBatchMockPaper
	if err := dao.ExamBatchMockPaper.Ctx(ctx).WhereIn("batch_id", batchIDs).OrderAsc("mock_examination_paper_id").Scan(&rows); err != nil {
		return nil, err
	}
	m := make(map[int64][]int64)
	for _, r := range rows {
		m[r.BatchId] = append(m[r.BatchId], r.MockExaminationPaperId)
	}
	return m, nil
}

func (s *sExam) toBatchAdminItem(b examentity.ExamBatch, paperIDs []int64, memberCount int) bo.ExamBatchAdminItem {
	return bo.ExamBatchAdminItem{
		Id:                      b.Id,
		MockExaminationPaperIds: paperIDs,
		Title:                   b.Title,
		ExamStartAt:             b.ExamStartAt,
		ExamEndAt:               b.ExamEndAt,
		MemberCount:             memberCount,
		CreateTime:              b.CreateTime,
	}
}

// ExamBatchList 管理端批次分页；mockPaperID=0 时不按卷筛选。
func (s *sExam) ExamBatchList(ctx context.Context, mockPaperID int64, page, size int) (list []bo.ExamBatchAdminItem, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	q := dao.ExamBatch.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if mockPaperID > 0 {
		var ebpRows []struct {
			BatchId int64 `orm:"batch_id"`
		}
		if err := dao.ExamBatchMockPaper.Ctx(ctx).Fields("batch_id").Where("mock_examination_paper_id", mockPaperID).Scan(&ebpRows); err != nil {
			return nil, 0, err
		}
		if len(ebpRows) == 0 {
			return []bo.ExamBatchAdminItem{}, 0, nil
		}
		batchIDs := make([]int64, len(ebpRows))
		for i := range ebpRows {
			batchIDs[i] = ebpRows[i].BatchId
		}
		q = q.WhereIn("id", batchIDs)
	}
	total, err = q.Count()
	if err != nil {
		return nil, 0, err
	}
	var batches []examentity.ExamBatch
	if err = q.Page(page, size).OrderDesc("id").Scan(&batches); err != nil {
		return nil, 0, err
	}
	ids := make([]int64, len(batches))
	for i := range batches {
		ids[i] = batches[i].Id
	}
	paperMap, err := s.mockPapersByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	cntMap, err := s.memberCountByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	list = make([]bo.ExamBatchAdminItem, 0, len(batches))
	for _, b := range batches {
		pids := paperMap[b.Id]
		if pids == nil {
			pids = []int64{}
		}
		list = append(list, s.toBatchAdminItem(b, pids, cntMap[b.Id]))
	}
	return list, total, nil
}

// ExamBatchDetail 批次详情（含 Mock 卷 id 列表与学员数）。
func (s *sExam) ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error) {
	b, err := s.examBatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	pids, err := s.loadMockPaperIDsForBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	cntMap, err := s.memberCountByBatchIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	item := s.toBatchAdminItem(b, pids, cntMap[id])
	return &item, nil
}

// ExamBatchCreate 创建批次并写入多选 Mock 卷（须已导入 exam_paper）。
func (s *sExam) ExamBatchCreate(ctx context.Context, title, examStartAt, examEndAt string, mockPaperIDs []int64, creator string) (int64, error) {
	paperIDs := dedupPositiveInt64(mockPaperIDs)
	if len(paperIDs) == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	if err := s.ensureMockPapersImported(ctx, paperIDs); err != nil {
		return 0, err
	}
	st, err := parseBatchTime(examStartAt)
	if err != nil {
		return 0, err
	}
	en, err := parseBatchTime(examEndAt)
	if err != nil {
		return 0, err
	}
	if !en.After(st) {
		return 0, gerror.NewCode(consts.CodeExamBatchTimeInvalid, "")
	}
	var newID int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		row := examentity.ExamBatch{
			Title:       strings.TrimSpace(title),
			ExamStartAt: st,
			ExamEndAt:   en,
			Creator:     creator,
			Updater:     creator,
			CreateTime:  gtime.Now(),
			UpdateTime:  gtime.Now(),
			DeleteFlag:  consts.DeleteFlagNotDeleted,
		}
		pid, err := tx.Model(dao.ExamBatch.Table()).Ctx(ctx).InsertAndGetId(row)
		if err != nil {
			return err
		}
		newID = pid
		batch := make([]gdb.Map, 0, len(paperIDs))
		for _, mid := range paperIDs {
			batch = append(batch, gdb.Map{
				dao.ExamBatchMockPaper.Columns().BatchId:                newID,
				dao.ExamBatchMockPaper.Columns().MockExaminationPaperId: mid,
			})
		}
		if _, err := tx.Model(dao.ExamBatchMockPaper.Table()).Ctx(ctx).Data(batch).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (s *sExam) batchHasMembersOutsidePaperSet(ctx context.Context, batchID int64, paperIDs []int64) (bool, error) {
	if len(paperIDs) == 0 {
		n, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID).Count()
		return n > 0, err
	}
	n, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID).WhereNotIn("mock_examination_paper_id", paperIDs).Count()
	return n > 0, err
}

// ExamBatchUpdate 更新批次时间与 Mock 卷多选（全量替换卷关联）；若移除的卷上仍有成员绑定则失败。
func (s *sExam) ExamBatchUpdate(ctx context.Context, id int64, title, examStartAt, examEndAt string, mockPaperIDs []int64, updater string) error {
	if _, err := s.examBatchByID(ctx, id); err != nil {
		return err
	}
	st, err := parseBatchTime(examStartAt)
	if err != nil {
		return err
	}
	en, err := parseBatchTime(examEndAt)
	if err != nil {
		return err
	}
	if !en.After(st) {
		return gerror.NewCode(consts.CodeExamBatchTimeInvalid, "")
	}
	paperIDs := dedupPositiveInt64(mockPaperIDs)
	if len(paperIDs) == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	if err := s.ensureMockPapersImported(ctx, paperIDs); err != nil {
		return err
	}
	hasOrphan, err := s.batchHasMembersOutsidePaperSet(ctx, id, paperIDs)
	if err != nil {
		return err
	}
	if hasOrphan {
		return gerror.NewCode(consts.CodeExamBatchPaperHasMembers, "")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.ExamBatch.Table()).Ctx(ctx).Where("id", id).Data(g.Map{
			dao.ExamBatch.Columns().Title:       strings.TrimSpace(title),
			dao.ExamBatch.Columns().ExamStartAt: st,
			dao.ExamBatch.Columns().ExamEndAt:   en,
			dao.ExamBatch.Columns().Updater:     updater,
			dao.ExamBatch.Columns().UpdateTime:  gtime.Now(),
		}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.ExamBatchMockPaper.Table()).Ctx(ctx).Where("batch_id", id).Delete(); err != nil {
			return err
		}
		batch := make([]gdb.Map, 0, len(paperIDs))
		for _, mid := range paperIDs {
			batch = append(batch, gdb.Map{
				dao.ExamBatchMockPaper.Columns().BatchId:                id,
				dao.ExamBatchMockPaper.Columns().MockExaminationPaperId: mid,
			})
		}
		if _, err := tx.Model(dao.ExamBatchMockPaper.Table()).Ctx(ctx).Data(batch).Insert(); err != nil {
			return err
		}
		return nil
	})
}

// ExamBatchDelete 逻辑删除批次（不删学员关联与卷行，便于审计；列表已过滤）。
func (s *sExam) ExamBatchDelete(ctx context.Context, id int64) error {
	if _, err := s.examBatchByID(ctx, id); err != nil {
		return err
	}
	_, err := dao.ExamBatch.Ctx(ctx).Where("id", id).Data(g.Map{
		dao.ExamBatch.Columns().DeleteFlag: consts.DeleteFlagDeleted,
		dao.ExamBatch.Columns().UpdateTime: gtime.Now(),
	}).Update()
	return err
}

// ExamBatchMembersImport 导入学员（指定批次内 Mock 卷）；已存在 (batch,member,paper) 则跳过。
func (s *sExam) ExamBatchMembersImport(ctx context.Context, batchID int64, mockPaperID int64, memberIDs []int64, creator string) (inserted int, err error) {
	if _, err := s.examBatchByID(ctx, batchID); err != nil {
		return 0, err
	}
	if err := s.ensureMockPaperInBatch(ctx, batchID, mockPaperID); err != nil {
		return 0, err
	}
	ids := dedupPositiveInt64(memberIDs)
	if len(ids) == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	if err := s.ensureMembersExist(ctx, ids); err != nil {
		return 0, err
	}
	var existing []struct {
		MemberId int64 `orm:"member_id"`
	}
	if err := dao.ExamBatchMember.Ctx(ctx).Fields("member_id").Where("batch_id", batchID).Where("mock_examination_paper_id", mockPaperID).WhereIn("member_id", ids).Scan(&existing); err != nil {
		return 0, err
	}
	have := make(map[int64]struct{}, len(existing))
	for _, e := range existing {
		have[e.MemberId] = struct{}{}
	}
	toAdd := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := have[id]; ok {
			continue
		}
		toAdd = append(toAdd, id)
	}
	if len(toAdd) == 0 {
		return 0, nil
	}
	now := gtime.Now()
	rows := make([]gdb.Map, 0, len(toAdd))
	for _, mid := range toAdd {
		rows = append(rows, gdb.Map{
			dao.ExamBatchMember.Columns().BatchId:                batchID,
			dao.ExamBatchMember.Columns().MemberId:               mid,
			dao.ExamBatchMember.Columns().MockExaminationPaperId: mockPaperID,
			dao.ExamBatchMember.Columns().Creator:                creator,
			dao.ExamBatchMember.Columns().CreateTime:             now,
		})
	}
	_, err = dao.ExamBatchMember.Ctx(ctx).Data(rows).Insert()
	if err != nil {
		return 0, err
	}
	return len(toAdd), nil
}

// ExamBatchMemberList 批次内学员分页。
func (s *sExam) ExamBatchMemberList(ctx context.Context, batchID int64, page, size int) (list []bo.ExamBatchMemberAdminRow, total int, err error) {
	if _, err := s.examBatchByID(ctx, batchID); err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	q := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID)
	total, err = q.Count()
	if err != nil {
		return nil, 0, err
	}
	var links []examentity.ExamBatchMember
	if err = q.Page(page, size).OrderDesc("create_time").Scan(&links); err != nil {
		return nil, 0, err
	}
	if len(links) == 0 {
		return []bo.ExamBatchMemberAdminRow{}, total, nil
	}
	memberIDs := make([]int64, len(links))
	type linkKey struct {
		mid int64
		pid int64
	}
	linkTime := make(map[linkKey]*gtime.Time, len(links))
	for i, l := range links {
		memberIDs[i] = l.MemberId
		linkTime[linkKey{l.MemberId, l.MockExaminationPaperId}] = l.CreateTime
	}
	var users []sysentity.SysMember
	if err = dao.SysMember.Ctx(ctx).WhereIn("id", memberIDs).Scan(&users); err != nil {
		return nil, 0, err
	}
	um := make(map[int64]sysentity.SysMember, len(users))
	for _, u := range users {
		um[u.Id] = u
	}
	list = make([]bo.ExamBatchMemberAdminRow, 0, len(links))
	for _, l := range links {
		u := um[l.MemberId]
		list = append(list, bo.ExamBatchMemberAdminRow{
			MemberId:               l.MemberId,
			MockExaminationPaperId: l.MockExaminationPaperId,
			Username:               u.Username,
			Nickname:               u.Nickname,
			ImportTime:             linkTime[linkKey{l.MemberId, l.MockExaminationPaperId}],
		})
	}
	return list, total, nil
}

// ExamBatchMembersRemove 从批次移除学员（指定 mock_examination_paper_id）。
func (s *sExam) ExamBatchMembersRemove(ctx context.Context, batchID int64, mockPaperID int64, memberIDs []int64) (removed int, err error) {
	if _, err := s.examBatchByID(ctx, batchID); err != nil {
		return 0, err
	}
	if mockPaperID <= 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	ids := dedupPositiveInt64(memberIDs)
	if len(ids) == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	r, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID).Where("mock_examination_paper_id", mockPaperID).WhereIn("member_id", ids).Delete()
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// ExamBatchMemberDetail 指定批次、学员、Mock 卷的成员绑定行。
func (s *sExam) ExamBatchMemberDetail(ctx context.Context, batchID int64, userID int64, mockPaperID int64) (*examentity.ExamBatchMember, error) {
	var row examentity.ExamBatchMember
	err := dao.ExamBatchMember.Ctx(ctx).
		Where("batch_id", batchID).
		Where("member_id", userID).
		Where("mock_examination_paper_id", mockPaperID).
		Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.BatchId == 0 {
		return nil, gerror.NewCode(consts.CodeExamBatchMemberNotFound, "")
	}
	return &row, nil
}
