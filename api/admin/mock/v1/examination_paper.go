package v1

import "github.com/gogf/gf/v2/frame/g"

type ExaminationPaperListReq struct {
	g.Meta       `path:"/mock/examination-paper/list" method:"get" tags:"Mock-管理" summary:"模拟卷列表（管理端）；import_status 筛选是否已导入 exam"`
	LevelId      int64  `json:"level_id" dc:"等级ID"`
	ImportStatus string `json:"import_status" dc:"空或 all：全部；imported：仅已导入；not_imported：仅未导入"`
}

type ExaminationPaperListRes struct {
	List []*ExaminationPaperItem `json:"list" dc:"模拟卷列表"`
}

type ExaminationPaperItem struct {
	Id          int64  `json:"id" dc:"模拟卷ID"`
	LevelId     int64  `json:"level_id" dc:"等级ID"`
	Name        string `json:"name" dc:"卷名"`
	ResourceUrl string `json:"resource_url" dc:"资源包 URL（.zip），导入时替换为 /index.json"`
	ScoreFull   int    `json:"score_full" dc:"满分"`
	TimeFull    int    `json:"time_full" dc:"考试时长(分钟)"`
	Status      int    `json:"status" dc:"状态"`
	PaperType   int    `json:"paper_type" dc:"试卷类型"`
	MockType    int    `json:"mock_type" dc:"模拟类型"`
	Imported    bool   `json:"imported" dc:"是否已在 exam 域导入（存在未删除的 exam_paper）"`
}

