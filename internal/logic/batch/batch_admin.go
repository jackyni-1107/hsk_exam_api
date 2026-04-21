package batch

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/utility"
)

// ExamBatchList 分页查询考试批次列表
func (s *sBatch) ExamBatchList(ctx context.Context, examPaperID int64, page int, size int, key string) (list []bo.ExamBatchAdminItem, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	q := dao.ExamBatch.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if examPaperID > 0 {
		var ebpRows []struct {
			BatchId int64 `orm:"batch_id"`
		}
		if err := dao.ExamBatchPaper.Ctx(ctx).Fields("batch_id").Where("exam_paper_id", examPaperID).Scan(&ebpRows); err != nil {
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
	paperMap, err := examPapersByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	cntMap, err := memberCountByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	list = make([]bo.ExamBatchAdminItem, 0, len(batches))
	for _, b := range batches {
		pids := paperMap[b.Id]
		if pids == nil {
			pids = []int64{}
		}
		list = append(list, toBatchAdminItem(b, pids, cntMap[b.Id]))
	}
	return list, total, nil
}

// ExamBatchDetail 批次详情（含 exam_paper.id 列表与学员数）。
func (s *sBatch) ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error) {
	b, err := examBatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	pids, err := loadExamPaperIDsForBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	cntMap, err := memberCountByBatchIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	item := toBatchAdminItem(b, pids, cntMap[id])
	return &item, nil
}

// ExamBatchCreate 创建考试批次
func (s *sBatch) ExamBatchCreate(ctx context.Context, title, examStartAt, examEndAt string, examPaperIDs []int64, creator string) (int64, error) {
	paperIDs := s.dedupIDs(examPaperIDs)
	if len(paperIDs) == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	if err := ensureExamPapersExist(ctx, paperIDs); err != nil {
		return 0, err
	}
	epMap, err := examPaperMapByIDs(ctx, paperIDs)
	if err != nil {
		return 0, err
	}
	st := s.parseTime(examStartAt)

	en := s.parseTime(examEndAt)
	if !en.After(st) {
		return 0, gerror.NewCode(consts.CodeExamBatchTimeInvalid)
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
		for _, epID := range paperIDs {
			p := epMap[epID]
			batch = append(batch, gdb.Map{
				dao.ExamBatchPaper.Columns().BatchId:                newID,
				dao.ExamBatchPaper.Columns().ExamPaperId:            epID,
				dao.ExamBatchPaper.Columns().MockExaminationPaperId: p.MockExaminationPaperId,
			})
		}
		if _, err := tx.Model(dao.ExamBatchPaper.Table()).Ctx(ctx).Data(batch).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if b, err := examBatchByID(ctx, newID); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.ExamBatch.Table(), newID, nil, &b)
	}
	if pids, err := loadExamPaperIDsForBatch(ctx, newID); err == nil {
		afterStr := utility.JoinSortedInt64IDs(pids)
		auditutil.RecordMapDiff(ctx, dao.ExamBatch.Table(), newID,
			map[string]interface{}{"exam_paper_ids": ""},
			map[string]interface{}{"exam_paper_ids": afterStr})
	}
	return newID, nil
}

// ExamBatchUpdate 更新考试批次
func (s *sBatch) ExamBatchUpdate(ctx context.Context, id int64, title, examStartAt, examEndAt string, examPaperIDs []int64, updater string) error {
	if _, err := examBatchByID(ctx, id); err != nil {
		return err
	}
	st := s.parseTime(examStartAt)
	en := s.parseTime(examEndAt)
	if !en.After(st) {
		return gerror.NewCode(consts.CodeExamBatchTimeInvalid)
	}
	paperIDs := s.dedupIDs(examPaperIDs)
	if len(paperIDs) == 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	if err := ensureExamPapersExist(ctx, paperIDs); err != nil {
		return err
	}
	epMap, err := examPaperMapByIDs(ctx, paperIDs)
	if err != nil {
		return err
	}
	hasOrphan, err := batchHasMembersOutsidePaperSet(ctx, id, paperIDs)
	if err != nil {
		return err
	}
	if hasOrphan {
		return gerror.NewCode(consts.CodeExamBatchPaperHasMembers)
	}
	beforeBatch, err := examBatchByID(ctx, id)
	if err != nil {
		return err
	}
	beforePapers, err := loadExamPaperIDsForBatch(ctx, id)
	if err != nil {
		return err
	}
	beforePaperStr := utility.JoinSortedInt64IDs(beforePapers)
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.ExamBatch.Table()).Ctx(ctx).Where("id", id).Data(g.Map{
			dao.ExamBatch.Columns().Title:       strings.TrimSpace(title),
			dao.ExamBatch.Columns().ExamStartAt: st,
			dao.ExamBatch.Columns().ExamEndAt:   en,
			dao.ExamBatch.Columns().Updater:     updater,
			dao.ExamBatch.Columns().UpdateTime:  gtime.Now(),
		}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.ExamBatchPaper.Table()).Ctx(ctx).Where("batch_id", id).Delete(); err != nil {
			return err
		}
		batch := make([]gdb.Map, 0, len(paperIDs))
		for _, epID := range paperIDs {
			p := epMap[epID]
			batch = append(batch, gdb.Map{
				dao.ExamBatchPaper.Columns().BatchId:                id,
				dao.ExamBatchPaper.Columns().ExamPaperId:            epID,
				dao.ExamBatchPaper.Columns().MockExaminationPaperId: p.MockExaminationPaperId,
			})
		}
		if _, err := tx.Model(dao.ExamBatchPaper.Table()).Ctx(ctx).Data(batch).Insert(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if afterBatch, err := examBatchByID(ctx, id); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.ExamBatch.Table(), id, &beforeBatch, &afterBatch)
	}
	if afterPapers, err := loadExamPaperIDsForBatch(ctx, id); err == nil {
		afterPaperStr := utility.JoinSortedInt64IDs(afterPapers)
		auditutil.RecordMapDiff(ctx, dao.ExamBatch.Table(), id,
			map[string]interface{}{"exam_paper_ids": beforePaperStr},
			map[string]interface{}{"exam_paper_ids": afterPaperStr})
	}
	return nil
}

// ExamBatchDelete 删除考试批次
func (s *sBatch) ExamBatchDelete(ctx context.Context, id int64) error {
	before, err := examBatchByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = dao.ExamBatch.Ctx(ctx).Where("id", id).Data(g.Map{
		dao.ExamBatch.Columns().DeleteFlag: consts.DeleteFlagDeleted,
		dao.ExamBatch.Columns().UpdateTime: gtime.Now(),
	}).Update()
	if err != nil {
		return err
	}
	var after examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.ExamBatch.Table(), id, &before, &after)
	}
	return nil
}

// ExamBatchMembersAdd 批量向指定批次和试卷（exam_paper.id）添加学员
func (s *sBatch) ExamBatchMembersAdd(ctx context.Context, batchID int64, examPaperID int64, memberIDs []int64, creator string) (inserted int, err error) {
	if _, err := examBatchByID(ctx, batchID); err != nil {
		return 0, err
	}
	if err := ensureExamPaperInBatch(ctx, batchID, examPaperID); err != nil {
		return 0, err
	}
	ids := s.dedupIDs(memberIDs)
	if len(ids) == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	if err := ensureMembersExist(ctx, ids); err != nil {
		return 0, err
	}
	var epPaper examentity.ExamPaper
	if err := dao.ExamPaper.Ctx(ctx).
		Where("id", examPaperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&epPaper); err != nil {
		return 0, err
	}
	if epPaper.Id == 0 || epPaper.MockExaminationPaperId <= 0 {
		return 0, gerror.NewCode(consts.CodeExamPaperNotFound)
	}
	var existing []struct {
		MemberId int64 `orm:"member_id"`
	}
	if err := dao.ExamBatchMember.Ctx(ctx).Fields("member_id").Where("batch_id", batchID).Where("exam_paper_id", examPaperID).WhereIn("member_id", ids).Scan(&existing); err != nil {
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
			dao.ExamBatchMember.Columns().ExamPaperId:            examPaperID,
			dao.ExamBatchMember.Columns().MockExaminationPaperId: epPaper.MockExaminationPaperId,
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

// ExamBatchMembersRemove 从批次中移除学员
func (s *sBatch) ExamBatchMembersRemove(ctx context.Context, batchID int64, examPaperID int64, memberIDs []int64) (int, error) {
	ids := s.dedupIDs(memberIDs)
	if len(ids) == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchMemberNotFound)
	}

	r, err := dao.ExamBatchMember.Ctx(ctx).
		Where("batch_id", batchID).
		Where("exam_paper_id", examPaperID).
		WhereIn("member_id", ids).Delete()

	if err != nil {
		return 0, gerror.NewCode(consts.CodeExamBatchMemberNotFound)
	}
	n, _ := r.RowsAffected()
	return int(n), nil
}

// ExamBatchMemberList 查询批次内的成员列表（关联系统用户表）
func (s *sBatch) ExamBatchMemberList(ctx context.Context, batchID int64, page, size int) (list []bo.ExamBatchMemberAdminRow, total int, err error) {
	if _, err := examBatchByID(ctx, batchID); err != nil {
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
		linkTime[linkKey{l.MemberId, l.ExamPaperId}] = l.CreateTime
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
			MemberId:    l.MemberId,
			ExamPaperId: l.ExamPaperId,
			Username:    u.Username,
			Nickname:    u.Nickname,
			ImportTime:  linkTime[linkKey{l.MemberId, l.ExamPaperId}],
		})
	}
	return list, total, nil
}

func examPapersByBatchIDs(ctx context.Context, batchIDs []int64) (map[int64][]int64, error) {
	if len(batchIDs) == 0 {
		return map[int64][]int64{}, nil
	}
	var rows []examentity.ExamBatchPaper
	if err := dao.ExamBatchPaper.Ctx(ctx).WhereIn("batch_id", batchIDs).OrderAsc("exam_paper_id").Scan(&rows); err != nil {
		return nil, err
	}
	m := make(map[int64][]int64)
	for _, r := range rows {
		m[r.BatchId] = append(m[r.BatchId], r.ExamPaperId)
	}
	return m, nil
}

func batchHasMembersOutsidePaperSet(ctx context.Context, batchID int64, paperIDs []int64) (bool, error) {
	if len(paperIDs) == 0 {
		n, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID).Count()
		return n > 0, err
	}
	n, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID).WhereNotIn("exam_paper_id", paperIDs).Count()
	return n > 0, err
}

func toBatchAdminItem(b examentity.ExamBatch, paperIDs []int64, memberCount int) bo.ExamBatchAdminItem {
	return bo.ExamBatchAdminItem{
		Id:           b.Id,
		ExamPaperIds: paperIDs,
		Title:        b.Title,
		ExamStartAt:  b.ExamStartAt,
		ExamEndAt:    b.ExamEndAt,
		MemberCount:  memberCount,
		CreateTime:   b.CreateTime,
	}
}

func memberCountByBatchIDs(ctx context.Context, batchIDs []int64) (map[int64]int, error) {
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

func loadExamPaperIDsForBatch(ctx context.Context, batchID int64) ([]int64, error) {
	var rows []examentity.ExamBatchPaper
	if err := dao.ExamBatchPaper.Ctx(ctx).Where("batch_id", batchID).OrderAsc("exam_paper_id").Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.ExamPaperId)
	}
	return out, nil
}

func examBatchByID(ctx context.Context, id int64) (examentity.ExamBatch, error) {
	var b examentity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&b); err != nil {
		return b, err
	}
	if b.Id == 0 {
		return b, gerror.NewCode(consts.CodeExamBatchNotFound)
	}
	return b, nil
}

func ensureExamPaperInBatch(ctx context.Context, batchID, examPaperID int64) error {
	if examPaperID <= 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	c, err := dao.ExamBatchPaper.Ctx(ctx).Where("batch_id", batchID).Where("exam_paper_id", examPaperID).Count()
	if err != nil {
		return err
	}
	if c == 0 {
		return gerror.NewCode(consts.CodeExamBatchPaperNotInBatch)
	}
	return nil
}

func ensureMembersExist(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	c, err := dao.SysMember.Ctx(ctx).WhereIn("id", ids).Where("delete_flag", consts.DeleteFlagNotDeleted).Count()
	if err != nil {
		return err
	}
	if c != len(ids) {
		return gerror.NewCode(consts.CodeUserNotFound)
	}
	return nil
}

func ensureExamPapersExist(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if id <= 0 {
			return gerror.NewCode(consts.CodeInvalidParams)
		}
		n, err := dao.ExamPaper.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Count()
		if err != nil {
			return err
		}
		if n == 0 {
			return gerror.NewCode(consts.CodeExamPaperNotFound)
		}
	}
	return nil
}

// examPaperMapByIDs 批量加载 exam_paper；用于写入 exam_batch_paper / exam_batch_member 时绑定 mock_examination_paper_id。
func examPaperMapByIDs(ctx context.Context, ids []int64) (map[int64]examentity.ExamPaper, error) {
	if len(ids) == 0 {
		return map[int64]examentity.ExamPaper{}, nil
	}
	var list []examentity.ExamPaper
	if err := dao.ExamPaper.Ctx(ctx).
		WhereIn("id", ids).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&list); err != nil {
		return nil, err
	}
	m := make(map[int64]examentity.ExamPaper, len(list))
	for _, p := range list {
		m[p.Id] = p
		if p.MockExaminationPaperId <= 0 {
			return nil, gerror.NewCode(consts.CodeExamPaperNotFound)
		}
	}
	if len(m) != len(ids) {
		return nil, gerror.NewCode(consts.CodeExamPaperNotFound)
	}
	return m, nil
}
