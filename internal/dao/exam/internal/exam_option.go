// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamOptionDao is the data access object for the table exam_option.
type ExamOptionDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExamOptionColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExamOptionColumns defines and stores column names for the table exam_option.
type ExamOptionColumns struct {
	Id         string // 主键
	QuestionId string // 试题ID exam_question.id
	Flag       string // 选项标识 A/B/C/T/F
	SortOrder  string // 对应 answers.index
	IsCorrect  string // 是否正确
	OptionType string // text/image/pinyin 等
	Content    string // 文本或资源文件名
	Creator    string // 创建者
	CreateTime string // 创建时间
	Updater    string // 更新者
	UpdateTime string // 更新时间
	DeleteFlag string // 逻辑删除
}

// examOptionColumns holds the columns for the table exam_option.
var examOptionColumns = ExamOptionColumns{
	Id:         "id",
	QuestionId: "question_id",
	Flag:       "flag",
	SortOrder:  "sort_order",
	IsCorrect:  "is_correct",
	OptionType: "option_type",
	Content:    "content",
	Creator:    "creator",
	CreateTime: "create_time",
	Updater:    "updater",
	UpdateTime: "update_time",
	DeleteFlag: "delete_flag",
}

// NewExamOptionDao creates and returns a new DAO object for table data access.
func NewExamOptionDao(handlers ...gdb.ModelHandler) *ExamOptionDao {
	return &ExamOptionDao{
		group:    "default",
		table:    "exam_option",
		columns:  examOptionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamOptionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamOptionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamOptionDao) Columns() ExamOptionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamOptionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamOptionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamOptionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
