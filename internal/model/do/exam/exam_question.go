// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamQuestion is the golang structure of table exam_question for DAO operations like Where/Data.
type ExamQuestion struct {
	g.Meta                  `orm:"table:exam_question, do:true"`
	Id                      any         // 主键
	ExamPaperId             any         // 试卷ID，冗余便于按卷查询
	MockExaminationPaperId  any         // 冗余 mock_examination_paper.id
	BlockId                 any         // 题块ID exam_question_block.id
	SortInBlock             any         // 块内顺序
	QuestionNo              any         // 卷面题号（JSON index，如 1-40）
	Score                   any         // 分值
	IsExample               any         // 是否例题
	IsSubjective            any         // 是否主观题：0否 1是（主观题不参与客观自动分）
	ContentType             any         // 题干内容类型，如 audio
	AudioFile               any         // 音频 content 文件名
	AudioHlsPrefix          any         // OSS 桶内 HLS 目录前缀（无首尾/），其下为分片与可选密钥文件
	AudioHlsSegmentCount    any         // 分片总数，合法索引 0..count-1；0 表示未配置 HLS
	AudioHlsSegmentPattern  any         // 分片文件名 fmt 格式，空则默认 %%05d.ts
	AudioHlsKeyObject       any         // 密钥对象相对 prefix 的路径，空表示不加密
	AudioHlsIvHex           any         // AES-128 IV 十六进制（可选），写入 #EXT-X-KEY
	AudioHlsSegmentDuration any         // #EXTINF 时长秒
	StemText                any         // content_sentence
	ScreenTextJson          any         // screen_text 数组
	AnalysisJson            any         // analysis 多语言
	QuestionDescriptionJson any         // 小题级 question_description_obj
	RawJson                 any         // 单题原始 JSON 备份
	Creator                 any         // 创建者
	CreateTime              *gtime.Time // 创建时间
	Updater                 any         // 更新者
	UpdateTime              *gtime.Time // 更新时间
	DeleteFlag              any         // 逻辑删除
}
