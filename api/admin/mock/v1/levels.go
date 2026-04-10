package v1

import "github.com/gogf/gf/v2/frame/g"

type MockLevelsListReq struct {
	g.Meta `path:"/mock/levels/list" method:"get" tags:"Mock" summary:"Mock 等级列表"`
}

type MockLevelsListRes struct {
	List []*MockLevelItem `json:"list" dc:"等级列表"`
}

type MockLevelItem struct {
	Id                 int64  `json:"id" dc:"记录ID"`
	LevelId            int    `json:"level_id" dc:"等级ID"`
	LevelType          int    `json:"level_type" dc:"等级类型"`
	TypeName           string `json:"type_name" dc:"类型名称"`
	LevelName          string `json:"level_name" dc:"等级名称"`
	AppLevelName       string `json:"app_level_name" dc:"APP等级名称"`
	ExamShowStatus     int    `json:"exam_show_status" dc:"考试显示状态"`
	HomeworkShowStatus int    `json:"homework_show_status" dc:"作业显示状态"`
}
