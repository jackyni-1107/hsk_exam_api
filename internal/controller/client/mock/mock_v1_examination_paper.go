package mock

import (
	"context"

	v1 "exam/api/client/mock/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	mockentity "exam/internal/model/entity/mock"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) ExaminationPaperList(ctx context.Context, req *v1.ExaminationPaperListReq) (res *v1.ExaminationPaperListRes, err error) {
	m := dao.MockExaminationPaper.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted).Where("status", 1)
	if req.LevelId > 0 {
		m = m.Where("level_id", req.LevelId)
	}
	var list []mockentity.MockExaminationPaper
	err = m.OrderAsc("seq").OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.ExaminationPaperItem, 0, len(list))
	for _, e := range list {
		items = append(items, &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
		})
	}
	return &v1.ExaminationPaperListRes{List: items}, nil
}

func (c *ControllerV1) ExaminationPaperDetail(ctx context.Context, req *v1.ExaminationPaperDetailReq) (res *v1.ExaminationPaperDetailRes, err error) {
	var e mockentity.MockExaminationPaper
	err = dao.MockExaminationPaper.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.not_found")
	}
	return &v1.ExaminationPaperDetailRes{
		Paper: &v1.ExaminationPaperItem{
			Id: e.Id, LevelId: e.LevelId, Name: e.Name, ScoreFull: e.ScoreFull,
			TimeFull: e.TimeFull, Status: e.Status, PaperType: e.PaperType, MockType: e.MockType,
		},
	}, nil
}
