package paper

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	"exam/internal/dao"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"
	"exam/internal/utility/exampaper"
)

func (s *sPaper) PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (map[string]interface{}, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	var sec examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("id", sectionId).
		Where("exam_paper_id", paper.Id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&sec); err != nil {
		return nil, err
	}
	if sec.Id == 0 {
		return nil, gerror.NewCode(consts.CodeExamSectionNotFound)
	}
	if sec.TopicJson == "" {
		return nil, gerror.NewCode(consts.CodeExamSectionTopicEmpty)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(sec.TopicJson), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *sPaper) PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	t, err := PaperDetailForExamInit(ctx, paper.Id)
	if err != nil {
		return nil, err
	}
	out := paperDetailForExamInitTreeToBO(t)
	return &out, nil
}

func (s *sPaper) PaperPrepareSegments(ctx context.Context, mockPaperID int64) ([]exambo.PaperPrepareSegment, error) {
	var mockPaper mockentity.MockExaminationPaper
	if err := dao.MockExaminationPaper.Ctx(ctx).
		Where("id", mockPaperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&mockPaper); err != nil {
		return nil, err
	}
	if mockPaper.Id == 0 {
		return []exambo.PaperPrepareSegment{}, nil
	}

	segments := make([]mockentity.MockExaminationSegment, 0)
	if err := dao.MockExaminationSegment.Ctx(ctx).
		Where("level_id", mockPaper.LevelId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("seq,id").
		Scan(&segments); err != nil {
		return nil, err
	}
	if len(segments) == 0 {
		return []exambo.PaperPrepareSegment{}, nil
	}

	segmentIDs := make([]int64, 0, len(segments))
	partsBySegment := make(map[int64][]mockentity.MockExaminationPart, len(segments))
	for _, seg := range segments {
		segmentIDs = append(segmentIDs, seg.Id)
	}

	parts := make([]mockentity.MockExaminationPart, 0)
	if err := dao.MockExaminationPart.Ctx(ctx).
		WhereIn("segment_id", segmentIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("segment_id,code,id").
		Scan(&parts); err != nil {
		return nil, err
	}
	for _, part := range parts {
		partsBySegment[part.SegmentId] = append(partsBySegment[part.SegmentId], part)
	}

	out := make([]exambo.PaperPrepareSegment, 0, len(segments))
	for _, seg := range segments {
		boSeg := exambo.PaperPrepareSegment{
			SegmentCode:   seg.SegmentCode,
			TotalScore:    seg.ScoreFull,
			QuestionCount: seg.QuestionCount,
			Duration:      seg.Duration,
			Seq:           seg.Seq,
			SegmentDesc:   seg.SegmentDesc,
			Parts:         make([]exambo.PaperPreparePartItem, 0, len(partsBySegment[seg.Id])),
		}
		for _, p := range partsBySegment[seg.Id] {
			partRate := 0.0
			if seg.ScoreFull > 0 {
				partRate = p.PartScore / float64(seg.ScoreFull)
			}
			boSeg.Parts = append(boSeg.Parts, exambo.PaperPreparePartItem{
				PartCode:                p.Code,
				PartName:                p.PartName,
				PartNameTrans:           p.PartNameTrans,
				PartRate:                partRate,
				PartScore:               p.PartScore,
				QuestionCount:           p.QuestionCount,
				ObjectiveQuestionCount:  p.ObjectiveQuestionCount,
				SubjectiveQuestionCount: p.SubjectiveQuestionCount,
				PartAnswerTime:          p.AnswerTime,
				ScoreTotal:              0,
				CorrectCount:            0,
				CorrectRate:             0,
				Practiced:               false,
				QuestionType:            "",
			})
		}
		out = append(out, boSeg)
	}
	return out, nil
}

func (s *sPaper) PaperSectionDetailForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*exambo.SectionDetailForExamView, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	t, err := PaperSectionDetailForExam(ctx, paper.Id, sectionId)
	if err != nil {
		return nil, err
	}
	out := sectionDetailForExamViewToBO(t)
	return &out, nil
}

func paperDetailForExamInitTreeToBO(t *PaperDetailForExamInitTree) exambo.PaperDetailForExamInitTree {
	if t == nil {
		return exambo.PaperDetailForExamInitTree{}
	}
	out := exambo.PaperDetailForExamInitTree{
		Paper: exambo.PaperHeadForExamView{
			Id:                 t.Paper.Id,
			Level:              t.Paper.Level,
			PaperId:            t.Paper.PaperId,
			Title:              t.Paper.Title,
			PrepareInstruction: t.Paper.PrepareInstruction,
			PrepareAudioFile:   t.Paper.PrepareAudioFile,
			SourceBaseUrl:      t.Paper.SourceBaseUrl,
			IndexJson:          t.Paper.IndexJson,
			DurationSeconds:    t.Paper.DurationSeconds,
			CreateTime:         t.Paper.CreateTime,
		},
		Sections: make([]exambo.SectionOutlineForExamView, len(t.Sections)),
	}
	for i, sec := range t.Sections {
		blocks := make([]exambo.BlockOutlineForExamView, len(sec.Blocks))
		for j, b := range sec.Blocks {
			blocks[j] = exambo.BlockOutlineForExamView{
				Id:                      b.Id,
				BlockOrder:              b.BlockOrder,
				GroupIndex:              b.GroupIndex,
				QuestionDescriptionJson: b.QuestionDescriptionJson,
				QuestionCount:           b.QuestionCount,
			}
		}
		out.Sections[i] = exambo.SectionOutlineForExamView{
			Id:             sec.Id,
			SortOrder:      sec.SortOrder,
			TopicTitle:     sec.TopicTitle,
			TopicSubtitle:  sec.TopicSubtitle,
			TopicType:      sec.TopicType,
			PartCode:       sec.PartCode,
			SegmentCode:    sec.SegmentCode,
			TopicItemsFile: sec.TopicItemsFile,
			TopicJson:      sec.TopicJson,
			Blocks:         blocks,
		}
	}
	return out
}

func sectionDetailForExamViewToBO(v *SectionDetailForExamView) exambo.SectionDetailForExamView {
	if v == nil {
		return exambo.SectionDetailForExamView{}
	}
	out := exambo.SectionDetailForExamView{
		Id:             v.Id,
		SortOrder:      v.SortOrder,
		TopicTitle:     v.TopicTitle,
		TopicSubtitle:  v.TopicSubtitle,
		TopicType:      v.TopicType,
		PartCode:       v.PartCode,
		SegmentCode:    v.SegmentCode,
		TopicItemsFile: v.TopicItemsFile,
		TopicJson:      v.TopicJson,
		Blocks:         make([]exambo.BlockDetailForExamView, len(v.Blocks)),
	}
	for i, b := range v.Blocks {
		out.Blocks[i] = blockDetailForExamViewToBO(b)
	}
	return out
}

func blockDetailForExamViewToBO(b BlockDetailForExamView) exambo.BlockDetailForExamView {
	qs := make([]exambo.QuestionDetailForExamView, len(b.Questions))
	for i, q := range b.Questions {
		qs[i] = questionDetailForExamViewToBO(q)
	}
	return exambo.BlockDetailForExamView{
		Id:                      b.Id,
		BlockOrder:              b.BlockOrder,
		GroupIndex:              b.GroupIndex,
		QuestionDescriptionJson: b.QuestionDescriptionJson,
		Questions:               qs,
	}
}

func questionDetailForExamViewToBO(q QuestionDetailForExamView) exambo.QuestionDetailForExamView {
	opts := make([]exambo.OptionDetailForExamView, len(q.Options))
	for i, o := range q.Options {
		opts[i] = exambo.OptionDetailForExamView{
			Id:         o.Id,
			Flag:       o.Flag,
			SortOrder:  o.SortOrder,
			OptionType: o.OptionType,
			Content:    o.Content,
		}
	}
	return exambo.QuestionDetailForExamView{
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
		Options:                 opts,
	}
}
