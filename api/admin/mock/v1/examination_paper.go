package v1

import "github.com/gogf/gf/v2/frame/g"

type ExaminationPaperListReq struct {
	g.Meta       `path:"/mock/examination-paper/list" method:"get" tags:"Mock-管理" summary:"模拟卷列表（管理端）；import_status 筛选是否已导入 exam"`
	LevelId      int64  `json:"level_id"`
	ImportStatus string `json:"import_status" dc:"空或 all：全部；imported：仅已导入；not_imported：仅未导入"`
}

type ExaminationPaperListRes struct {
	List []*ExaminationPaperItem `json:"list"`
}

type ExaminationPaperItem struct {
	Id        int64  `json:"id"`
	LevelId   int64  `json:"level_id"`
	Name      string `json:"name"`
	ScoreFull int    `json:"score_full"`
	TimeFull  int    `json:"time_full"`
	Status    int    `json:"status"`
	PaperType int    `json:"paper_type"`
	MockType  int    `json:"mock_type"`
	Imported  bool   `json:"imported" dc:"是否已在 exam 域导入（存在未删除的 exam_paper）"`
}

type ExaminationPaperDetailReq struct {
	g.Meta `path:"/mock/examination-paper/{id}" method:"get" tags:"Mock-管理" summary:"模拟卷详情（管理端）"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type ExaminationPaperDetailRes struct {
	Paper *ExaminationPaperItem `json:"paper"`
}
