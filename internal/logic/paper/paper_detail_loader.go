package paper

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	examdao "exam/internal/dao/exam"
	mockdao "exam/internal/dao/mock"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"
)

type paperDetailData struct {
	paper             examentity.ExamPaper
	mockPaper         mockentity.MockExaminationPaper
	section           examentity.ExamSection
	sections          []examentity.ExamSection
	blocksBySection   map[int64][]examentity.ExamQuestionBlock
	questionsByBlock  map[int64][]examentity.ExamQuestion
	optionsByQuestion map[int64][]examentity.ExamOption
}

func loadPaperDetailData(ctx context.Context, examPaperID int64) (*paperDetailData, error) {
	paper, err := loadExamPaperByID(ctx, examPaperID)
	if err != nil {
		return nil, err
	}
	mockPaper, err := loadMockPaperForPaper(ctx, paper.MockExaminationPaperId)
	if err != nil {
		return nil, err
	}
	sections, err := loadPaperSections(ctx, examPaperID)
	if err != nil {
		return nil, err
	}
	blocksBySection, questionsByBlock, optionsByQuestion, err := loadPaperQuestionGraph(ctx, sections, 0)
	if err != nil {
		return nil, err
	}
	return &paperDetailData{
		paper:             paper,
		mockPaper:         mockPaper,
		sections:          sections,
		blocksBySection:   blocksBySection,
		questionsByBlock:  questionsByBlock,
		optionsByQuestion: optionsByQuestion,
	}, nil
}

func loadPaperSectionDetailData(ctx context.Context, examPaperID, sectionID int64) (*paperDetailData, error) {
	paper, err := loadExamPaperByID(ctx, examPaperID)
	if err != nil {
		return nil, err
	}
	section, err := loadPaperSectionByID(ctx, examPaperID, sectionID)
	if err != nil {
		return nil, err
	}
	sections := []examentity.ExamSection{section}
	blocksBySection, questionsByBlock, optionsByQuestion, err := loadPaperQuestionGraph(ctx, sections, sectionID)
	if err != nil {
		return nil, err
	}
	return &paperDetailData{
		paper:             paper,
		section:           section,
		sections:          sections,
		blocksBySection:   blocksBySection,
		questionsByBlock:  questionsByBlock,
		optionsByQuestion: optionsByQuestion,
	}, nil
}

func loadExamPaperByID(ctx context.Context, examPaperID int64) (examentity.ExamPaper, error) {
	var paper examentity.ExamPaper
	if err := examdao.ExamPaper.Ctx(ctx).
		Where("id", examPaperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper); err != nil {
		return paper, err
	}
	if paper.Id == 0 {
		return paper, gerror.NewCode(consts.CodeExamPaperNotFound)
	}
	return paper, nil
}

func loadMockPaperForPaper(ctx context.Context, mockPaperID int64) (mockentity.MockExaminationPaper, error) {
	var mockPaper mockentity.MockExaminationPaper
	if err := mockdao.MockExaminationPaper.Ctx(ctx).
		Fields(mockdao.MockExaminationPaper.Columns().Id, mockdao.MockExaminationPaper.Columns().Name).
		Where(mockdao.MockExaminationPaper.Columns().Id, mockPaperID).
		Where(mockdao.MockExaminationPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Scan(&mockPaper); err != nil {
		return mockPaper, err
	}
	return mockPaper, nil
}

func loadPaperSections(ctx context.Context, examPaperID int64) ([]examentity.ExamSection, error) {
	var sections []examentity.ExamSection
	if err := examdao.ExamSection.Ctx(ctx).
		Where("exam_paper_id", examPaperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&sections); err != nil {
		return nil, err
	}
	return sections, nil
}

func loadPaperSectionByID(ctx context.Context, examPaperID, sectionID int64) (examentity.ExamSection, error) {
	var section examentity.ExamSection
	if err := examdao.ExamSection.Ctx(ctx).
		Where("id", sectionID).
		Where("exam_paper_id", examPaperID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&section); err != nil {
		return section, err
	}
	if section.Id == 0 {
		return section, gerror.NewCode(consts.CodeExamSectionNotFound)
	}
	return section, nil
}

func loadPaperQuestionGraph(ctx context.Context, sections []examentity.ExamSection, singleSectionID int64) (
	map[int64][]examentity.ExamQuestionBlock,
	map[int64][]examentity.ExamQuestion,
	map[int64][]examentity.ExamOption,
	error,
) {
	blocksBySection := make(map[int64][]examentity.ExamQuestionBlock, len(sections))
	questionsByBlock := map[int64][]examentity.ExamQuestion{}
	optionsByQuestion := map[int64][]examentity.ExamOption{}
	if len(sections) == 0 {
		return blocksBySection, questionsByBlock, optionsByQuestion, nil
	}

	sectionIDs := make([]interface{}, len(sections))
	for i := range sections {
		sectionIDs[i] = sections[i].Id
	}

	var blocks []examentity.ExamQuestionBlock
	blockModel := examdao.ExamQuestionBlock.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted)
	if singleSectionID > 0 {
		blockModel = blockModel.Where("section_id", singleSectionID).OrderAsc("block_order").OrderAsc("id")
	} else {
		blockModel = blockModel.WhereIn("section_id", sectionIDs)
	}
	if err := blockModel.Scan(&blocks); err != nil {
		return nil, nil, nil, err
	}
	for _, block := range blocks {
		blocksBySection[block.SectionId] = append(blocksBySection[block.SectionId], block)
	}
	for sectionID, sectionBlocks := range blocksBySection {
		sort.Slice(sectionBlocks, func(i, j int) bool {
			if sectionBlocks[i].BlockOrder != sectionBlocks[j].BlockOrder {
				return sectionBlocks[i].BlockOrder < sectionBlocks[j].BlockOrder
			}
			return sectionBlocks[i].Id < sectionBlocks[j].Id
		})
		blocksBySection[sectionID] = sectionBlocks
	}

	blockIDs := make([]interface{}, len(blocks))
	for i := range blocks {
		blockIDs[i] = blocks[i].Id
	}
	if len(blockIDs) == 0 {
		return blocksBySection, questionsByBlock, optionsByQuestion, nil
	}

	var questions []examentity.ExamQuestion
	if err := examdao.ExamQuestion.Ctx(ctx).
		WhereIn("block_id", blockIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&questions); err != nil {
		return nil, nil, nil, err
	}
	for _, question := range questions {
		questionsByBlock[question.BlockId] = append(questionsByBlock[question.BlockId], question)
	}
	for blockID, blockQuestions := range questionsByBlock {
		sort.Slice(blockQuestions, func(i, j int) bool {
			if blockQuestions[i].SortInBlock != blockQuestions[j].SortInBlock {
				return blockQuestions[i].SortInBlock < blockQuestions[j].SortInBlock
			}
			return blockQuestions[i].Id < blockQuestions[j].Id
		})
		questionsByBlock[blockID] = blockQuestions
	}

	questionIDs := make([]interface{}, len(questions))
	for i := range questions {
		questionIDs[i] = questions[i].Id
	}
	if len(questionIDs) == 0 {
		return blocksBySection, questionsByBlock, optionsByQuestion, nil
	}

	var options []examentity.ExamOption
	if err := examdao.ExamOption.Ctx(ctx).
		WhereIn("question_id", questionIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&options); err != nil {
		return nil, nil, nil, err
	}
	for _, option := range options {
		optionsByQuestion[option.QuestionId] = append(optionsByQuestion[option.QuestionId], option)
	}
	for questionID, questionOptions := range optionsByQuestion {
		sort.Slice(questionOptions, func(i, j int) bool {
			if questionOptions[i].SortOrder != questionOptions[j].SortOrder {
				return questionOptions[i].SortOrder < questionOptions[j].SortOrder
			}
			return questionOptions[i].Id < questionOptions[j].Id
		})
		optionsByQuestion[questionID] = questionOptions
	}

	return blocksBySection, questionsByBlock, optionsByQuestion, nil
}
