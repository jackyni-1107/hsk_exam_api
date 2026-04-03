// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamAttemptAnswerDao is the data access object for the table exam_attempt_answer.
type ExamAttemptAnswerDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  ExamAttemptAnswerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// ExamAttemptAnswerColumns defines and stores column names for the table exam_attempt_answer.
type ExamAttemptAnswerColumns struct {
	Id             string // 主键
	AttemptId      string // exam_attempt.id
	ExamQuestionId string // exam_question.id
	AnswerJson     string // 用户答案JSON
	AwardedScore   string // 主观题人工得分；NULL 表示未评
	Version        string // 乐观锁版本
	Creator        string // 创建者
	CreateTime     string // 创建时间
	Updater        string // 更新者
	UpdateTime     string // 更新时间
	DeleteFlag     string // 逻辑删除
}

// examAttemptAnswerColumns holds the columns for the table exam_attempt_answer.
var examAttemptAnswerColumns = ExamAttemptAnswerColumns{
	Id:             "id",
	AttemptId:      "attempt_id",
	ExamQuestionId: "exam_question_id",
	AnswerJson:     "answer_json",
	AwardedScore:   "awarded_score",
	Version:        "version",
	Creator:        "creator",
	CreateTime:     "create_time",
	Updater:        "updater",
	UpdateTime:     "update_time",
	DeleteFlag:     "delete_flag",
}

// NewExamAttemptAnswerDao creates and returns a new DAO object for table data access.
func NewExamAttemptAnswerDao(handlers ...gdb.ModelHandler) *ExamAttemptAnswerDao {
	return &ExamAttemptAnswerDao{
		group:    "default",
		table:    "exam_attempt_answer",
		columns:  examAttemptAnswerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamAttemptAnswerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamAttemptAnswerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamAttemptAnswerDao) Columns() ExamAttemptAnswerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamAttemptAnswerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamAttemptAnswerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamAttemptAnswerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
