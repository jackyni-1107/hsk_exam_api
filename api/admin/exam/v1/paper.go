package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PaperListReq struct {
	g.Meta `path:"/exam/paper/list" method:"get" tags:"试卷管理" summary:"试卷列表"`
	Level  string `json:"level" dc:"级别筛选，如 hsk1"`
	Page   int    `json:"page" dc:"页码"`
	Size   int    `json:"size" dc:"每页条数"`
}

type PaperListRes struct {
	List  []*PaperListItem `json:"list" dc:"列表"`
	Total int              `json:"total" dc:"总数"`
}

type PaperListItem struct {
	Id                      int64   `json:"id" dc:"mock_examination_paper.id"`
	Level                   string  `json:"level" dc:"试卷级别"`
	PaperId                 string  `json:"paper_id" dc:"远程试卷ID"`
	Title                   string  `json:"title" dc:"试卷标题"`
	SourceBaseUrl           string  `json:"source_base_url" dc:"资源基础URL"`
	AudioHlsPrefix          string  `json:"audio_hls_prefix" dc:"听力 HLS 目录前缀"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count" dc:"HLS 分片总数"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern" dc:"HLS 分片文件名模式"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object" dc:"HLS 密钥对象路径"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex" dc:"AES-128 IV 十六进制"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration" dc:"HLS 分片时长(秒)"`
	CreateTime              string  `json:"create_time" dc:"创建时间"`
}

type PaperImportReq struct {
	g.Meta                 `path:"/exam/paper/import" method:"post" tags:"试卷管理" summary:"从远程 index.json 导入试卷"`
	MockExaminationPaperId int64 `json:"mock_examination_paper_id" v:"required|min:1#err.invalid_params" dc:"mock 卷 id，业务主键"`
	// 二选一：index_url 拉取；或 index_json 粘贴（需同时传 level、paper_id、source_base_url）
	IndexUrl       string `json:"index_url" dc:"完整 index.json URL"`
	IndexJson      string `json:"index_json" dc:"可选，直接传 index.json 字符串"`
	Level          string `json:"level" dc:"粘贴 index_json 时必填"`
	PaperId        string `json:"paper_id" dc:"粘贴 index_json 时必填"`
	SourceBaseUrl  string `json:"source_base_url" dc:"粘贴 index_json 时必填，以 / 结尾的基址"`
	AudioHlsPrefix string `json:"audio_hls_prefix" dc:"听力 HLS 目录前缀（无首尾/），动态 m3u8 拼接时使用"`
	// fail：已存在则返回 conflict；overwrite：覆盖；new_copy：使用 new_paper_id 作为新试卷
	ConflictMode string `json:"conflict_mode" dc:"fail|overwrite|new_copy，默认 fail"`
	NewPaperId   string `json:"new_paper_id" dc:"new_copy 时新的远程 paper_id"`
}

type PaperImportRes struct {
	ExaminationPaperId         int64 `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	Conflict                   bool  `json:"conflict" dc:"true 表示未导入，因已存在且 conflict_mode=fail"`
	ExistingExaminationPaperId int64 `json:"existing_examination_paper_id" dc:"冲突时已存在的 mock 卷 id"`
	SectionCount               int   `json:"section_count" dc:"大题数"`
	QuestionCount              int   `json:"question_count" dc:"试题数"`
}

type PaperUpdateReq struct {
	g.Meta `path:"/exam/paper/update" method:"post" tags:"试卷管理" summary:"修改听力 HLS 配置；答题时长以 mock 卷为准。id 为 mock 卷 id，未导入时 code=11114"`
	Id     int64 `json:"id" v:"required|min:1#err.invalid_params" dc:"mock_examination_paper.id"`
	// 听力 HLS：与 exam_paper 表字段一致
	AudioHlsPrefix          string  `json:"audio_hls_prefix" dc:"目录前缀（无首尾/）"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count" v:"min:0#err.invalid_params" dc:"分片总数，0 表示未配置"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern" dc:"分片文件名 fmt，空则默认 %%05d.ts"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object" dc:"密钥对象相对 prefix 的路径"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex" dc:"AES-128 IV 十六进制"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration" v:"min:0#err.invalid_params" dc:"#EXTINF 时长秒"`
}

type PaperUpdateRes struct{}

type PaperDetailReq struct {
	g.Meta `path:"/exam/paper/{id}" method:"get" tags:"试卷管理" summary:"试卷详情（大题/题块/试题）；path id 为 mock_examination_paper.id。未导入时 code=11114，mock 不存在时 code=11201"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"mock_examination_paper.id"`
}

type PaperDetailRes struct {
	Paper    PaperDetailPaper     `json:"paper" dc:"试卷信息"`
	Sections []PaperDetailSection `json:"sections" dc:"大题列表"`
}

type PaperDetailPaper struct {
	Id                      int64   `json:"id" dc:"mock_examination_paper.id"`
	Level                   string  `json:"level" dc:"试卷级别"`
	PaperId                 string  `json:"paper_id" dc:"远程试卷ID"`
	Title                   string  `json:"title" dc:"废弃 试卷标题"`
	Name                    string  `json:"name" dc:"试卷标题"`
	PrepareTitle            string  `json:"prepare_title" dc:"准备阶段标题"`
	PrepareInstruction      string  `json:"prepare_instruction" dc:"准备阶段说明"`
	PrepareAudioFile        string  `json:"prepare_audio_file" dc:"准备阶段音频文件"`
	SourceBaseUrl           string  `json:"source_base_url" dc:"资源基础URL"`
	AudioHlsPrefix          string  `json:"audio_hls_prefix" dc:"听力 HLS 目录前缀"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count" dc:"HLS 分片总数"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern" dc:"HLS 分片文件名模式"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object" dc:"HLS 密钥对象路径"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex" dc:"AES-128 IV 十六进制"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration" dc:"HLS 分片时长(秒)"`
	IndexJson               string  `json:"index_json" dc:"原始 index.json 内容"`
	CreateTime              string  `json:"create_time" dc:"创建时间"`
}

type PaperDetailSection struct {
	Id             int64              `json:"id" dc:"大题ID"`
	SortOrder      int                `json:"sort_order" dc:"排序"`
	TopicTitle     string             `json:"topic_title" dc:"大题标题"`
	TopicSubtitle  string             `json:"topic_subtitle" dc:"大题副标题"`
	TopicType      string             `json:"topic_type" dc:"题型"`
	PartCode       int                `json:"part_code" dc:"部分编号"`
	SegmentCode    string             `json:"segment_code" dc:"段落编号"`
	TopicItemsFile string             `json:"topic_items_file" dc:"题目文件名"`
	TopicJson      string             `json:"topic_json" dc:"大题JSON"`
	Blocks         []PaperDetailBlock `json:"blocks" dc:"题块列表"`
}

type PaperDetailBlock struct {
	Id                      int64                 `json:"id" dc:"题块ID"`
	BlockOrder              int                   `json:"block_order" dc:"题块排序"`
	GroupIndex              int                   `json:"group_index" dc:"组索引"`
	QuestionDescriptionJson string                `json:"question_description_json" dc:"题块描述(JSON)"`
	Questions               []PaperDetailQuestion `json:"questions" dc:"试题列表"`
}

type PaperDetailQuestion struct {
	Id                      int64               `json:"id" dc:"试题ID"`
	SortInBlock             int                 `json:"sort_in_block" dc:"块内排序"`
	QuestionNo              int                 `json:"question_no" dc:"题号"`
	Score                   float64             `json:"score" dc:"分值"`
	IsExample               int                 `json:"is_example" dc:"是否例题：0否 1是"`
	ContentType             string              `json:"content_type" dc:"内容类型"`
	AudioFile               string              `json:"audio_file" dc:"音频文件"`
	StemText                string              `json:"stem_text" dc:"题干文本"`
	ScreenTextJson          string              `json:"screen_text_json" dc:"屏幕文本(JSON)"`
	AnalysisJson            string              `json:"analysis_json" dc:"解析(JSON)"`
	QuestionDescriptionJson string              `json:"question_description_json" dc:"题目描述(JSON)"`
	RawJson                 string              `json:"raw_json" dc:"原始JSON"`
	Options                 []PaperDetailOption `json:"options" dc:"选项列表"`
}

type PaperDetailOption struct {
	Id         int64  `json:"id" dc:"选项ID"`
	Flag       string `json:"flag" dc:"选项标识"`
	SortOrder  int    `json:"sort_order" dc:"排序"`
	IsCorrect  int    `json:"is_correct" dc:"是否正确：0否 1是"`
	OptionType string `json:"option_type" dc:"选项类型"`
	Content    string `json:"content" dc:"选项内容"`
}
