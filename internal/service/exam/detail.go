package exam

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	"exam/internal/dao"
	examentity "exam/internal/model/entity/exam"
)

// PaperDetailForExam 客户端考前拉题：不含选项正误，含 is_subjective。
func PaperDetailForExam(ctx context.Context, examPaperId int64) (*PaperDetailForExamTree, error) {
	var paper examentity.ExamPaper
	err := dao.ExamPaper.Ctx(ctx).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper)
	if err != nil {
		return nil, err
	}
	if paper.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.exam_paper_not_found")
	}

	var sections []examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&sections); err != nil {
		return nil, err
	}

	out := &PaperDetailForExamTree{
		Paper: paperHeadForExam(paper),
	}

	for _, sec := range sections {
		sv := SectionDetailForExamView{
			Id:             sec.Id,
			SortOrder:      sec.SortOrder,
			TopicTitle:     sec.TopicTitle,
			TopicSubtitle:  sec.TopicSubtitle,
			TopicType:      sec.TopicType,
			PartCode:       sec.PartCode,
			SegmentCode:    sec.SegmentCode,
			TopicItemsFile: sec.TopicItemsFile,
			TopicJson:      sec.TopicJson,
		}

		var blocks []examentity.ExamQuestionBlock
		_ = dao.ExamQuestionBlock.Ctx(ctx).
			Where("section_id", sec.Id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("block_order").
			Scan(&blocks)

		for _, blk := range blocks {
			bv := BlockDetailForExamView{
				Id:                      blk.Id,
				BlockOrder:              blk.BlockOrder,
				GroupIndex:              blk.GroupIndex,
				QuestionDescriptionJson: blk.QuestionDescriptionJson,
			}

			var qs []examentity.ExamQuestion
			_ = dao.ExamQuestion.Ctx(ctx).
				Where("block_id", blk.Id).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				OrderAsc("sort_in_block").
				Scan(&qs)

			for _, q := range qs {
				qv := QuestionDetailForExamView{
					Id:                      q.Id,
					SortInBlock:             q.SortInBlock,
					QuestionNo:              q.QuestionNo,
					Score:                   q.Score,
					IsExample:               q.IsExample,
					IsSubjective:            q.IsSubjective,
					ContentType:             q.ContentType,
					AudioFile:               q.AudioFile,
					StemText:                q.StemText,
					ScreenTextJson:          q.ScreenTextJson,
					AnalysisJson:            q.AnalysisJson,
					QuestionDescriptionJson: q.QuestionDescriptionJson,
					RawJson:                 q.RawJson,
				}

				var opts []examentity.ExamOption
				_ = dao.ExamOption.Ctx(ctx).
					Where("question_id", q.Id).
					Where("delete_flag", consts.DeleteFlagNotDeleted).
					OrderAsc("sort_order").
					Scan(&opts)

				for _, o := range opts {
					qv.Options = append(qv.Options, OptionDetailForExamView{
						Id:         o.Id,
						Flag:       o.Flag,
						SortOrder:  o.SortOrder,
						OptionType: o.OptionType,
						Content:    o.Content,
					})
				}
				bv.Questions = append(bv.Questions, qv)
			}
			sv.Blocks = append(sv.Blocks, bv)
		}
		out.Sections = append(out.Sections, sv)
	}

	return out, nil
}

func paperHeadForExam(p examentity.ExamPaper) PaperHeadForExamView {
	v := PaperHeadForExamView{
		Id:                 p.MockExaminationPaperId,
		Level:              p.Level,
		PaperId:            p.PaperId,
		Title:              p.Title,
		PrepareInstruction: p.PrepareInstruction,
		PrepareAudioFile:   p.PrepareAudioFile,
		SourceBaseUrl:      p.SourceBaseUrl,
		IndexJson:          p.IndexJson,
		DurationSeconds:    p.DurationSeconds,
	}
	if p.CreateTime != nil {
		v.CreateTime = p.CreateTime.String()
	}
	return v
}

// PaperDetailForExamTree 客户端试卷（脱敏）。
type PaperDetailForExamTree struct {
	Paper    PaperHeadForExamView       `json:"paper"`
	Sections []SectionDetailForExamView `json:"sections"`
}

// PaperHeadForExamView 含考试时长。
type PaperHeadForExamView struct {
	Id                 int64  `json:"id"`
	Level              string `json:"level"`
	PaperId            string `json:"paper_id"`
	Title              string `json:"title"`
	PrepareInstruction string `json:"prepare_instruction"`
	PrepareAudioFile   string `json:"prepare_audio_file"`
	SourceBaseUrl      string `json:"source_base_url"`
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

// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
func PaperDetail(ctx context.Context, examPaperId int64) (*PaperDetailTree, error) {
	var paper examentity.ExamPaper
	err := dao.ExamPaper.Ctx(ctx).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper)
	if err != nil {
		return nil, err
	}
	if paper.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.exam_paper_not_found")
	}

	var sections []examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&sections); err != nil {
		return nil, err
	}

	out := &PaperDetailTree{
		Paper: paperToView(paper),
	}

	for _, sec := range sections {
		sv := SectionDetailView{
			Id:             sec.Id,
			SortOrder:      sec.SortOrder,
			TopicTitle:     sec.TopicTitle,
			TopicSubtitle:  sec.TopicSubtitle,
			TopicType:      sec.TopicType,
			PartCode:       sec.PartCode,
			SegmentCode:    sec.SegmentCode,
			TopicItemsFile: sec.TopicItemsFile,
			TopicJson:      sec.TopicJson,
		}

		var blocks []examentity.ExamQuestionBlock
		_ = dao.ExamQuestionBlock.Ctx(ctx).
			Where("section_id", sec.Id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("block_order").
			Scan(&blocks)

		for _, blk := range blocks {
			bv := BlockDetailView{
				Id:                      blk.Id,
				BlockOrder:              blk.BlockOrder,
				GroupIndex:              blk.GroupIndex,
				QuestionDescriptionJson: blk.QuestionDescriptionJson,
			}

			var qs []examentity.ExamQuestion
			_ = dao.ExamQuestion.Ctx(ctx).
				Where("block_id", blk.Id).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				OrderAsc("sort_in_block").
				Scan(&qs)

			for _, q := range qs {
				qv := QuestionDetailView{
					Id:                      q.Id,
					SortInBlock:             q.SortInBlock,
					QuestionNo:              q.QuestionNo,
					Score:                   q.Score,
					IsExample:               q.IsExample,
					ContentType:             q.ContentType,
					AudioFile:               q.AudioFile,
					StemText:                q.StemText,
					ScreenTextJson:          q.ScreenTextJson,
					AnalysisJson:            q.AnalysisJson,
					QuestionDescriptionJson: q.QuestionDescriptionJson,
					RawJson:                 q.RawJson,
				}

				var opts []examentity.ExamOption
				_ = dao.ExamOption.Ctx(ctx).
					Where("question_id", q.Id).
					Where("delete_flag", consts.DeleteFlagNotDeleted).
					OrderAsc("sort_order").
					Scan(&opts)

				for _, o := range opts {
					qv.Options = append(qv.Options, OptionDetailView{
						Id:         o.Id,
						Flag:       o.Flag,
						SortOrder:  o.SortOrder,
						IsCorrect:  o.IsCorrect,
						OptionType: o.OptionType,
						Content:    o.Content,
					})
				}
				bv.Questions = append(bv.Questions, qv)
			}
			sv.Blocks = append(sv.Blocks, bv)
		}
		out.Sections = append(out.Sections, sv)
	}

	return out, nil
}

// PaperDetailTree 试卷详情（树形）。
type PaperDetailTree struct {
	Paper    PaperHeadView       `json:"paper"`
	Sections []SectionDetailView `json:"sections"`
}

// PaperHeadView 试卷主表字段。
type PaperHeadView struct {
	Id                 int64  `json:"id"`
	Level              string `json:"level"`
	PaperId            string `json:"paper_id"`
	Title              string `json:"title"`
	PrepareInstruction string `json:"prepare_instruction"`
	PrepareAudioFile   string `json:"prepare_audio_file"`
	SourceBaseUrl      string `json:"source_base_url"`
	IndexJson          string `json:"index_json"`
	CreateTime         string `json:"create_time"`
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

func paperToView(p examentity.ExamPaper) PaperHeadView {
	v := PaperHeadView{
		Id:                 p.Id,
		Level:              p.Level,
		PaperId:            p.PaperId,
		Title:              p.Title,
		PrepareInstruction: p.PrepareInstruction,
		PrepareAudioFile:   p.PrepareAudioFile,
		SourceBaseUrl:      p.SourceBaseUrl,
		IndexJson:          p.IndexJson,
	}
	if p.CreateTime != nil {
		v.CreateTime = p.CreateTime.String()
	}
	return v
}
