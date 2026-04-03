// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamPaper is the golang structure of table exam_paper for DAO operations like Where/Data.
type ExamPaper struct {
	g.Meta                  `orm:"table:exam_paper, do:true"`
	Id                      any         // 主键
	Level                   any         // 级别，如 hsk1
	PaperId                 any         // 远程试卷目录ID，如 0d26e5c778ad4ca8
	MockExaminationPaperId  any         // mock 真源 mock_examination_paper.id；NULL=非 mock 导入
	Title                   any         // 试卷标题
	PrepareTitle            any         // prepare.title
	PrepareInstruction      any         // 考前说明 instruction
	PrepareAudioFile        any         // prepare.audio_file
	SourceBaseUrl           any         // 资源基址，可拼 index 与媒体 URL
	AudioHlsPrefix          any         // OSS 桶内 HLS 目录前缀，动态 m3u8 拼接时使用
	AudioHlsSegmentCount    any         // 分片总数，合法索引 0..count-1；0 表示未配置 HLS
	AudioHlsSegmentPattern  any         // 分片文件名 fmt 格式，空则默认 %%05d.ts
	AudioHlsKeyObject       any         // 密钥对象相对 prefix 的路径，空表示不加密
	AudioHlsIvHex           any         // AES-128 IV 十六进制（可选），写入 #EXT-X-KEY
	AudioHlsSegmentDuration any         // #EXTINF 时长秒
	IndexJson               any         // index.json 全文快照
	DurationSeconds         any         // 考试时长秒，0=使用系统默认 exam.defaultDurationSeconds
	Creator                 any         // 创建者
	CreateTime              *gtime.Time // 创建时间
	Updater                 any         // 更新者
	UpdateTime              *gtime.Time // 更新时间
	DeleteFlag              any         // 逻辑删除：0-否，1-是
}
