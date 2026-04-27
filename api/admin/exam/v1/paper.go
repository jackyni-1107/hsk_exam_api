package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PaperListReq struct {
	g.Meta      `path:"/exam/paper/list" method:"get" tags:"试卷管理" summary:"exam_paper 试卷列表" permission:"exam:paper:list"`
	Level       string `json:"level" dc:"级别字符串筛选（兼容）；优先 mock_level_id"`
	MockLevelId int64  `json:"mock_level_id" dc:"mock_levels.id，按 mock 卷 level_id 筛，0 不限"`
	Page        int    `json:"page" dc:"页码"`
	Size        int    `json:"size" dc:"每页条数"`
}

type PaperListRes struct {
	List  []*PaperListItem `json:"list" dc:"列表"`
	Total int              `json:"total" dc:"总数"`
}

type PaperListItem struct {
	Id                      int64   `json:"id" dc:"exam_paper.id"`
	MockExaminationPaperId  int64   `json:"mock_examination_paper_id" dc:"mock_examination_paper.id"`
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
	g.Meta                 `path:"/exam/paper/import" method:"post" tags:"试卷管理" summary:"导入" permission:"exam:paper:import"`
	MockExaminationPaperId int64  `json:"mock_examination_paper_id" v:"required|min:1#err.invalid_params" dc:"mock 卷 id，业务主键"`
	Title                  string `json:"title" dc:"可选，试卷名称；填写则写入 exam_paper.title，否则用 mock 卷 name"`
	AudioHlsPrefix         string `json:"audio_hls_prefix" dc:"听力 HLS 目录前缀（无首尾/），动态 m3u8 拼接时使用"`
	// fail：已存在则返回 conflict；overwrite / new：伪删旧题目树后在原 exam_paper 上全量更新（二者实现一致，见 docs）
	ConflictMode string `json:"conflict_mode" dc:"fail|overwrite|new，默认 fail；兼容 new_copy 等同于 new"`
}

type PaperImportRes struct {
	ExaminationPaperId         int64 `json:"examination_paper_id" dc:"mock_examination_paper.id"`
	Conflict                   bool  `json:"conflict" dc:"true 表示未导入，因已存在且 conflict_mode=fail"`
	ExistingExaminationPaperId int64 `json:"existing_examination_paper_id" dc:"冲突时已存在的 mock 卷 id"`
	SectionCount               int   `json:"section_count" dc:"大题数"`
	QuestionCount              int   `json:"question_count" dc:"试题数"`
}

type PaperUpdateReq struct {
	g.Meta `path:"/exam/paper/update" method:"post" tags:"试卷管理" summary:"修改听力 HLS 配置；答题时长以 exam_paper.duration_seconds 为准。id 为 exam_paper.id"`
	Id     int64 `json:"id" v:"required|min:1#err.invalid_params" dc:"exam_paper.id"`
	// 听力 HLS：与 exam_paper 表字段一致
	AudioHlsPrefix          string  `json:"audio_hls_prefix" dc:"目录前缀（无首尾/）"`
	AudioHlsSegmentCount    int     `json:"audio_hls_segment_count" v:"min:0#err.invalid_params" dc:"分片总数，0 表示未配置"`
	AudioHlsSegmentPattern  string  `json:"audio_hls_segment_pattern" dc:"分片文件名 fmt，空则默认 %%05d.ts"`
	AudioHlsKeyObject       string  `json:"audio_hls_key_object" dc:"密钥对象相对 prefix 的路径"`
	AudioHlsIvHex           string  `json:"audio_hls_iv_hex" dc:"AES-128 IV 十六进制"`
	AudioHlsSegmentDuration float64 `json:"audio_hls_segment_duration" v:"min:0#err.invalid_params" dc:"#EXTINF 时长秒"`
}

type PaperUpdateRes struct{}

type PaperEditReq struct {
	g.Meta `path:"/exam/paper/edit" method:"post" tags:"试卷管理" summary:"编辑试卷元数据（exam_paper）"`
	// exam_paper.id，与列表项 id 一致
	ExamPaperId int64  `json:"exam_paper_id" v:"required|min:1#err.invalid_params" dc:"exam_paper.id"`
	Title       string `json:"title" dc:"试卷标题"`
	// 考前：与 index prepare 对应
	PrepareTitle       string `json:"prepare_title" dc:"考前标题"`
	PrepareInstruction string `json:"prepare_instruction" dc:"考前说明"`
	PrepareAudioFile   string `json:"prepare_audio_file" dc:"考前音频"`
	SourceBaseUrl      string `json:"source_base_url" dc:"资源基址，建议以 / 结尾"`
	// 0 表示使用系统默认 exam.defaultDurationSeconds
	DurationSeconds int `json:"duration_seconds" v:"min:0#err.invalid_params" dc:"考试时长秒"`
}

type PaperEditRes struct{}

// PaperPurgeReq 物理删除 exam_paper（不可恢复）。confirm_text 须与「DELETE:{exam_paper_id}」完全一致（含大小写与冒号）。
type PaperPurgeReq struct {
	g.Meta      `path:"/exam/paper/purge" method:"post" tags:"试卷管理" summary:"物理删除试卷" permission:"exam:paper:purge"`
	ExamPaperId int64  `json:"exam_paper_id" v:"required|min:1#err.invalid_params" dc:"exam_paper.id"`
	ConfirmText string `json:"confirm_text" v:"required#err.invalid_params" dc:"二次确认，须为 DELETE:试卷主键"`
}

type PaperPurgeRes struct{}

type PaperDetailReq struct {
	g.Meta `path:"/exam/paper/{id}" method:"get" tags:"试卷管理" summary:"试卷详情（大题/题块/试题）；path id 为 exam_paper.id，未找到 code=11114"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"exam_paper.id"`
}

type PaperSectionDetailReq struct {
	g.Meta    `path:"/exam/paper/{id}/sections/{sectionId}" method:"get" tags:"试卷管理" summary:"单大题详情（与试卷详情中一节一致，含选项正误）；id 为 exam_paper.id"`
	Id        int64 `json:"id" in:"path" v:"required|min:1#err.invalid_params" dc:"exam_paper.id"`
	SectionId int64 `json:"sectionId" in:"path" v:"required|min:1#err.invalid_params" dc:"大题 exam_section.id"`
}

type PaperSectionDetailRes struct {
	Section PaperDetailSection `json:"section" dc:"大题（题块/试题/选项）"`
}

type PaperDetailRes struct {
	Paper    PaperDetailPaper     `json:"paper" dc:"试卷信息"`
	Sections []PaperDetailSection `json:"sections" dc:"大题列表"`
}

type PaperDetailPaper struct {
	ExamPaperId             int64   `json:"exam_paper_id" dc:"exam_paper.id"`
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
	DurationSeconds         int     `json:"duration_seconds" dc:"考试时长秒，0 为系统默认"`
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
