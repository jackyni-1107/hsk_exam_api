package exam

import (
	"context"
	"sort"
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

// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
func (s *sExam) PaperDetail(ctx context.Context, examPaperId int64) (*exambo.PaperDetailTree, error) {
	var paper examentity.ExamPaper
	err := examdao.ExamPaper.Ctx(ctx).
		Where("id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper)
	if err != nil {
		return nil, err
	}
	if paper.Id == 0 {
		return nil, gerror.NewCode(consts.CodeExamPaperNotFound)
	}

	var sections []examentity.ExamSection
	if err := examdao.ExamSection.Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&sections); err != nil {
		return nil, err
	}

	sectionIDs := make([]interface{}, 0, len(sections))
	for _, sec := range sections {
		sectionIDs = append(sectionIDs, sec.Id)
	}

	var allBlocks []examentity.ExamQuestionBlock
	if len(sectionIDs) > 0 {
		if err := examdao.ExamQuestionBlock.Ctx(ctx).
			WhereIn("section_id", sectionIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&allBlocks); err != nil {
			return nil, err
		}
	}
	blocksBySection := make(map[int64][]examentity.ExamQuestionBlock, len(sectionIDs))
	for _, blk := range allBlocks {
		blocksBySection[blk.SectionId] = append(blocksBySection[blk.SectionId], blk)
	}
	for sid, blks := range blocksBySection {
		sort.Slice(blks, func(i, j int) bool {
			if blks[i].BlockOrder != blks[j].BlockOrder {
				return blks[i].BlockOrder < blks[j].BlockOrder
			}
			return blks[i].Id < blks[j].Id
		})
		blocksBySection[sid] = blks
	}

	blockIDs := make([]interface{}, 0, len(allBlocks))
	for _, blk := range allBlocks {
		blockIDs = append(blockIDs, blk.Id)
	}

	var allQuestions []examentity.ExamQuestion
	if len(blockIDs) > 0 {
		if err := examdao.ExamQuestion.Ctx(ctx).
			WhereIn("block_id", blockIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&allQuestions); err != nil {
			return nil, err
		}
	}
	questionsByBlock := make(map[int64][]examentity.ExamQuestion, len(blockIDs))
	for _, q := range allQuestions {
		questionsByBlock[q.BlockId] = append(questionsByBlock[q.BlockId], q)
	}
	for bid, qs := range questionsByBlock {
		sort.Slice(qs, func(i, j int) bool {
			if qs[i].SortInBlock != qs[j].SortInBlock {
				return qs[i].SortInBlock < qs[j].SortInBlock
			}
			return qs[i].Id < qs[j].Id
		})
		questionsByBlock[bid] = qs
	}

	qIDs := make([]interface{}, 0, len(allQuestions))
	for _, q := range allQuestions {
		qIDs = append(qIDs, q.Id)
	}

	var allOptions []examentity.ExamOption
	if len(qIDs) > 0 {
		if err := examdao.ExamOption.Ctx(ctx).
			WhereIn("question_id", qIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&allOptions); err != nil {
			return nil, err
		}
	}
	optionsByQuestion := make(map[int64][]examentity.ExamOption, len(qIDs))
	for _, o := range allOptions {
		optionsByQuestion[o.QuestionId] = append(optionsByQuestion[o.QuestionId], o)
	}
	for qid, opts := range optionsByQuestion {
		sort.Slice(opts, func(i, j int) bool {
			if opts[i].SortOrder != opts[j].SortOrder {
				return opts[i].SortOrder < opts[j].SortOrder
			}
			return opts[i].Id < opts[j].Id
		})
		optionsByQuestion[qid] = opts
	}

	out := &exambo.PaperDetailTree{
		Paper: examPaperEntityToBOHead(paper),
	}

	for _, sec := range sections {
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
		for _, blk := range blocksBySection[sec.Id] {
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
		out.Sections = append(out.Sections, sv)
	}

	return out, nil
}

// UpdatePaperSettings 修改试卷听力 HLS 配置（答题时长以 mock_examination_paper 为准）。
func (s *sExam) UpdatePaperSettings(ctx context.Context, examPaperId int64, in exambo.PaperHlsExamAdminUpdate, updater string) error {
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
	s.InvalidatePaperForExamCache(ctx, examPaperId)
	return nil
}

func examPaperEntityToBOHead(p examentity.ExamPaper) exambo.PaperHeadView {
	v := exambo.PaperHeadView{
		Id:                      p.MockExaminationPaperId,
		Level:                   p.Level,
		PaperId:                 p.PaperId,
		Title:                   p.Title,
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
	}
	v.CreateTime = utility.ToRFC3339UTC(p.CreateTime)
	return v
}
