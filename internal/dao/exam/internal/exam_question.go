// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExamQuestionDao is the data access object for the table exam_question.
type ExamQuestionDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ExamQuestionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ExamQuestionColumns defines and stores column names for the table exam_question.
type ExamQuestionColumns struct {
	Id                      string // 主键
	ExamPaperId             string // 试卷ID，冗余便于按卷查询
	MockExaminationPaperId  string // 冗余 mock_examination_paper.id
	BlockId                 string // 题块ID exam_question_block.id
	SortInBlock             string // 块内顺序
	QuestionNo              string // 卷面题号（JSON index，如 1-40）
	Score                   string // 分值
	IsExample               string // 是否例题
	IsSubjective            string // 是否主观题：0否 1是（主观题不参与客观自动分）
	ContentType             string // 题干内容类型，如 audio
	AudioFile               string // 音频 content 文件名
	AudioHlsPrefix          string // OSS 桶内 HLS 目录前缀（无首尾/），其下为分片与可选密钥文件
	AudioHlsSegmentCount    string // 分片总数，合法索引 0..count-1；0 表示未配置 HLS
	AudioHlsSegmentPattern  string // 分片文件名 fmt 格式，空则默认 %%05d.ts
	AudioHlsKeyObject       string // 密钥对象相对 prefix 的路径，空表示不加密
	AudioHlsIvHex           string // AES-128 IV 十六进制（可选），写入 #EXT-X-KEY
	AudioHlsSegmentDuration string // #EXTINF 时长秒
	StemText                string // content_sentence
	ScreenTextJson          string // screen_text 数组
	AnalysisJson            string // analysis 多语言
	QuestionDescriptionJson string // 小题级 question_description_obj
	RawJson                 string // 单题原始 JSON 备份
	Creator                 string // 创建者
	CreateTime              string // 创建时间
	Updater                 string // 更新者
	UpdateTime              string // 更新时间
	DeleteFlag              string // 逻辑删除
}

// examQuestionColumns holds the columns for the table exam_question.
var examQuestionColumns = ExamQuestionColumns{
	Id:                      "id",
	ExamPaperId:             "exam_paper_id",
	MockExaminationPaperId:  "mock_examination_paper_id",
	BlockId:                 "block_id",
	SortInBlock:             "sort_in_block",
	QuestionNo:              "question_no",
	Score:                   "score",
	IsExample:               "is_example",
	IsSubjective:            "is_subjective",
	ContentType:             "content_type",
	AudioFile:               "audio_file",
	AudioHlsPrefix:          "audio_hls_prefix",
	AudioHlsSegmentCount:    "audio_hls_segment_count",
	AudioHlsSegmentPattern:  "audio_hls_segment_pattern",
	AudioHlsKeyObject:       "audio_hls_key_object",
	AudioHlsIvHex:           "audio_hls_iv_hex",
	AudioHlsSegmentDuration: "audio_hls_segment_duration",
	StemText:                "stem_text",
	ScreenTextJson:          "screen_text_json",
	AnalysisJson:            "analysis_json",
	QuestionDescriptionJson: "question_description_json",
	RawJson:                 "raw_json",
	Creator:                 "creator",
	CreateTime:              "create_time",
	Updater:                 "updater",
	UpdateTime:              "update_time",
	DeleteFlag:              "delete_flag",
}

// NewExamQuestionDao creates and returns a new DAO object for table data access.
func NewExamQuestionDao(handlers ...gdb.ModelHandler) *ExamQuestionDao {
	return &ExamQuestionDao{
		group:    "default",
		table:    "exam_question",
		columns:  examQuestionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExamQuestionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExamQuestionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExamQuestionDao) Columns() ExamQuestionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExamQuestionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExamQuestionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExamQuestionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
