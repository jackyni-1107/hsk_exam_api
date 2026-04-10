package mock

import (
	"context"

	v1 "exam/api/admin/mock/v1"
	"exam/internal/consts"
	mocksvc "exam/internal/service/mock"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) ExaminationPaperList(ctx context.Context, req *v1.ExaminationPaperListReq) (res *v1.ExaminationPaperListRes, err error) {
	list, err := mocksvc.Mock().ExaminationPaperList(ctx, req.LevelId, req.ImportStatus)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ExaminationPaperItem, 0, len(list))
	for _, e := range list {
		items = append(items, &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
			Imported: e.Imported,
		})
	}
	return &v1.ExaminationPaperListRes{List: items}, nil
}

func (c *ControllerV1) ExaminationPaperDetail(ctx context.Context, req *v1.ExaminationPaperDetailReq) (res *v1.ExaminationPaperDetailRes, err error) {
	e, err := mocksvc.Mock().ExaminationPaperDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, gerror.NewCode(consts.CodeResourceNotFound)
	}
	return &v1.ExaminationPaperDetailRes{
		Paper: &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
			Imported: e.Imported,
		},
	}, nil
}
