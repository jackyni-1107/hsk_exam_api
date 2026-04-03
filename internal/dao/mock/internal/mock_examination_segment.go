// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MockExaminationSegmentDao is the data access object for the table mock_examination_segment.
type MockExaminationSegmentDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  MockExaminationSegmentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// MockExaminationSegmentColumns defines and stores column names for the table mock_examination_segment.
type MockExaminationSegmentColumns struct {
	Id               string // id 主键
	LevelId          string // HSK等级ID
	SegmentCode      string // 环节编码
	SegmentName      string // 环节名称
	SegmentNameTrans string // 环节名称国际化
	SegmentDesc      string // segment说明
	ScoreFull        string // 环节满分
	QuestionCount    string // 题目数量
	Duration         string // 环节时长 分钟数
	Seq              string // 顺序号
	DeleteFlag       string // 是否删除
	CreateTime       string // 创建时间
	UpdateTime       string // 更新时间
}

// mockExaminationSegmentColumns holds the columns for the table mock_examination_segment.
var mockExaminationSegmentColumns = MockExaminationSegmentColumns{
	Id:               "id",
	LevelId:          "level_id",
	SegmentCode:      "segment_code",
	SegmentName:      "segment_name",
	SegmentNameTrans: "segment_name_trans",
	SegmentDesc:      "segment_desc",
	ScoreFull:        "score_full",
	QuestionCount:    "question_count",
	Duration:         "duration",
	Seq:              "seq",
	DeleteFlag:       "delete_flag",
	CreateTime:       "create_time",
	UpdateTime:       "update_time",
}

// NewMockExaminationSegmentDao creates and returns a new DAO object for table data access.
func NewMockExaminationSegmentDao(handlers ...gdb.ModelHandler) *MockExaminationSegmentDao {
	return &MockExaminationSegmentDao{
		group:    "default",
		table:    "mock_examination_segment",
		columns:  mockExaminationSegmentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MockExaminationSegmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MockExaminationSegmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MockExaminationSegmentDao) Columns() MockExaminationSegmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MockExaminationSegmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MockExaminationSegmentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MockExaminationSegmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
