package mock

import (
	"context"

	v1 "exam/api/client/mock/v1"
)

type IMock interface {
	ExaminationPaperList(ctx context.Context, req *v1.ExaminationPaperListReq) (res *v1.ExaminationPaperListRes, err error)
	ExaminationPaperDetail(ctx context.Context, req *v1.ExaminationPaperDetailReq) (res *v1.ExaminationPaperDetailRes, err error)
}
