// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamPaper is the golang structure for table exam_paper.
type ExamPaper struct {
	Id                      int64       `json:"id"                         orm:"id"                         description:"主键"`                                               // 主键
	Level                   string      `json:"level"                      orm:"level"                      description:"级别，如 hsk1"`                                        // 级别，如 hsk1
	PaperId                 string      `json:"paper_id"                   orm:"paper_id"                   description:"远程试卷目录ID，如 0d26e5c778ad4ca8"`                      // 远程试卷目录ID，如 0d26e5c778ad4ca8
	MockExaminationPaperId  int64       `json:"mock_examination_paper_id"  orm:"mock_examination_paper_id"  description:"mock 真源 mock_examination_paper.id（业务卷标识）"` // mock 真源 mock_examination_paper.id（业务卷标识）
	Title                   string      `json:"title"                      orm:"title"                      description:"试卷标题"`                                             // 试卷标题
	PrepareTitle            string      `json:"prepare_title"              orm:"prepare_title"              description:"prepare.title"`                                    // prepare.title
	PrepareInstruction      string      `json:"prepare_instruction"        orm:"prepare_instruction"        description:"考前说明 instruction"`                                 // 考前说明 instruction
	PrepareAudioFile        string      `json:"prepare_audio_file"         orm:"prepare_audio_file"         description:"prepare.audio_file"`                               // prepare.audio_file
	SourceBaseUrl           string      `json:"source_base_url"            orm:"source_base_url"            description:"资源基址，可拼 index 与媒体 URL"`                            // 资源基址，可拼 index 与媒体 URL
	AudioHlsPrefix          string      `json:"audio_hls_prefix"           orm:"audio_hls_prefix"           description:"OSS 桶内 HLS 目录前缀，动态 m3u8 拼接时使用"`                    // OSS 桶内 HLS 目录前缀，动态 m3u8 拼接时使用
	AudioHlsSegmentCount    int         `json:"audio_hls_segment_count"    orm:"audio_hls_segment_count"    description:"分片总数，合法索引 0..count-1；0 表示未配置 HLS"`                 // 分片总数，合法索引 0..count-1；0 表示未配置 HLS
	AudioHlsSegmentPattern  string      `json:"audio_hls_segment_pattern"  orm:"audio_hls_segment_pattern"  description:"分片文件名 fmt 格式，空则默认 %%05d.ts"`                       // 分片文件名 fmt 格式，空则默认 %%05d.ts
	AudioHlsKeyObject       string      `json:"audio_hls_key_object"       orm:"audio_hls_key_object"       description:"密钥对象相对 prefix 的路径，空表示不加密"`                         // 密钥对象相对 prefix 的路径，空表示不加密
	AudioHlsIvHex           string      `json:"audio_hls_iv_hex"           orm:"audio_hls_iv_hex"           description:"AES-128 IV 十六进制（可选），写入 #EXT-X-KEY"`                // AES-128 IV 十六进制（可选），写入 #EXT-X-KEY
	AudioHlsSegmentDuration float64     `json:"audio_hls_segment_duration" orm:"audio_hls_segment_duration" description:"#EXTINF 时长秒"`                                      // #EXTINF 时长秒
	IndexJson               string      `json:"index_json"                 orm:"index_json"                 description:"index.json 全文快照"`                                  // index.json 全文快照
	DurationSeconds         int         `json:"duration_seconds"           orm:"duration_seconds"           description:"考试时长秒，0=使用系统默认 exam.defaultDurationSeconds"`       // 考试时长秒，0=使用系统默认 exam.defaultDurationSeconds
	Creator                 string      `json:"creator"                    orm:"creator"                    description:"创建者"`                                              // 创建者
	CreateTime              *gtime.Time `json:"create_time"                orm:"create_time"                description:"创建时间"`                                             // 创建时间
	Updater                 string      `json:"updater"                    orm:"updater"                    description:"更新者"`                                              // 更新者
	UpdateTime              *gtime.Time `json:"update_time"                orm:"update_time"                description:"更新时间"`                                             // 更新时间
	DeleteFlag              int         `json:"delete_flag"                orm:"delete_flag"                description:"逻辑删除：0-否，1-是"`                                     // 逻辑删除：0-否，1-是
}
