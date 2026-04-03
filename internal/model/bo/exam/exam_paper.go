package exam

type PaperDetailTree struct {
	Paper    PaperHeadView       `json:"paper"`
	Sections []SectionDetailView `json:"sections"`
}

type PaperHeadView struct {
	Id                      int64   `json:"id"`
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

// PaperHlsExamAdminUpdate 管理端听力 HLS 配置（答题时长以 mock 卷为准，不在此更新）。
type PaperHlsExamAdminUpdate struct {
	AudioHlsPrefix          string
	AudioHlsSegmentCount    int
	AudioHlsSegmentPattern  string
	AudioHlsKeyObject       string
	AudioHlsIvHex           string
	AudioHlsSegmentDuration float64
}

type SectionDetailView struct {
	Id             int64             `json:"id"`
	SortOrder      int               `json:"sort_order"`
	TopicTitle     string            `json:"topic_title"`
	TopicSubtitle  string            `json:"topic_subtitle"`
	TopicType      string            `json:"topic_type"`
	PartCode       int               `json:"part_code"`
	SegmentCode    string            `json:"segment_code"`
	TopicItemsFile string            `json:"topic_items_file"`
	TopicJson      string            `json:"topic_json"`
	Blocks         []BlockDetailView `json:"blocks"`
}

type BlockDetailView struct {
	Id                      int64                `json:"id"`
	BlockOrder              int                  `json:"block_order"`
	GroupIndex              int                  `json:"group_index"`
	QuestionDescriptionJson string               `json:"question_description_json"`
	Questions               []QuestionDetailView `json:"questions"`
}

type QuestionDetailView struct {
	Id                      int64              `json:"id"`
	SortInBlock             int                `json:"sort_in_block"`
	QuestionNo              int                `json:"question_no"`
	Score                   float64            `json:"score"`
	IsExample               int                `json:"is_example"`
	ContentType             string             `json:"content_type"`
	AudioFile               string             `json:"audio_file"`
	StemText                string             `json:"stem_text"`
	ScreenTextJson          string             `json:"screen_text_json"`
	AnalysisJson            string             `json:"analysis_json"`
	QuestionDescriptionJson string             `json:"question_description_json"`
	RawJson                 string             `json:"raw_json"`
	Options                 []OptionDetailView `json:"options"`
}

type OptionDetailView struct {
	Id         int64  `json:"id"`
	Flag       string `json:"flag"`
	SortOrder  int    `json:"sort_order"`
	IsCorrect  int    `json:"is_correct"`
	OptionType string `json:"option_type"`
	Content    string `json:"content"`
}

type PaperDetailForExamInitTree struct {
	Paper    PaperHeadForExamView        `json:"paper"`
	Sections []SectionOutlineForExamView `json:"sections"`
}

type SectionOutlineForExamView struct {
	Id             int64                     `json:"id"`
	SortOrder      int                       `json:"sort_order"`
	TopicTitle     string                    `json:"topic_title"`
	TopicSubtitle  string                    `json:"topic_subtitle"`
	TopicType      string                    `json:"topic_type"`
	PartCode       int                       `json:"part_code"`
	SegmentCode    string                    `json:"segment_code"`
	TopicItemsFile string                    `json:"topic_items_file"`
	TopicJson      string                    `json:"topic_json"`
	Blocks         []BlockOutlineForExamView `json:"blocks"`
}

type BlockOutlineForExamView struct {
	Id                      int64  `json:"id"`
	BlockOrder              int    `json:"block_order"`
	GroupIndex              int    `json:"group_index"`
	QuestionDescriptionJson string `json:"question_description_json"`
	QuestionCount           int    `json:"question_count"`
}

type PaperHeadForExamView struct {
	Id                 int64  `json:"id"`
	Level              string `json:"level"`
	PaperId            string `json:"paper_id"`
	Title              string `json:"title"`
	PrepareTitle       string `json:"prepare_title"`
	PrepareInstruction string `json:"prepare_instruction"`
	PrepareAudioFile   string `json:"prepare_audio_file"`
	SourceBaseUrl      string `json:"source_base_url"`
	AudioHlsPrefix     string `json:"audio_hls_prefix"`
	IndexJson          string `json:"index_json"`
	DurationSeconds    int    `json:"duration_seconds"`
	CreateTime         string `json:"create_time"`
}

type SectionDetailForExamView struct {
	Id             int64                    `json:"id"`
	SortOrder      int                      `json:"sort_order"`
	TopicTitle     string                   `json:"topic_title"`
	TopicSubtitle  string                   `json:"topic_subtitle"`
	TopicType      string                   `json:"topic_type"`
	PartCode       int                      `json:"part_code"`
	SegmentCode    string                   `json:"segment_code"`
	TopicItemsFile string                   `json:"topic_items_file"`
	TopicJson      string                   `json:"topic_json"`
	Blocks         []BlockDetailForExamView `json:"blocks"`
}

type BlockDetailForExamView struct {
	Id                      int64                       `json:"id"`
	BlockOrder              int                         `json:"block_order"`
	GroupIndex              int                         `json:"group_index"`
	QuestionDescriptionJson string                      `json:"question_description_json"`
	Questions               []QuestionDetailForExamView `json:"questions"`
}

type QuestionDetailForExamView struct {
	Id                      int64                     `json:"id"`
	SortInBlock             int                       `json:"sort_in_block"`
	QuestionNo              int                       `json:"question_no"`
	Score                   float64                   `json:"score"`
	IsExample               int                       `json:"is_example"`
	IsSubjective            int                       `json:"is_subjective"`
	ContentType             string                    `json:"content_type"`
	AudioFile               string                    `json:"audio_file"`
	StemText                string                    `json:"stem_text"`
	ScreenTextJson          string                    `json:"screen_text_json"`
	AnalysisJson            string                    `json:"analysis_json"`
	QuestionDescriptionJson string                    `json:"question_description_json"`
	RawJson                 string                    `json:"raw_json"`
	Options                 []OptionDetailForExamView `json:"options"`
}

type OptionDetailForExamView struct {
	Id         int64  `json:"id"`
	Flag       string `json:"flag"`
	SortOrder  int    `json:"sort_order"`
	OptionType string `json:"option_type"`
	Content    string `json:"content"`
}
