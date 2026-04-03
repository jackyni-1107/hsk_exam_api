// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamResultDao is the data access object for the table exam_result.
type ExamResultDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExamResultColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExamResultColumns defines and stores column names for the table exam_result.
type ExamResultColumns struct {
	AttemptId              string // exam_attempt.id
	MemberId               string // sys_member.id
	ExamPaperId            string // exam_paper.id
	MockExaminationPaperId string // 冗余 mock_examination_paper.id
	ExamBatchId            string // 冗余 exam_batch.id
	MockLevelId            string // 冗余 mock_levels.id
	Status                 string // 与 exam_attempt.status，列表阶段为已结束 4
	ObjectiveScore         string //
	SubjectiveScore        string //
	TotalScore             string //
	HasSubjective          string //
	StartedAt              string //
	SubmittedAt            string //
	EndedAt                string //
	CreateTime             string // 同步自会话创建时间
	UpdateTime             string //
	DeleteFlag             string //
}

// examResultColumns holds the columns for the table exam_result.
var examResultColumns = ExamResultColumns{
	AttemptId:              "attempt_id",
	MemberId:               "member_id",
	ExamPaperId:            "exam_paper_id",
	MockExaminationPaperId: "mock_examination_paper_id",
	ExamBatchId:            "exam_batch_id",
	MockLevelId:            "mock_level_id",
	Status:                 "status",
	ObjectiveScore:         "objective_score",
	SubjectiveScore:        "subjective_score",
	TotalScore:             "total_score",
	HasSubjective:          "has_subjective",
	StartedAt:              "started_at",
	SubmittedAt:            "submitted_at",
	EndedAt:                "ended_at",
	CreateTime:             "create_time",
	UpdateTime:             "update_time",
	DeleteFlag:             "delete_flag",
}

// NewExamResultDao creates and returns a new DAO object for table data access.
func NewExamResultDao(handlers ...gdb.ModelHandler) *ExamResultDao {
	return &ExamResultDao{
		group:    "default",
		table:    "exam_result",
		columns:  examResultColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamResultDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamResultDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamResultDao) Columns() ExamResultColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamResultDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamResultDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamResultDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
