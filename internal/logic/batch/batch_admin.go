package batch

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"
)

// ExamBatchList 分页查询考试批次列表
func (s *sBatch) ExamBatchList(ctx context.Context, page, size int, key string) (list []examentity.ExamBatch, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := dao.ExamBatch.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if key != "" {
		m = m.WhereLike("name", "%"+key+"%")
	}

	total, err = m.Count()
	if err != nil || total == 0 {
		return nil, 0, gerror.NewCode(consts.CodeDataNotFound)
	}

	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.NewCode(consts.CodeDataNotFound)
	}

	return list, total, nil
}

// ExamBatchCreate 创建考试批次
func (s *sBatch) ExamBatchCreate(ctx context.Context, name string, startAt, endAt string, creator string) (int64, error) {
	st := s.parseTime(startAt)
	et := s.parseTime(endAt)
	if st == nil || et == nil {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}

	r, err := dao.ExamBatch.Ctx(ctx).Insert(g.Map{
		"name":          name,
		"exam_start_at": st,
		"exam_end_at":   et,
		"creator":       creator,
		"delete_flag":   consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.NewCode(consts.CodeDataNotFound)
	}
	return r.LastInsertId()
}

// ExamBatchUpdate 更新考试批次
func (s *sBatch) ExamBatchUpdate(ctx context.Context, id int64, name string, startAt, endAt string, updater string) error {
	st := s.parseTime(startAt)
	et := s.parseTime(endAt)

	data := g.Map{"updater": updater}
	if name != "" {
		data["name"] = name
	}
	if st != nil {
		data["exam_start_at"] = st
	}
	if et != nil {
		data["exam_end_at"] = et
	}

	_, err := dao.ExamBatch.Ctx(ctx).Where("id", id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()

	if err != nil {
		return gerror.NewCode(consts.CodeExamBatchNotFound)
	}
	return nil
}

// ExamBatchDelete 删除考试批次
func (s *sBatch) ExamBatchDelete(ctx context.Context, id int64, updater string) error {
	_, err := dao.ExamBatch.Ctx(ctx).Where("id", id).Data(g.Map{
		"delete_flag": consts.DeleteFlagDeleted,
		"updater":     updater,
	}).Update()

	if err != nil {
		return gerror.NewCode(consts.CodeExamBatchNotFound)
	}
	return nil
}

// ExamBatchMembersAdd 批量向指定批次和 Mock 卷添加学员
func (s *sBatch) ExamBatchMembersAdd(ctx context.Context, batchID int64, mockPaperID int64, memberIDs []int64, creator string) error {
	// 验证批次是否存在 (调用 init.go 中的公用函数)
	if _, err := s.GetBatchByID(ctx, batchID); err != nil {
		return err
	}

	ids := s.dedupIDs(memberIDs)
	if len(ids) == 0 {
		return gerror.NewCode(consts.CodeExamBatchMemberNotFound)
	}

	return dao.ExamBatchMember.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, mid := range ids {
			// 检查是否已存在，防止重复插入
			count, _ := dao.ExamBatchMember.Ctx(ctx).
				Where(g.Map{"batch_id": batchID, "mock_examination_paper_id": mockPaperID, "member_id": mid}).
				Count()
			if count > 0 {
				continue
			}

			_, err := dao.ExamBatchMember.Ctx(ctx).Insert(g.Map{
				"batch_id":                  batchID,
				"mock_examination_paper_id": mockPaperID,
				"member_id":                 mid,
				"creator":                   creator,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// ExamBatchMembersRemove 从批次中移除学员
func (s *sBatch) ExamBatchMembersRemove(ctx context.Context, batchID int64, mockPaperID int64, memberIDs []int64) (int, error) {
	ids := s.dedupIDs(memberIDs)
	if len(ids) == 0 {
		return 0, gerror.NewCode(consts.CodeExamBatchMemberNotFound)
	}

	r, err := dao.ExamBatchMember.Ctx(ctx).
		Where("batch_id", batchID).
		Where("mock_examination_paper_id", mockPaperID).
		WhereIn("member_id", ids).Delete()

	if err != nil {
		return 0, gerror.NewCode(consts.CodeExamBatchMemberNotFound)
	}
	n, _ := r.RowsAffected()
	return int(n), nil
}

// ExamBatchMemberList 查询批次内的成员列表（关联系统用户表）
func (s *sBatch) ExamBatchMemberList(ctx context.Context, page, size int, batchID int64, mockPaperID int64) (list []sysentity.SysUser, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	m := g.DB().Ctx(ctx).Model("sys_user u").
		InnerJoin("exam_batch_member ebm", "u.id = ebm.member_id").
		Where("ebm.batch_id", batchID).
		Where("ebm.mock_examination_paper_id", mockPaperID)

	total, err = m.Count()
	if err != nil || total == 0 {
		return nil, 0, nil
	}

	err = m.Page(page, size).OrderAsc("u.id").Scan(&list)
	return list, total, err
}
