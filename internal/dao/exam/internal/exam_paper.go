// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamPaperDao is the data access object for the table exam_paper.
type ExamPaperDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExamPaperColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExamPaperColumns defines and stores column names for the table exam_paper.
type ExamPaperColumns struct {
	Id                      string // 主键
	Level                   string // 级别，如 hsk1
	PaperId                 string // 远程试卷目录ID，如 0d26e5c778ad4ca8
	MockExaminationPaperId  string // mock 真源 mock_examination_paper.id
	Title                   string // 试卷标题
	PrepareTitle            string // prepare.title
	PrepareInstruction      string // 考前说明 instruction
	PrepareAudioFile        string // prepare.audio_file
	SourceBaseUrl           string // 资源基址，可拼 index 与媒体 URL
	AudioHlsPrefix          string // OSS 桶内 HLS 目录前缀，动态 m3u8 拼接时使用
	AudioHlsSegmentCount    string // 分片总数，合法索引 0..count-1；0 表示未配置 HLS
	AudioHlsSegmentPattern  string // 分片文件名 fmt 格式，空则默认 %%05d.ts
	AudioHlsKeyObject       string // 密钥对象相对 prefix 的路径，空表示不加密
	AudioHlsIvHex           string // AES-128 IV 十六进制（可选），写入 #EXT-X-KEY
	AudioHlsSegmentDuration string // #EXTINF 时长秒
	IndexJson               string // index.json 全文快照
	DurationSeconds         string // 考试时长秒，0=使用系统默认 exam.defaultDurationSeconds
	Creator                 string // 创建者
	CreateTime              string // 创建时间
	Updater                 string // 更新者
	UpdateTime              string // 更新时间
	DeleteFlag              string // 逻辑删除：0-否，1-是
}

// examPaperColumns holds the columns for the table exam_paper.
var examPaperColumns = ExamPaperColumns{
	Id:                      "id",
	Level:                   "level",
	PaperId:                 "paper_id",
	MockExaminationPaperId:  "mock_examination_paper_id",
	Title:                   "title",
	PrepareTitle:            "prepare_title",
	PrepareInstruction:      "prepare_instruction",
	PrepareAudioFile:        "prepare_audio_file",
	SourceBaseUrl:           "source_base_url",
	AudioHlsPrefix:          "audio_hls_prefix",
	AudioHlsSegmentCount:    "audio_hls_segment_count",
	AudioHlsSegmentPattern:  "audio_hls_segment_pattern",
	AudioHlsKeyObject:       "audio_hls_key_object",
	AudioHlsIvHex:           "audio_hls_iv_hex",
	AudioHlsSegmentDuration: "audio_hls_segment_duration",
	IndexJson:               "index_json",
	DurationSeconds:         "duration_seconds",
	Creator:                 "creator",
	CreateTime:              "create_time",
	Updater:                 "updater",
	UpdateTime:              "update_time",
	DeleteFlag:              "delete_flag",
}

// NewExamPaperDao creates and returns a new DAO object for table data access.
func NewExamPaperDao(handlers ...gdb.ModelHandler) *ExamPaperDao {
	return &ExamPaperDao{
		group:    "default",
		table:    "exam_paper",
		columns:  examPaperColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamPaperDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamPaperDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamPaperDao) Columns() ExamPaperColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamPaperDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamPaperDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamPaperDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
