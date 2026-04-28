// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamAttemptCheatEventDao is the data access object for the table exam_attempt_cheat_event.
type ExamAttemptCheatEventDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  ExamAttemptCheatEventColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// ExamAttemptCheatEventColumns defines and stores column names for the table exam_attempt_cheat_event.
type ExamAttemptCheatEventColumns struct {
	Id          string // 主键
	AttemptId   string // exam_attempt.id
	MemberId    string // sys_member.id
	EventType   string // 作弊事件类型，如 switch_screen/screen_record
	EventAt     string // 事件发生时间（服务端写入，与请求到达时刻一致）
	SegmentCode string // 发生时环节编码 listen/read/write
	Detail      string // 事件详情
	ClientIp    string // 客户端IP
	ClientAgent string // User-Agent
	Creator     string // 创建者
	CreateTime  string // 创建时间
	DeleteFlag  string // 逻辑删除
}

// examAttemptCheatEventColumns holds the columns for the table exam_attempt_cheat_event.
var examAttemptCheatEventColumns = ExamAttemptCheatEventColumns{
	Id:          "id",
	AttemptId:   "attempt_id",
	MemberId:    "member_id",
	EventType:   "event_type",
	EventAt:     "event_at",
	SegmentCode: "segment_code",
	Detail:      "detail",
	ClientIp:    "client_ip",
	ClientAgent: "client_agent",
	Creator:     "creator",
	CreateTime:  "create_time",
	DeleteFlag:  "delete_flag",
}

// NewExamAttemptCheatEventDao creates and returns a new DAO object for table data access.
func NewExamAttemptCheatEventDao(handlers ...gdb.ModelHandler) *ExamAttemptCheatEventDao {
	return &ExamAttemptCheatEventDao{
		group:    "default",
		table:    "exam_attempt_cheat_event",
		columns:  examAttemptCheatEventColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamAttemptCheatEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamAttemptCheatEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamAttemptCheatEventDao) Columns() ExamAttemptCheatEventColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamAttemptCheatEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamAttemptCheatEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamAttemptCheatEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
