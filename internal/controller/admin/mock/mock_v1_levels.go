package mock

import (
	"context"

	v1 "exam/api/admin/mock/v1"
	mocksvc "exam/internal/service/mock"
)

func (c *ControllerV1) MockLevelsList(ctx context.Context, req *v1.MockLevelsListReq) (res *v1.MockLevelsListRes, err error) {
	list, err := mocksvc.Mock().MockLevelsList(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.MockLevelItem, 0, len(list))
	for _, e := range list {
		items = append(items, &v1.MockLevelItem{
			Id:                 e.Id,
			LevelId:            e.LevelId,
			LevelType:          e.LevelType,
			TypeName:           e.TypeName,
			LevelName:          e.LevelName,
			AppLevelName:       e.AppLevelName,
			ExamShowStatus:     e.ExamShowStatus,
			HomeworkShowStatus: e.HomeworkShowStatus,
		})
	}
	return &v1.MockLevelsListRes{List: items}, nil
}
