package exam

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/entity"
	"exam/internal/util"
)

// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
func PaperDetail(ctx context.Context, examPaperId int64) (*PaperDetailTree, error) {
	var paper entity.ExamPaper
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

	var sections []entity.ExamSection
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

		var blocks []entity.ExamQuestionBlock
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

			var qs []entity.ExamQuestion
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

				var opts []entity.ExamOption
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

func paperToView(p entity.ExamPaper) PaperHeadView {
	v := PaperHeadView{
		Id:                 p.MockExaminationPaperId,
		Level:              p.Level,
		PaperId:            p.PaperId,
		Title:              p.Title,
		PrepareInstruction: p.PrepareInstruction,
		PrepareAudioFile:   p.PrepareAudioFile,
		SourceBaseUrl:      p.SourceBaseUrl,
		IndexJson:          p.IndexJson,
	}
	v.CreateTime = util.ToRFC3339UTC(p.CreateTime)
	return v
}
