// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamQuestionBlockDao is the data access object for the table exam_question_block.
type ExamQuestionBlockDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  ExamQuestionBlockColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// ExamQuestionBlockColumns defines and stores column names for the table exam_question_block.
type ExamQuestionBlockColumns struct {
	Id                      string // 主键
	SectionId               string // 大题ID exam_section.id
	BlockOrder              string // 对应 topic JSON 中 items 下标
	GroupIndex              string // 套题外层 index（若存在）
	QuestionDescriptionJson string // 块级 question_description_obj 等
	Creator                 string // 创建者
	CreateTime              string // 创建时间
	Updater                 string // 更新者
	UpdateTime              string // 更新时间
	DeleteFlag              string // 逻辑删除
}

// examQuestionBlockColumns holds the columns for the table exam_question_block.
var examQuestionBlockColumns = ExamQuestionBlockColumns{
	Id:                      "id",
	SectionId:               "section_id",
	BlockOrder:              "block_order",
	GroupIndex:              "group_index",
	QuestionDescriptionJson: "question_description_json",
	Creator:                 "creator",
	CreateTime:              "create_time",
	Updater:                 "updater",
	UpdateTime:              "update_time",
	DeleteFlag:              "delete_flag",
}

// NewExamQuestionBlockDao creates and returns a new DAO object for table data access.
func NewExamQuestionBlockDao(handlers ...gdb.ModelHandler) *ExamQuestionBlockDao {
	return &ExamQuestionBlockDao{
		group:    "default",
		table:    "exam_question_block",
		columns:  examQuestionBlockColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamQuestionBlockDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamQuestionBlockDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamQuestionBlockDao) Columns() ExamQuestionBlockColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamQuestionBlockDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamQuestionBlockDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamQuestionBlockDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
