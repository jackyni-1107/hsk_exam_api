package exam

import (
	"context"

	v1 "exam/api/admin/exam/v1"
)

// IExam 管理端考试接口（与 controller/admin/exam 对齐）。
type IExam interface {
	PaperList(ctx context.Context, req *v1.PaperListReq) (res *v1.PaperListRes, err error)
	PaperDetail(ctx context.Context, req *v1.PaperDetailReq) (res *v1.PaperDetailRes, err error)
	PaperImport(ctx context.Context, req *v1.PaperImportReq) (res *v1.PaperImportRes, err error)
	PaperUpdate(ctx context.Context, req *v1.PaperUpdateReq) (res *v1.PaperUpdateRes, err error)
	PaperEdit(ctx context.Context, req *v1.PaperEditReq) (res *v1.PaperEditRes, err error)
	PaperPurge(ctx context.Context, req *v1.PaperPurgeReq) (res *v1.PaperPurgeRes, err error)
	AttemptList(ctx context.Context, req *v1.AttemptListReq) (res *v1.AttemptListRes, err error)
	AttemptDetail(ctx context.Context, req *v1.AttemptDetailReq) (res *v1.AttemptDetailRes, err error)
	AttemptSubjectiveScores(ctx context.Context, req *v1.AttemptSubjectiveScoresReq) (res *v1.AttemptSubjectiveScoresRes, err error)
	BatchList(ctx context.Context, req *v1.BatchListReq) (res *v1.BatchListRes, err error)
	BatchDetail(ctx context.Context, req *v1.BatchDetailReq) (res *v1.BatchDetailRes, err error)
	BatchCreate(ctx context.Context, req *v1.BatchCreateReq) (res *v1.BatchCreateRes, err error)
	BatchUpdate(ctx context.Context, req *v1.BatchUpdateReq) (res *v1.BatchUpdateRes, err error)
	BatchDelete(ctx context.Context, req *v1.BatchDeleteReq) (res *v1.BatchDeleteRes, err error)
	BatchMembersImport(ctx context.Context, req *v1.BatchMembersImportReq) (res *v1.BatchMembersImportRes, err error)
	BatchMemberList(ctx context.Context, req *v1.BatchMemberListReq) (res *v1.BatchMemberListRes, err error)
	BatchMembersRemove(ctx context.Context, req *v1.BatchMembersRemoveReq) (res *v1.BatchMembersRemoveRes, err error)
}
