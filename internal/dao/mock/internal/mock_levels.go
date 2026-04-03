// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MockLevelsDao is the data access object for the table mock_levels.
type MockLevelsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MockLevelsColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MockLevelsColumns defines and stores column names for the table mock_levels.
type MockLevelsColumns struct {
	Id                 string // 主键
	LevelId            string // 在指定类型下的等级
	LevelType          string // 等级类型 （1 hsk 2 hskk 3 yct）
	TypeName           string // 类型名称
	LevelName          string // hsk等级名称
	AppLevelName       string // app专用等级名称
	LevelNameTrans     string // 等级名称多语言
	SignUpUrl          string // 报名地址
	LevelDesc          string // 等级简介
	CnDesc             string // 中文等级描述
	EnDesc             string // 英文等级描述
	LevelExplainAudio  string // 模拟考试音音频地址
	DeleteFlag         string // 逻辑删除标识
	MarkFlag           string // 是否需要批改
	CreateTime         string // 创建时间
	UpdateTime         string // 更新时间
	ExamShowStatus     string // 模拟考显示状态，0表示不显示，1表示显示
	HomeworkShowStatus string // 作业显示状态，0表示不显示，1表示显示
	PaperAnswerTime    string // 该等级下试卷作答时间
	ResourceType       string // 1：考试列表 2：文件资源列表
}

// mockLevelsColumns holds the columns for the table mock_levels.
var mockLevelsColumns = MockLevelsColumns{
	Id:                 "id",
	LevelId:            "level_id",
	LevelType:          "level_type",
	TypeName:           "type_name",
	LevelName:          "level_name",
	AppLevelName:       "app_level_name",
	LevelNameTrans:     "level_name_trans",
	SignUpUrl:          "sign_up_url",
	LevelDesc:          "level_desc",
	CnDesc:             "cn_desc",
	EnDesc:             "en_desc",
	LevelExplainAudio:  "level_explain_audio",
	DeleteFlag:         "delete_flag",
	MarkFlag:           "mark_flag",
	CreateTime:         "create_time",
	UpdateTime:         "update_time",
	ExamShowStatus:     "exam_show_status",
	HomeworkShowStatus: "homework_show_status",
	PaperAnswerTime:    "paper_answer_time",
	ResourceType:       "resource_type",
}

// NewMockLevelsDao creates and returns a new DAO object for table data access.
func NewMockLevelsDao(handlers ...gdb.ModelHandler) *MockLevelsDao {
	return &MockLevelsDao{
		group:    "default",
		table:    "mock_levels",
		columns:  mockLevelsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MockLevelsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MockLevelsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MockLevelsDao) Columns() MockLevelsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MockLevelsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MockLevelsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *MockLevelsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
