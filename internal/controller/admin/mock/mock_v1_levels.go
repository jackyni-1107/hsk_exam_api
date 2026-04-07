package mock

import (
	"context"

	v1 "exam/api/admin/mock/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	mockentity "exam/internal/model/entity/mock"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MockLevelsList(ctx context.Context, req *v1.MockLevelsListReq) (res *v1.MockLevelsListRes, err error) {
	var list []mockentity.MockLevels
	err = dao.MockLevels.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
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
