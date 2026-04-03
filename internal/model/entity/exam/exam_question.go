// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamQuestion is the golang structure for table exam_question.
type ExamQuestion struct {
	Id                      int64       `json:"id"                         orm:"id"                         description:"主键"`                                 // 主键
	ExamPaperId             int64       `json:"exam_paper_id"              orm:"exam_paper_id"              description:"试卷ID，冗余便于按卷查询"`                      // 试卷ID，冗余便于按卷查询
	MockExaminationPaperId  int64       `json:"mock_examination_paper_id"  orm:"mock_examination_paper_id"  description:"冗余 mock_examination_paper.id"`       // 冗余 mock_examination_paper.id
	BlockId                 int64       `json:"block_id"                   orm:"block_id"                   description:"题块ID exam_question_block.id"`        // 题块ID exam_question_block.id
	SortInBlock             int         `json:"sort_in_block"              orm:"sort_in_block"              description:"块内顺序"`                               // 块内顺序
	QuestionNo              int         `json:"question_no"                orm:"question_no"                description:"卷面题号（JSON index，如 1-40）"`            // 卷面题号（JSON index，如 1-40）
	Score                   float64     `json:"score"                      orm:"score"                      description:"分值"`                                 // 分值
	IsExample               int         `json:"is_example"                 orm:"is_example"                 description:"是否例题"`                               // 是否例题
	IsSubjective            int         `json:"is_subjective"              orm:"is_subjective"              description:"是否主观题：0否 1是（主观题不参与客观自动分）"`           // 是否主观题：0否 1是（主观题不参与客观自动分）
	ContentType             string      `json:"content_type"               orm:"content_type"               description:"题干内容类型，如 audio"`                     // 题干内容类型，如 audio
	AudioFile               string      `json:"audio_file"                 orm:"audio_file"                 description:"音频 content 文件名"`                     // 音频 content 文件名
	AudioHlsPrefix          string      `json:"audio_hls_prefix"           orm:"audio_hls_prefix"           description:"OSS 桶内 HLS 目录前缀（无首尾/），其下为分片与可选密钥文件"` // OSS 桶内 HLS 目录前缀（无首尾/），其下为分片与可选密钥文件
	AudioHlsSegmentCount    int         `json:"audio_hls_segment_count"    orm:"audio_hls_segment_count"    description:"分片总数，合法索引 0..count-1；0 表示未配置 HLS"`   // 分片总数，合法索引 0..count-1；0 表示未配置 HLS
	AudioHlsSegmentPattern  string      `json:"audio_hls_segment_pattern"  orm:"audio_hls_segment_pattern"  description:"分片文件名 fmt 格式，空则默认 %%05d.ts"`         // 分片文件名 fmt 格式，空则默认 %%05d.ts
	AudioHlsKeyObject       string      `json:"audio_hls_key_object"       orm:"audio_hls_key_object"       description:"密钥对象相对 prefix 的路径，空表示不加密"`           // 密钥对象相对 prefix 的路径，空表示不加密
	AudioHlsIvHex           string      `json:"audio_hls_iv_hex"           orm:"audio_hls_iv_hex"           description:"AES-128 IV 十六进制（可选），写入 #EXT-X-KEY"`  // AES-128 IV 十六进制（可选），写入 #EXT-X-KEY
	AudioHlsSegmentDuration float64     `json:"audio_hls_segment_duration" orm:"audio_hls_segment_duration" description:"#EXTINF 时长秒"`                        // #EXTINF 时长秒
	StemText                string      `json:"stem_text"                  orm:"stem_text"                  description:"content_sentence"`                   // content_sentence
	ScreenTextJson          string      `json:"screen_text_json"           orm:"screen_text_json"           description:"screen_text 数组"`                     // screen_text 数组
	AnalysisJson            string      `json:"analysis_json"              orm:"analysis_json"              description:"analysis 多语言"`                       // analysis 多语言
	QuestionDescriptionJson string      `json:"question_description_json"  orm:"question_description_json"  description:"小题级 question_description_obj"`       // 小题级 question_description_obj
	RawJson                 string      `json:"raw_json"                   orm:"raw_json"                   description:"单题原始 JSON 备份"`                       // 单题原始 JSON 备份
	Creator                 string      `json:"creator"                    orm:"creator"                    description:"创建者"`                                // 创建者
	CreateTime              *gtime.Time `json:"create_time"                orm:"create_time"                description:"创建时间"`                               // 创建时间
	Updater                 string      `json:"updater"                    orm:"updater"                    description:"更新者"`                                // 更新者
	UpdateTime              *gtime.Time `json:"update_time"                orm:"update_time"                description:"更新时间"`                               // 更新时间
	DeleteFlag              int         `json:"delete_flag"                orm:"delete_flag"                description:"逻辑删除"`                               // 逻辑删除
}
