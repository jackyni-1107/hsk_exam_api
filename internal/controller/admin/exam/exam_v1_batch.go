package exam

import (
	"context"

	v1 "exam/api/admin/exam/v1"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	"exam/internal/service/batch"
	"exam/internal/utility"
)

func batchListItemPtr(b *bo.ExamBatchAdminItem) *v1.BatchListItem {
	ids := b.ExamPaperIds
	if ids == nil {
		ids = []int64{}
	}
	formattedStartAt := utility.ToRFC3339UTC(b.ExamStartAt)
	formattedEndAt := utility.ToRFC3339UTC(b.ExamEndAt)
	return &v1.BatchListItem{
		Id:                    b.Id,
		ExamPaperIds:          ids,
		Title:                 b.Title,
		ExamStartAt:           formattedStartAt,
		ExamEndAt:             formattedEndAt,
		BatchKind:             b.BatchKind,
		AllowMultipleAttempts: b.AllowMultipleAttempts,
		MaxAttemptsPerMember:  b.MaxAttemptsPerMember,
		SkipScoring:           b.SkipScoring,
		AutoSubmitOnDeadline:  b.AutoSubmitOnDeadline,
		MemberCount:           b.MemberCount,
		CreateTime:            utility.ToRFC3339UTC(b.CreateTime),
	}
}

func (c *ControllerV1) BatchList(ctx context.Context, req *v1.BatchListReq) (res *v1.BatchListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	rows, total, err := batch.Batch().ExamBatchList(ctx, req.ExamPaperId, req.Page, req.Size, req.TimeFrom, req.TimeTo)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.BatchListItem, 0, len(rows))
	for i := range rows {
		list = append(list, batchListItemPtr(&rows[i]))
	}
	return &v1.BatchListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) BatchDetail(ctx context.Context, req *v1.BatchDetailReq) (res *v1.BatchDetailRes, err error) {
	b, err := batch.Batch().ExamBatchDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.BatchDetailRes{Batch: *batchListItemPtr(b)}, nil
}

func (c *ControllerV1) BatchCreate(ctx context.Context, req *v1.BatchCreateReq) (res *v1.BatchCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	autoSubmit := 1
	if req.AutoSubmitOnDeadline != nil {
		autoSubmit = *req.AutoSubmitOnDeadline
	}
	policy := bo.ExamBatchPolicyInput{
		BatchKind:             req.BatchKind,
		AllowMultipleAttempts: req.AllowMultipleAttempts,
		MaxAttemptsPerMember:  req.MaxAttemptsPerMember,
		SkipScoring:           req.SkipScoring,
		AutoSubmitOnDeadline:  autoSubmit,
	}
	id, err := batch.Batch().ExamBatchCreate(ctx, req.Title, req.ExamStartAt, req.ExamEndAt, req.ExamPaperIds, creator, policy)
	if err != nil {
		return nil, err
	}
	return &v1.BatchCreateRes{Id: id}, nil
}

func (c *ControllerV1) BatchUpdate(ctx context.Context, req *v1.BatchUpdateReq) (res *v1.BatchUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	detail, err := batch.Batch().ExamBatchDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	autoSubmit := detail.AutoSubmitOnDeadline
	if req.AutoSubmitOnDeadline != nil {
		autoSubmit = *req.AutoSubmitOnDeadline
	}
	policy := bo.ExamBatchPolicyInput{
		BatchKind:             req.BatchKind,
		AllowMultipleAttempts: req.AllowMultipleAttempts,
		MaxAttemptsPerMember:  req.MaxAttemptsPerMember,
		SkipScoring:           req.SkipScoring,
		AutoSubmitOnDeadline:  autoSubmit,
	}
	if err := batch.Batch().ExamBatchUpdate(ctx, req.Id, req.Title, req.ExamStartAt, req.ExamEndAt, req.ExamPaperIds, updater, policy); err != nil {
		return nil, err
	}
	return &v1.BatchUpdateRes{}, nil
}

func (c *ControllerV1) BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error) {
	if err := batch.Batch().ExamBatchDelete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.BatchDeleteRes{}, nil
}

func (c *ControllerV1) BatchMembersImport(ctx context.Context, req *v1.BatchMembersImportReq) (res *v1.BatchMembersImportRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	n, err := batch.Batch().ExamBatchMembersAdd(ctx, req.Id, req.ExamPaperId, req.MemberIds, creator)
	if err != nil {
		return nil, err
	}
	return &v1.BatchMembersImportRes{Inserted: n}, nil
}

func (c *ControllerV1) BatchMemberList(ctx context.Context, req *v1.BatchMemberListReq) (res *v1.BatchMemberListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	rows, total, err := batch.Batch().ExamBatchMemberList(ctx, req.Id, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.BatchMemberListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, &v1.BatchMemberListItem{
			MemberId:    r.MemberId,
			ExamPaperId: r.ExamPaperId,
			Username:    r.Username,
			Nickname:    r.Nickname,
			ImportTime:  utility.ToRFC3339UTC(r.ImportTime),
		})
	}
	return &v1.BatchMemberListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) BatchMembersRemove(ctx context.Context, req *v1.BatchMembersRemoveReq) (res *v1.BatchMembersRemoveRes, err error) {
	n, err := batch.Batch().ExamBatchMembersRemove(ctx, req.Id, req.ExamPaperId, req.MemberIds)
	if err != nil {
		return nil, err
	}
	return &v1.BatchMembersRemoveRes{Removed: n}, nil
}
