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
	"exam/internal/model/entity"
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

func (s *sExam) ensureMockLevelsExist(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	c, err := dao.MockLevels.Ctx(ctx).WhereIn("id", ids).Where("delete_flag", consts.DeleteFlagNotDeleted).Count()
	if err != nil {
		return err
	}
	if c != len(ids) {
		return gerror.NewCode(consts.CodeMockLevelNotFound, "")
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

func (s *sExam) examBatchByID(ctx context.Context, id int64) (entity.ExamBatch, error) {
	var b entity.ExamBatch
	if err := dao.ExamBatch.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&b); err != nil {
		return b, err
	}
	if b.Id == 0 {
		return b, gerror.NewCode(consts.CodeExamBatchNotFound, "")
	}
	return b, nil
}

func (s *sExam) loadMockLevelIDsForBatch(ctx context.Context, batchID int64) ([]int64, error) {
	var rows []entity.ExamBatchMockLevel
	if err := dao.ExamBatchMockLevel.Ctx(ctx).Where("batch_id", batchID).OrderAsc("mock_level_id").Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]int64, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.MockLevelId)
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

func (s *sExam) mockLevelsByBatchIDs(ctx context.Context, batchIDs []int64) (map[int64][]int64, error) {
	if len(batchIDs) == 0 {
		return map[int64][]int64{}, nil
	}
	var rows []entity.ExamBatchMockLevel
	if err := dao.ExamBatchMockLevel.Ctx(ctx).WhereIn("batch_id", batchIDs).OrderAsc("mock_level_id").Scan(&rows); err != nil {
		return nil, err
	}
	m := make(map[int64][]int64)
	for _, r := range rows {
		m[r.BatchId] = append(m[r.BatchId], r.MockLevelId)
	}
	return m, nil
}

func (s *sExam) toBatchAdminItem(b entity.ExamBatch, levelIDs []int64, memberCount int) bo.ExamBatchAdminItem {
	return bo.ExamBatchAdminItem{
		Id:                     b.Id,
		MockExaminationPaperId: b.MockExaminationPaperId,
		Title:                  b.Title,
		ExamStartAt:            b.ExamStartAt,
		ExamEndAt:              b.ExamEndAt,
		MockLevelIds:           levelIDs,
		MemberCount:            memberCount,
		CreateTime:             b.CreateTime,
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
		q = q.Where("mock_examination_paper_id", mockPaperID)
	}
	total, err = q.Count()
	if err != nil {
		return nil, 0, err
	}
	var batches []entity.ExamBatch
	if err = q.Page(page, size).OrderDesc("id").Scan(&batches); err != nil {
		return nil, 0, err
	}
	ids := make([]int64, len(batches))
	for i := range batches {
		ids[i] = batches[i].Id
	}
	levelMap, err := s.mockLevelsByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	cntMap, err := s.memberCountByBatchIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	list = make([]bo.ExamBatchAdminItem, 0, len(batches))
	for _, b := range batches {
		lids := levelMap[b.Id]
		if lids == nil {
			lids = []int64{}
		}
		list = append(list, s.toBatchAdminItem(b, lids, cntMap[b.Id]))
	}
	return list, total, nil
}

// ExamBatchDetail 批次详情（含等级 id 列表与学员数）。
func (s *sExam) ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error) {
	b, err := s.examBatchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	lids, err := s.loadMockLevelIDsForBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	cntMap, err := s.memberCountByBatchIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	item := s.toBatchAdminItem(b, lids, cntMap[id])
	return &item, nil
}

// ExamBatchCreate 创建批次并写入多选 mock_levels。
func (s *sExam) ExamBatchCreate(ctx context.Context, mockPaperID int64, title, examStartAt, examEndAt string, mockLevelIds []int64, creator string) (int64, error) {
	if _, err := exampaper.ByMockID(ctx, mockPaperID); err != nil {
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
	levelIDs := dedupPositiveInt64(mockLevelIds)
	if err := s.ensureMockLevelsExist(ctx, levelIDs); err != nil {
		return 0, err
	}
	var newID int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		row := entity.ExamBatch{
			MockExaminationPaperId: mockPaperID,
			Title:                  strings.TrimSpace(title),
			ExamStartAt:            st,
			ExamEndAt:              en,
			Creator:                creator,
			Updater:                creator,
			CreateTime:             gtime.Now(),
			UpdateTime:             gtime.Now(),
			DeleteFlag:             consts.DeleteFlagNotDeleted,
		}
		pid, err := tx.Model(dao.ExamBatch.Table()).Ctx(ctx).InsertAndGetId(row)
		if err != nil {
			return err
		}
		newID = pid
		if len(levelIDs) > 0 {
			batch := make([]gdb.Map, 0, len(levelIDs))
			for _, lid := range levelIDs {
				batch = append(batch, gdb.Map{
					dao.ExamBatchMockLevel.Columns().BatchId:     newID,
					dao.ExamBatchMockLevel.Columns().MockLevelId: lid,
				})
			}
			if _, err := tx.Model(dao.ExamBatchMockLevel.Table()).Ctx(ctx).Data(batch).Insert(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// ExamBatchUpdate 更新批次时间与等级多选（全量替换等级关联）。
func (s *sExam) ExamBatchUpdate(ctx context.Context, id int64, title, examStartAt, examEndAt string, mockLevelIds []int64, updater string) error {
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
	levelIDs := dedupPositiveInt64(mockLevelIds)
	if err := s.ensureMockLevelsExist(ctx, levelIDs); err != nil {
		return err
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
		if _, err := tx.Model(dao.ExamBatchMockLevel.Table()).Ctx(ctx).Where("batch_id", id).Delete(); err != nil {
			return err
		}
		if len(levelIDs) > 0 {
			batch := make([]gdb.Map, 0, len(levelIDs))
			for _, lid := range levelIDs {
				batch = append(batch, gdb.Map{
					dao.ExamBatchMockLevel.Columns().BatchId:     id,
					dao.ExamBatchMockLevel.Columns().MockLevelId: lid,
				})
			}
			if _, err := tx.Model(dao.ExamBatchMockLevel.Table()).Ctx(ctx).Data(batch).Insert(); err != nil {
				return err
			}
		}
		return nil
	})
}

// ExamBatchDelete 逻辑删除批次（不删学员关联与等级行，便于审计；列表已过滤）。
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

// ExamBatchMembersImport 导入学员；已存在主键则跳过。
func (s *sExam) ExamBatchMembersImport(ctx context.Context, batchID int64, memberIDs []int64, creator string) (inserted int, err error) {
	if _, err := s.examBatchByID(ctx, batchID); err != nil {
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
	if err := dao.ExamBatchMember.Ctx(ctx).Fields("member_id").Where("batch_id", batchID).WhereIn("member_id", ids).Scan(&existing); err != nil {
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
			dao.ExamBatchMember.Columns().BatchId:    batchID,
			dao.ExamBatchMember.Columns().MemberId:   mid,
			dao.ExamBatchMember.Columns().Creator:    creator,
			dao.ExamBatchMember.Columns().CreateTime: now,
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
	var links []entity.ExamBatchMember
	if err = q.Page(page, size).OrderDesc("create_time").Scan(&links); err != nil {
		return nil, 0, err
	}
	if len(links) == 0 {
		return []bo.ExamBatchMemberAdminRow{}, total, nil
	}
	memberIDs := make([]int64, len(links))
	linkTime := make(map[int64]*gtime.Time, len(links))
	for i, l := range links {
		memberIDs[i] = l.MemberId
		linkTime[l.MemberId] = l.CreateTime
	}
	var users []entity.SystemMember
	if err = dao.SysMember.Ctx(ctx).WhereIn("id", memberIDs).Scan(&users); err != nil {
		return nil, 0, err
	}
	um := make(map[int64]entity.SystemMember, len(users))
	for _, u := range users {
		um[u.Id] = u
	}
	list = make([]bo.ExamBatchMemberAdminRow, 0, len(links))
	for _, l := range links {
		u := um[l.MemberId]
		list = append(list, bo.ExamBatchMemberAdminRow{
			MemberId:   l.MemberId,
			Username:   u.Username,
			Nickname:   u.Nickname,
			ImportTime: linkTime[l.MemberId],
		})
	}
	return list, total, nil
}

// ExamBatchMembersRemove 从批次移除学员。
func (s *sExam) ExamBatchMembersRemove(ctx context.Context, batchID int64, memberIDs []int64) (removed int, err error) {
	if _, err := s.examBatchByID(ctx, batchID); err != nil {
		return 0, err
	}
	ids := dedupPositiveInt64(memberIDs)
	if len(ids) == 0 {
		return 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	r, err := dao.ExamBatchMember.Ctx(ctx).Where("batch_id", batchID).WhereIn("member_id", ids).Delete()
	if err != nil {
		return 0, err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
