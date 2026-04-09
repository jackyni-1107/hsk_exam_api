package v1

import "github.com/gogf/gf/v2/frame/g"

type PaperListReq struct {
	g.Meta `path:"/exam/paper/list" method:"get" tags:"试卷管理" summary:"试卷列表"`
	Level  string `json:"level" dc:"级别筛选，如 hsk1"`
	Page   int    `json:"page" dc:"页码"`
	Size   int    `json:"size" dc:"每页条数"`
}

type PaperListRes struct {
	List  []*PaperListItem `json:"list"`
	Total int              `json:"total"`
}

type PaperListItem struct {
	Id                      int64   `json:"id" dc:"mock_examination_paper.id"`
	Level                   string  `json:"level"`
	PaperId                 string  `json:"paper_id"`
	Title                   string  `json:"title"`
	SourceBaseUrl           string  `json:"source_base_url"`
	AudioHlsPrefix          string  `json:"audio_hls_prefix"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration"`
	CreateTime              string  `json:"create_time"`
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
	SectionCount               int   `json:"section_count"`
	QuestionCount              int   `json:"question_count"`
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
	Paper    PaperDetailPaper     `json:"paper"`
	Sections []PaperDetailSection `json:"sections"`
}

type PaperDetailPaper struct {
	Id                      int64   `json:"id" dc:"mock_examination_paper.id"`
	Level                   string  `json:"level"`
	PaperId                 string  `json:"paper_id"`
	Title                   string  `json:"title"`
	PrepareTitle            string  `json:"prepare_title"`
	PrepareInstruction      string  `json:"prepare_instruction"`
	PrepareAudioFile        string  `json:"prepare_audio_file"`
	SourceBaseUrl           string  `json:"source_base_url"`
	AudioHlsPrefix          string  `json:"audio_hls_prefix"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration"`
	IndexJson               string  `json:"index_json"`
	CreateTime              string  `json:"create_time"`
}

type PaperDetailSection struct {
	Id             int64              `json:"id"`
	SortOrder      int                `json:"sort_order"`
	TopicTitle     string             `json:"topic_title"`
	TopicSubtitle  string             `json:"topic_subtitle"`
	TopicType      string             `json:"topic_type"`
	PartCode       int                `json:"part_code"`
	SegmentCode    string             `json:"segment_code"`
	TopicItemsFile string             `json:"topic_items_file"`
	TopicJson      string             `json:"topic_json"`
	Blocks         []PaperDetailBlock `json:"blocks"`
}

type PaperDetailBlock struct {
	Id                      int64                 `json:"id"`
	BlockOrder              int                   `json:"block_order"`
	GroupIndex              int                   `json:"group_index"`
	QuestionDescriptionJson string                `json:"question_description_json"`
	Questions               []PaperDetailQuestion `json:"questions"`
}

type PaperDetailQuestion struct {
	Id                      int64               `json:"id"`
	SortInBlock             int                 `json:"sort_in_block"`
	QuestionNo              int                 `json:"question_no"`
	Score                   float64             `json:"score"`
	IsExample               int                 `json:"is_example"`
	ContentType             string              `json:"content_type"`
	AudioFile               string              `json:"audio_file"`
	StemText                string              `json:"stem_text"`
	ScreenTextJson          string              `json:"screen_text_json"`
	AnalysisJson            string              `json:"analysis_json"`
	QuestionDescriptionJson string              `json:"question_description_json"`
	RawJson                 string              `json:"raw_json"`
	Options                 []PaperDetailOption `json:"options"`
}

type PaperDetailOption struct {
	Id         int64  `json:"id"`
	Flag       string `json:"flag"`
	SortOrder  int    `json:"sort_order"`
	IsCorrect  int    `json:"is_correct"`
	OptionType string `json:"option_type"`
	Content    string `json:"content"`
}
