// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamAttemptDao is the data access object for the table exam_attempt.
type ExamAttemptDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExamAttemptColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExamAttemptColumns defines and stores column names for the table exam_attempt.
type ExamAttemptColumns struct {
	Id                     string // 主键
	MemberId               string // sys_member.id
	ExamPaperId            string // exam_paper.id
	MockExaminationPaperId string // 冗余 mock_examination_paper.id
	Status                 string // 1=not_started 2=in_progress 3=submitted 4=ended
	DurationSeconds        string // 开考时快照时长秒
	StartedAt              string // 开考时间
	DeadlineAt             string // 截止时间
	SubmittedAt            string // 交卷时间
	EndedAt                string // 结束/阅卷完成时间
	ObjectiveScore         string // 客观题得分
	SubjectiveScore        string // 主观题得分（未批阅时为空）
	TotalScore             string // 总分（无主观或主观未批时可能为客观分）
	HasSubjective          string // 本卷是否含主观题（交卷快照）
	Creator                string // 创建者
	CreateTime             string // 创建时间
	Updater                string // 更新者
	UpdateTime             string // 更新时间
	DeleteFlag             string // 逻辑删除
}

// examAttemptColumns holds the columns for the table exam_attempt.
var examAttemptColumns = ExamAttemptColumns{
	Id:                     "id",
	MemberId:               "member_id",
	ExamPaperId:            "exam_paper_id",
	MockExaminationPaperId: "mock_examination_paper_id",
	Status:                 "status",
	DurationSeconds:        "duration_seconds",
	StartedAt:              "started_at",
	DeadlineAt:             "deadline_at",
	SubmittedAt:            "submitted_at",
	EndedAt:                "ended_at",
	ObjectiveScore:         "objective_score",
	SubjectiveScore:        "subjective_score",
	TotalScore:             "total_score",
	HasSubjective:          "has_subjective",
	Creator:                "creator",
	CreateTime:             "create_time",
	Updater:                "updater",
	UpdateTime:             "update_time",
	DeleteFlag:             "delete_flag",
}

// NewExamAttemptDao creates and returns a new DAO object for table data access.
func NewExamAttemptDao(handlers ...gdb.ModelHandler) *ExamAttemptDao {
	return &ExamAttemptDao{
		group:    "default",
		table:    "exam_attempt",
		columns:  examAttemptColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamAttemptDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamAttemptDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamAttemptDao) Columns() ExamAttemptColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamAttemptDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamAttemptDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamAttemptDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
