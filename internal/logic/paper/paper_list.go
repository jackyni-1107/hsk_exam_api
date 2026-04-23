package paper

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
)

// PaperList 分页查询试卷列表。
func (s *sPaper) PaperList(ctx context.Context, page, size int, level string, mockLevelId int64) (list []examentity.ExamPaper, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	m := dao.ExamPaper.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if mockLevelId > 0 {
		mockPaperIDs, err := mockPaperIDsByLevelID(ctx, mockLevelId)
		if err != nil {
			return nil, 0, err
		}
		if len(mockPaperIDs) == 0 {
			return []examentity.ExamPaper{}, 0, nil
		}
		m = m.WhereIn("mock_examination_paper_id", mockPaperIDs)
	} else if level != "" {
		m = m.Where("level", level)
	}

	n, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	total = n

	if err := m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func mockPaperIDsByLevelID(ctx context.Context, mockLevelID int64) ([]interface{}, error) {
	var mockRows []struct {
		Id int64 `json:"id"`
	}
	if err := dao.MockExaminationPaper.Ctx(ctx).
		Fields("id").
		Where("level_id", mockLevelID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&mockRows); err != nil {
		return nil, err
	}

	ids := make([]interface{}, len(mockRows))
	for i := range mockRows {
		ids[i] = mockRows[i].Id
	}
	return ids, nil
}
