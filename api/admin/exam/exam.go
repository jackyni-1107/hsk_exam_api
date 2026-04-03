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
	AttemptList(ctx context.Context, req *v1.AttemptListReq) (res *v1.AttemptListRes, err error)
	AttemptDetail(ctx context.Context, req *v1.AttemptDetailReq) (res *v1.AttemptDetailRes, err error)
	AttemptSubjectiveScores(ctx context.Context, req *v1.AttemptSubjectiveScoresReq) (res *v1.AttemptSubjectiveScoresRes, err error)
}
