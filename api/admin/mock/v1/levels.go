package v1

import "github.com/gogf/gf/v2/frame/g"

type MockLevelsListReq struct {
	g.Meta `path:"/mock/levels/list" method:"get" tags:"Mock" summary:"Mock 等级列表"`
}

type MockLevelsListRes struct {
	List []*MockLevelItem `json:"list"`
}

type MockLevelItem struct {
	Id                 int64  `json:"id"`
	LevelId            int    `json:"level_id"`
	LevelType          int    `json:"level_type"`
	TypeName           string `json:"type_name"`
	LevelName          string `json:"level_name"`
	AppLevelName       string `json:"app_level_name"`
	ExamShowStatus     int    `json:"exam_show_status"`
	HomeworkShowStatus int    `json:"homework_show_status"`
}
