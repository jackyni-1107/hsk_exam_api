// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamBatchDao is the data access object for the table exam_batch.
type ExamBatchDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExamBatchColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExamBatchColumns defines and stores column names for the table exam_batch.
type ExamBatchColumns struct {
	Id                     string // 主键
	MockExaminationPaperId string // mock 卷 id，与 exam_paper 业务主键一致
	Title                  string // 批次名称
	ExamStartAt            string // 考试允许开始时间
	ExamEndAt              string // 考试允许结束时间
	Creator                string // 创建者
	CreateTime             string // 创建时间
	Updater                string // 更新者
	UpdateTime             string // 更新时间
	DeleteFlag             string // 逻辑删除：0-否，1-是
}

// examBatchColumns holds the columns for the table exam_batch.
var examBatchColumns = ExamBatchColumns{
	Id:                     "id",
	MockExaminationPaperId: "mock_examination_paper_id",
	Title:                  "title",
	ExamStartAt:            "exam_start_at",
	ExamEndAt:              "exam_end_at",
	Creator:                "creator",
	CreateTime:             "create_time",
	Updater:                "updater",
	UpdateTime:             "update_time",
	DeleteFlag:             "delete_flag",
}

// NewExamBatchDao creates and returns a new DAO object for table data access.
func NewExamBatchDao(handlers ...gdb.ModelHandler) *ExamBatchDao {
	return &ExamBatchDao{
		group:    "default",
		table:    "exam_batch",
		columns:  examBatchColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamBatchDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamBatchDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamBatchDao) Columns() ExamBatchColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamBatchDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamBatchDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamBatchDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
