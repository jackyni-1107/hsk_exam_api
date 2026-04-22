package mock

import (
	"context"

	v1 "exam/api/admin/mock/v1"
	mocksvc "exam/internal/service/mock"
)

func (c *ControllerV1) ExaminationPaperList(ctx context.Context, req *v1.ExaminationPaperListReq) (res *v1.ExaminationPaperListRes, err error) {
	list, err := mocksvc.Mock().ExaminationPaperList(ctx, req.LevelId, req.ImportStatus)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ExaminationPaperItem, 0, len(list))
	for _, e := range list {
		items = append(items, &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ResourceUrl: e.ResourceUrl, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
			Imported: e.Imported,
		})
	}
	return &v1.ExaminationPaperListRes{List: items}, nil
}
