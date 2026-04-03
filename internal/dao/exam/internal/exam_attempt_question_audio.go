// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamAttemptQuestionAudioDao is the data access object for the table exam_attempt_question_audio.
type ExamAttemptQuestionAudioDao struct {
	table    string                          // table is the underlying table name of the DAO.
	group    string                          // group is the database configuration group name of the current DAO.
	columns  ExamAttemptQuestionAudioColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler              // handlers for customized model modification.
}

// ExamAttemptQuestionAudioColumns defines and stores column names for the table exam_attempt_question_audio.
type ExamAttemptQuestionAudioColumns struct {
	Id              string // 主键
	AttemptId       string // exam_attempt.id
	ExamQuestionId  string // exam_question.id
	MaxSegmentIndex string // 允许播放到该分片索引（含）；与 0..segment_count-1 求交
	Creator         string // 创建者
	CreateTime      string // 创建时间
	Updater         string // 更新者
	UpdateTime      string // 更新时间
	DeleteFlag      string // 逻辑删除
}

// examAttemptQuestionAudioColumns holds the columns for the table exam_attempt_question_audio.
var examAttemptQuestionAudioColumns = ExamAttemptQuestionAudioColumns{
	Id:              "id",
	AttemptId:       "attempt_id",
	ExamQuestionId:  "exam_question_id",
	MaxSegmentIndex: "max_segment_index",
	Creator:         "creator",
	CreateTime:      "create_time",
	Updater:         "updater",
	UpdateTime:      "update_time",
	DeleteFlag:      "delete_flag",
}

// NewExamAttemptQuestionAudioDao creates and returns a new DAO object for table data access.
func NewExamAttemptQuestionAudioDao(handlers ...gdb.ModelHandler) *ExamAttemptQuestionAudioDao {
	return &ExamAttemptQuestionAudioDao{
		group:    "default",
		table:    "exam_attempt_question_audio",
		columns:  examAttemptQuestionAudioColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamAttemptQuestionAudioDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamAttemptQuestionAudioDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamAttemptQuestionAudioDao) Columns() ExamAttemptQuestionAudioColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamAttemptQuestionAudioDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamAttemptQuestionAudioDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamAttemptQuestionAudioDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
