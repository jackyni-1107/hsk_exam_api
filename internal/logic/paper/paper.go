package paper

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	examdao "exam/internal/dao/exam"
	exambo "exam/internal/model/bo/exam"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/utility"
)

// PaperDetail 返回试卷及嵌套的大题、题块、题目、选项。
func (s *sPaper) PaperDetail(ctx context.Context, examPaperId int64) (*exambo.PaperDetailTree, error) {
	data, err := loadPaperDetailData(ctx, examPaperId)
	if err != nil {
		return nil, err
	}

	out := &exambo.PaperDetailTree{
		Paper: examPaperEntityToBOHead(data.paper, data.mockPaper.Name),
	}

	for _, sec := range data.sections {
		out.Sections = append(out.Sections, sectionDetailViewFromData(
			sec, data.blocksBySection[sec.Id], data.questionsByBlock, data.optionsByQuestion))
	}

	return out, nil
}

// sectionDetailViewFromData 将一个大题下的题块、题目、选项拼装为详情视图。
func sectionDetailViewFromData(
	sec examentity.ExamSection,
	blocks []examentity.ExamQuestionBlock,
	questionsByBlock map[int64][]examentity.ExamQuestion,
	optionsByQuestion map[int64][]examentity.ExamOption,
) exambo.SectionDetailView {
	sv := exambo.SectionDetailView{
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
	for _, blk := range blocks {
		bv := exambo.BlockDetailView{
			Id:                      blk.Id,
			BlockOrder:              blk.BlockOrder,
			GroupIndex:              blk.GroupIndex,
			QuestionDescriptionJson: blk.QuestionDescriptionJson,
		}
		for _, q := range questionsByBlock[blk.Id] {
			qv := exambo.QuestionDetailView{
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
			for _, o := range optionsByQuestion[q.Id] {
				qv.Options = append(qv.Options, exambo.OptionDetailView{
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
	return sv
}

// PaperDetailSection 只加载单个大题的完整详情。
func (s *sPaper) PaperDetailSection(ctx context.Context, examPaperId, sectionId int64) (*exambo.SectionDetailView, error) {
	data, err := loadPaperSectionDetailData(ctx, examPaperId, sectionId)
	if err != nil {
		return nil, err
	}
	sv := sectionDetailViewFromData(data.section, data.blocksBySection[data.section.Id], data.questionsByBlock, data.optionsByQuestion)
	return &sv, nil
}

// UpdatePaperSettings 修改试卷的 HLS 配置。
func (s *sPaper) UpdatePaperSettings(ctx context.Context, examPaperId int64, in exambo.PaperHlsExamAdminUpdate, updater string) error {
	if examPaperId <= 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	if in.AudioHlsSegmentCount < 0 || in.AudioHlsSegmentDuration < 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	var paper examentity.ExamPaper
	if err := examdao.ExamPaper.Ctx(ctx).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper); err != nil {
		return err
	}
	if paper.Id == 0 {
		return gerror.NewCode(consts.CodeExamPaperNotFound)
	}
	prefix := strings.Trim(in.AudioHlsPrefix, "/")
	_, err := examdao.ExamPaper.Ctx(ctx).
		Data(examdo.ExamPaper{
			AudioHlsPrefix:          prefix,
			AudioHlsSegmentCount:    in.AudioHlsSegmentCount,
			AudioHlsSegmentPattern:  in.AudioHlsSegmentPattern,
			AudioHlsKeyObject:       in.AudioHlsKeyObject,
			AudioHlsIvHex:           in.AudioHlsIvHex,
			AudioHlsSegmentDuration: in.AudioHlsSegmentDuration,
			Updater:                 updater,
			UpdateTime:              gtime.Now(),
		}).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Update()
	if err != nil {
		return err
	}
	invalidatePaperCaches(ctx, examPaperId, paper.MockExaminationPaperId)
	return nil
}

// UpdatePaperMeta 修改试卷元数据，不包含 HLS 与题目树。
func (s *sPaper) UpdatePaperMeta(ctx context.Context, examPaperId int64, in exambo.PaperMetaAdminUpdate, updater string) error {
	if examPaperId <= 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	if in.DurationSeconds < 0 {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	var paper examentity.ExamPaper
	if err := examdao.ExamPaper.Ctx(ctx).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper); err != nil {
		return err
	}
	if paper.Id == 0 {
		return gerror.NewCode(consts.CodeExamPaperNotFound)
	}
	baseURL := strings.TrimSpace(in.SourceBaseURL)
	if baseURL != "" {
		baseURL = strings.TrimRight(baseURL, "/") + "/"
	}
	_, err := examdao.ExamPaper.Ctx(ctx).
		Data(examdo.ExamPaper{
			Title:              in.Title,
			PrepareTitle:       in.PrepareTitle,
			PrepareInstruction: in.PrepareInstruction,
			PrepareAudioFile:   in.PrepareAudioFile,
			SourceBaseUrl:      baseURL,
			DurationSeconds:    in.DurationSeconds,
			Updater:            updater,
			UpdateTime:         gtime.Now(),
		}).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Update()
	if err != nil {
		return err
	}
	invalidatePaperCaches(ctx, examPaperId, paper.MockExaminationPaperId)
	return nil
}

func examPaperEntityToBOHead(p examentity.ExamPaper, name string) exambo.PaperHeadView {
	v := exambo.PaperHeadView{
		ExamPaperId:             p.Id,
		Id:                      p.MockExaminationPaperId,
		Level:                   p.Level,
		PaperId:                 p.PaperId,
		Title:                   p.Title,
		Name:                    name,
		PrepareTitle:            p.PrepareTitle,
		PrepareInstruction:      p.PrepareInstruction,
		PrepareAudioFile:        p.PrepareAudioFile,
		SourceBaseUrl:           p.SourceBaseUrl,
		AudioHlsPrefix:          p.AudioHlsPrefix,
		AudioHlsSegmentCount:    p.AudioHlsSegmentCount,
		AudioHlsSegmentPattern:  p.AudioHlsSegmentPattern,
		AudioHlsKeyObject:       p.AudioHlsKeyObject,
		AudioHlsIvHex:           p.AudioHlsIvHex,
		AudioHlsSegmentDuration: p.AudioHlsSegmentDuration,
		IndexJson:               p.IndexJson,
		DurationSeconds:         p.DurationSeconds,
	}
	v.CreateTime = utility.ToRFC3339UTC(p.CreateTime)
	return v
}
