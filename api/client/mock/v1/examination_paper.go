package v1

import "github.com/gogf/gf/v2/frame/g"

type ExaminationPaperListReq struct {
	g.Meta  `path:"/mock/examination-paper/list" method:"get" tags:"Mock" summary:"模拟卷列表"`
	LevelId int64 `json:"level_id"`
}

type ExaminationPaperListRes struct {
	List []*ExaminationPaperItem `json:"list"`
}

type ExaminationPaperItem struct {
	Id          int64  `json:"id"`
	LevelId     int64  `json:"level_id"`
	Name        string `json:"name"`
	ScoreFull   int    `json:"score_full"`
	TimeFull    int    `json:"time_full"`
	Status      int    `json:"status"`
	PaperType   int    `json:"paper_type"`
	MockType    int    `json:"mock_type"`
}

type ExaminationPaperDetailReq struct {
	g.Meta `path:"/mock/examination-paper/{id}" method:"get" tags:"Mock" summary:"模拟卷详情"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type ExaminationPaperDetailRes struct {
	Paper *ExaminationPaperItem `json:"paper"`
}
