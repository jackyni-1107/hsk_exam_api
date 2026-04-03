// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MockExaminationPartDao is the data access object for the table mock_examination_part.
type MockExaminationPartDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  MockExaminationPartColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// MockExaminationPartColumns defines and stores column names for the table mock_examination_part.
type MockExaminationPartColumns struct {
	Id                      string // id主键
	Code                    string // 部分编号
	SegmentId               string // 环节id
	QuestionCount           string // 题目数量
	ObjectiveQuestionCount  string // 客观题数量
	SubjectiveQuestionCount string // 主观题数量
	QuestionScore           string // 每个题目得分
	PartScore               string // 该题型分数
	AnswerTime              string // 每个题型的回答时间
	PartName                string // 只用于hskk
	PartNameTrans           string // 部分名称多语言
	DeleteFlag              string // 是否删除
	CreateTime              string // 创建时间
	UpdateTime              string // 更新时间
}

// mockExaminationPartColumns holds the columns for the table mock_examination_part.
var mockExaminationPartColumns = MockExaminationPartColumns{
	Id:                      "id",
	Code:                    "code",
	SegmentId:               "segment_id",
	QuestionCount:           "question_count",
	ObjectiveQuestionCount:  "objective_question_count",
	SubjectiveQuestionCount: "subjective_question_count",
	QuestionScore:           "question_score",
	PartScore:               "part_score",
	AnswerTime:              "answer_time",
	PartName:                "part_name",
	PartNameTrans:           "part_name_trans",
	DeleteFlag:              "delete_flag",
	CreateTime:              "create_time",
	UpdateTime:              "update_time",
}

// NewMockExaminationPartDao creates and returns a new DAO object for table data access.
func NewMockExaminationPartDao(handlers ...gdb.ModelHandler) *MockExaminationPartDao {
	return &MockExaminationPartDao{
		group:    "default",
		table:    "mock_examination_part",
		columns:  mockExaminationPartColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MockExaminationPartDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MockExaminationPartDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MockExaminationPartDao) Columns() MockExaminationPartColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MockExaminationPartDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MockExaminationPartDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MockExaminationPartDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
