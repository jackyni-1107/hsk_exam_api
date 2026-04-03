package mock

import (
	"context"

	v1 "exam/api/admin/mock/v1"
)

type IMock interface {
	MockLevelsList(ctx context.Context, req *v1.MockLevelsListReq) (res *v1.MockLevelsListRes, err error)
	ExaminationPaperList(ctx context.Context, req *v1.ExaminationPaperListReq) (res *v1.ExaminationPaperListRes, err error)
	ExaminationPaperDetail(ctx context.Context, req *v1.ExaminationPaperDetailReq) (res *v1.ExaminationPaperDetailRes, err error)
}

