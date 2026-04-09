package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamBatchMockPaperDao is the data access object for the table exam_batch_mock_paper.
type ExamBatchMockPaperDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  ExamBatchMockPaperColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// ExamBatchMockPaperColumns defines and stores column names for the table exam_batch_mock_paper.
type ExamBatchMockPaperColumns struct {
	BatchId                string // exam_batch.id
	MockExaminationPaperId string // mock_examination_paper.id
}

var examBatchMockPaperColumns = ExamBatchMockPaperColumns{
	BatchId:                "batch_id",
	MockExaminationPaperId: "mock_examination_paper_id",
}

// NewExamBatchMockPaperDao creates and returns a new DAO object for table data access.
func NewExamBatchMockPaperDao(handlers ...gdb.ModelHandler) *ExamBatchMockPaperDao {
	return &ExamBatchMockPaperDao{
		group:    "default",
		table:    "exam_batch_mock_paper",
		columns:  examBatchMockPaperColumns,
		handlers: handlers,
	}
}

func (dao *ExamBatchMockPaperDao) DB() gdb.DB {
	return g.DB(dao.group)
}

func (dao *ExamBatchMockPaperDao) Table() string {
	return dao.table
}

func (dao *ExamBatchMockPaperDao) Columns() ExamBatchMockPaperColumns {
	return dao.columns
}

func (dao *ExamBatchMockPaperDao) Group() string {
	return dao.group
}

func (dao *ExamBatchMockPaperDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

func (dao *ExamBatchMockPaperDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
