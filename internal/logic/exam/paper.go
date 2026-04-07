package exam

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	examdao "exam/internal/dao/exam"
	exam "exam/internal/model/bo/exam"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	"exam/internal/util"
)

// PaperDetail 返回试卷及嵌套大题、题块、小题、选项（只读查看）。
func (s *sExam) PaperDetail(ctx context.Context, examPaperId int64) (*exam.PaperDetailTree, error) {
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

	out := &exam.PaperDetailTree{
		Paper: examPaperEntityToBOHead(paper),
	}

	for _, sec := range sections {
		sv := exam.SectionDetailView{
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
		_ = examdao.ExamQuestionBlock.Ctx(ctx).
			Where("section_id", sec.Id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("block_order").
			Scan(&blocks)

		for _, blk := range blocks {
			bv := exam.BlockDetailView{
				Id:                      blk.Id,
				BlockOrder:              blk.BlockOrder,
				GroupIndex:              blk.GroupIndex,
				QuestionDescriptionJson: blk.QuestionDescriptionJson,
			}

			var qs []examentity.ExamQuestion
			_ = examdao.ExamQuestion.Ctx(ctx).
				Where("block_id", blk.Id).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				OrderAsc("sort_in_block").
				Scan(&qs)

			for _, q := range qs {
				qv := exam.QuestionDetailView{
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
				_ = examdao.ExamOption.Ctx(ctx).
					Where("question_id", q.Id).
					Where("delete_flag", consts.DeleteFlagNotDeleted).
					OrderAsc("sort_order").
					Scan(&opts)

				for _, o := range opts {
					qv.Options = append(qv.Options, exam.OptionDetailView{
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
func (s *sExam) UpdatePaperSettings(ctx context.Context, examPaperId int64, in exam.PaperHlsExamAdminUpdate, updater string) error {
	if examPaperId <= 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	if in.AudioHlsSegmentCount < 0 || in.AudioHlsSegmentDuration < 0 {
		return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
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

func examPaperEntityToBOHead(p examentity.ExamPaper) exam.PaperHeadView {
	v := exam.PaperHeadView{
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
	v.CreateTime = util.ToRFC3339UTC(p.CreateTime)
	return v
}
