// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MockExaminationPaperDao is the data access object for the table mock_examination_paper.
type MockExaminationPaperDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  MockExaminationPaperColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// MockExaminationPaperColumns defines and stores column names for the table mock_examination_paper.
type MockExaminationPaperColumns struct {
	Id                   string // id 主键
	LevelId              string // HSK等级
	Name                 string // 试卷名称
	NameTrans            string // 试卷名称国际化
	ScoreFull            string // 满分
	Seq                  string // 试卷顺序
	ExplainAudio         string // 说明音频
	TimeFull             string // 考试时长
	ListenReviewDuration string // 听力结束后回顾时间
	TimeSheet            string // 答题卡时间
	IconUrl              string // 图标路径
	ResourceUrl          string // 资源包路径
	Version              string // 资源版本号
	DeleteFlag           string // 是否删除 0 未删除 1 已删除
	MockType             string // 1 hsk 2 hskk 3 yct
	ProductBaseId        string // 学习需要付的元商品
	Credit               string // 所需元商品数量
	BuyProductId         string // 当付费券不足时, 需要购买此商品
	CreateTime           string // 创建时间
	UpdateTime           string // 更新时间
	Status               string // 试卷状态 0 未发布 1 发布 2 下架
	MemberResource       string // 0 非会员资源1 会员资源
	PaperType            string // 1 模拟考试卷 2 在线练习试卷
}

// mockExaminationPaperColumns holds the columns for the table mock_examination_paper.
var mockExaminationPaperColumns = MockExaminationPaperColumns{
	Id:                   "id",
	LevelId:              "level_id",
	Name:                 "name",
	NameTrans:            "name_trans",
	ScoreFull:            "score_full",
	Seq:                  "seq",
	ExplainAudio:         "explain_audio",
	TimeFull:             "time_full",
	ListenReviewDuration: "listen_review_duration",
	TimeSheet:            "time_sheet",
	IconUrl:              "icon_url",
	ResourceUrl:          "resource_url",
	Version:              "version",
	DeleteFlag:           "delete_flag",
	MockType:             "mock_type",
	ProductBaseId:        "product_base_id",
	Credit:               "credit",
	BuyProductId:         "buy_product_id",
	CreateTime:           "create_time",
	UpdateTime:           "update_time",
	Status:               "status",
	MemberResource:       "member_resource",
	PaperType:            "paper_type",
}

// NewMockExaminationPaperDao creates and returns a new DAO object for table data access.
func NewMockExaminationPaperDao(handlers ...gdb.ModelHandler) *MockExaminationPaperDao {
	return &MockExaminationPaperDao{
		group:    "default",
		table:    "mock_examination_paper",
		columns:  mockExaminationPaperColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MockExaminationPaperDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MockExaminationPaperDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MockExaminationPaperDao) Columns() MockExaminationPaperColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MockExaminationPaperDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MockExaminationPaperDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MockExaminationPaperDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
