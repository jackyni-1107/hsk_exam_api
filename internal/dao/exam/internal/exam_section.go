// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamSectionDao is the data access object for the table exam_section.
type ExamSectionDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExamSectionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExamSectionColumns defines and stores column names for the table exam_section.
type ExamSectionColumns struct {
	Id                     string // 主键
	ExamPaperId            string // 试卷ID exam_paper.id
	MockExaminationPaperId string // 冗余 mock_examination_paper.id
	SortOrder              string // 在 index.items 中的顺序
	TopicTitle     string // topic_title
	TopicSubtitle  string // topic_subtitle
	TopicType      string // 题型代码 pt/xp/xt/...
	PartCode       string // 大题内 part 序号
	SegmentCode    string // listen/read
	TopicItemsFile string // topic_items 文件名，如 pt.json
	TopicJson      string // 该 topic 文件全文快照
	Creator        string // 创建者
	CreateTime     string // 创建时间
	Updater        string // 更新者
	UpdateTime     string // 更新时间
	DeleteFlag     string // 逻辑删除
}

// examSectionColumns holds the columns for the table exam_section.
var examSectionColumns = ExamSectionColumns{
	Id:                     "id",
	ExamPaperId:            "exam_paper_id",
	MockExaminationPaperId: "mock_examination_paper_id",
	SortOrder:              "sort_order",
	TopicTitle:     "topic_title",
	TopicSubtitle:  "topic_subtitle",
	TopicType:      "topic_type",
	PartCode:       "part_code",
	SegmentCode:    "segment_code",
	TopicItemsFile: "topic_items_file",
	TopicJson:      "topic_json",
	Creator:        "creator",
	CreateTime:     "create_time",
	Updater:        "updater",
	UpdateTime:     "update_time",
	DeleteFlag:     "delete_flag",
}

// NewExamSectionDao creates and returns a new DAO object for table data access.
func NewExamSectionDao(handlers ...gdb.ModelHandler) *ExamSectionDao {
	return &ExamSectionDao{
		group:    "default",
		table:    "exam_section",
		columns:  examSectionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamSectionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamSectionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamSectionDao) Columns() ExamSectionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamSectionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamSectionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamSectionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
