package paper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/sync/singleflight"

	"exam/internal/consts"
	"exam/internal/dao"
	exambo "exam/internal/model/bo/exam"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"
	"exam/internal/utility/exampaper"
)

const mockMetaCacheTTL = 10 * time.Minute

type cachedMockPaperEntry struct {
	data     mockentity.MockExaminationPaper
	cachedAt time.Time
}

type cachedSegmentsEntry struct {
	data     []mockentity.MockExaminationSegment
	cachedAt time.Time
}

type cachedPartsEntry struct {
	data     []mockentity.MockExaminationPart
	cachedAt time.Time
}

var (
	mockPaperMetaCache sync.Map
	mockSegmentsCache  sync.Map
	mockPartsCache     sync.Map

	mockPaperMetaSF singleflight.Group
	mockSegmentsSF  singleflight.Group
	mockPartsSF     singleflight.Group
)

func (s *sPaper) PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (map[string]interface{}, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	rkey := paperSectionTopicRedisKey(paper.Id, sectionId)
	if cached := redisGetSectionTopicJSON(ctx, rkey); cached != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &m); err == nil {
			return m, nil
		}
	}
	v, err, _ := paperSectionTopicSF.Do(rkey, func() (interface{}, error) {
		if cached := redisGetSectionTopicJSON(ctx, rkey); cached != "" {
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(cached), &m); err == nil {
				return m, nil
			}
		}
		m, err := paperSectionTopicFromDB(ctx, paper.Id, sectionId)
		if err != nil {
			return nil, err
		}
		if b, mErr := json.Marshal(m); mErr == nil {
			redisSetSectionTopicJSON(ctx, rkey, string(b))
		}
		return m, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(map[string]interface{}), nil
}

func paperSectionTopicFromDB(ctx context.Context, examPaperId int64, sectionId int64) (map[string]interface{}, error) {
	var sec examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("id", sectionId).
		Where("exam_paper_id", examPaperId).
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
	var blocks []examentity.ExamQuestionBlock
	if err := dao.ExamQuestionBlock.Ctx(ctx).
		Where("section_id", sec.Id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("block_order,id").
		Scan(&blocks); err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		return m, nil
	}
	blockIDs := make([]interface{}, 0, len(blocks))
	for _, b := range blocks {
		blockIDs = append(blockIDs, b.Id)
	}
	var questions []examentity.ExamQuestion
	if err := dao.ExamQuestion.Ctx(ctx).
		Fields("id", "block_id", "sort_in_block", "is_example").
		Where("exam_paper_id", examPaperId).
		WhereIn("block_id", blockIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("block_id,sort_in_block,id").
		Scan(&questions); err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return m, nil
	}
	questionsByBlock := make(map[int64][]examentity.ExamQuestion, len(blocks))
	questionIDs := make([]interface{}, 0, len(questions))
	for _, q := range questions {
		questionsByBlock[q.BlockId] = append(questionsByBlock[q.BlockId], q)
		questionIDs = append(questionIDs, q.Id)
	}
	for bid := range questionsByBlock {
		sort.Slice(questionsByBlock[bid], func(i, j int) bool {
			if questionsByBlock[bid][i].SortInBlock != questionsByBlock[bid][j].SortInBlock {
				return questionsByBlock[bid][i].SortInBlock < questionsByBlock[bid][j].SortInBlock
			}
			return questionsByBlock[bid][i].Id < questionsByBlock[bid][j].Id
		})
	}
	var options []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		Fields("id", "question_id", "sort_order").
		WhereIn("question_id", questionIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("question_id,sort_order,id").
		Scan(&options); err != nil {
		return nil, err
	}
	optionsByQuestion := make(map[int64][]examentity.ExamOption, len(questions))
	for _, opt := range options {
		optionsByQuestion[opt.QuestionId] = append(optionsByQuestion[opt.QuestionId], opt)
	}
	for qid := range optionsByQuestion {
		sort.Slice(optionsByQuestion[qid], func(i, j int) bool {
			if optionsByQuestion[qid][i].SortOrder != optionsByQuestion[qid][j].SortOrder {
				return optionsByQuestion[qid][i].SortOrder < optionsByQuestion[qid][j].SortOrder
			}
			return optionsByQuestion[qid][i].Id < optionsByQuestion[qid][j].Id
		})
	}
	enrichTopicWithExamIDs(m, blocks, questionsByBlock, optionsByQuestion)
	stripSensitiveExamFields(m)
	return m, nil
}

func enrichTopicWithExamIDs(
	topic map[string]interface{},
	blocks []examentity.ExamQuestionBlock,
	questionsByBlock map[int64][]examentity.ExamQuestion,
	optionsByQuestion map[int64][]examentity.ExamOption,
) {
	rawItems, ok := topic["items"]
	if !ok {
		return
	}
	items, ok := rawItems.([]interface{})
	if !ok {
		return
	}
	for i, rawItem := range items {
		if i >= len(blocks) {
			return
		}
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}
		blockQs := questionsByBlock[blocks[i].Id]
		if len(blockQs) == 0 {
			continue
		}
		if rawQuestions, ok := item["questions"]; ok {
			questions, ok := rawQuestions.([]interface{})
			if !ok {
				continue
			}
			for qi, rawQuestion := range questions {
				if qi >= len(blockQs) {
					break
				}
				question, ok := rawQuestion.(map[string]interface{})
				if !ok {
					continue
				}
				enrichQuestionWithExamIDs(question, blockQs[qi], optionsByQuestion[blockQs[qi].Id])
			}
			continue
		}
		enrichQuestionWithExamIDs(item, blockQs[0], optionsByQuestion[blockQs[0].Id])
	}
}

func enrichQuestionWithExamIDs(question map[string]interface{}, q examentity.ExamQuestion, options []examentity.ExamOption) {
	question["eqid"] = q.Id
	rawAnswers, ok := question["answers"]
	if !ok {
		return
	}
	answers, ok := rawAnswers.([]interface{})
	if !ok {
		return
	}
	optionIDBySortOrder := make(map[int]int64, len(options))
	optionIDByFlag := make(map[string]int64, len(options))
	optionIDByIDString := make(map[string]int64, len(options))
	for _, opt := range options {
		optionIDBySortOrder[opt.SortOrder] = opt.Id
		if opt.Flag != "" {
			optionIDByFlag[opt.Flag] = opt.Id
		}
		optionIDByIDString[strconv.FormatInt(opt.Id, 10)] = opt.Id
	}
	for ai, rawAnswer := range answers {
		answer, ok := rawAnswer.(map[string]interface{})
		if !ok {
			continue
		}
		var eaid int64
		matched := false
		if fv, ok := answer["flag"]; ok {
			if flag, ok := fv.(string); ok && flag != "" {
				if v, ok := optionIDByFlag[flag]; ok {
					eaid = v
					matched = true
				}
			}
		}
		rawID := ""
		if !matched {
			if iv, ok := answer["id"]; ok {
				if idStr, ok := iv.(string); ok {
					rawID = idStr
				}
				if rawID != "" {
					if v, ok := optionIDByIDString[rawID]; ok {
						eaid = v
						matched = true
					}
				}
			}
		}
		if !matched && rawID == "" {
			sortOrder := ai
			if iv, ok := answer["index"]; ok {
				if n, ok := iv.(float64); ok {
					sortOrder = int(n)
				}
			}
			if v, ok := optionIDBySortOrder[sortOrder]; ok {
				eaid = v
				matched = true
			} else if v, ok := optionIDBySortOrder[sortOrder+1]; ok {
				// 兼容 index 起点与 DB sort_order 的 0/1 差异。
				eaid = v
				matched = true
			}
		}
		if matched {
			answer["eaid"] = eaid
		}
	}
}

func stripSensitiveExamFields(node interface{}) {
	switch v := node.(type) {
	case map[string]interface{}:
		// 移除正确答案与分数相关字段，避免考试前泄露答案或计分信息。
		delete(v, "correct_answer")
		delete(v, "correct")
		delete(v, "answer")
		delete(v, "score")
		delete(v, "total_score")
		delete(v, "score_total")
		delete(v, "question_score")
		delete(v, "part_score")
		delete(v, "part_rate")
		delete(v, "objective_score")
		delete(v, "subjective_score")
		delete(v, "correct_count")
		delete(v, "correct_rate")
		for _, child := range v {
			stripSensitiveExamFields(child)
		}
	case []interface{}:
		for _, child := range v {
			stripSensitiveExamFields(child)
		}
	}
}

func (s *sPaper) PaperDetailForExamInit(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	mockPaper, err := loadMockPaperByID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	t, err := PaperDetailForExamInit(ctx, paper.Id)
	if err != nil {
		return nil, err
	}
	out := paperDetailForExamInitTreeToBO(t)
	out.Paper.ListenReviewDuration = mockPaper.ListenReviewDuration
	return &out, nil
}

func (s *sPaper) PaperBootstrapForExam(ctx context.Context, mockPaperID int64) (*exambo.PaperDetailForExamInitTree, []exambo.PaperPrepareSegment, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, nil, err
	}
	mockPaper, err := loadMockPaperByID(ctx, mockPaperID)
	if err != nil {
		return nil, nil, err
	}
	t, err := PaperDetailForExamInit(ctx, paper.Id)
	if err != nil {
		return nil, nil, err
	}
	detail := paperDetailForExamInitTreeToBO(t)
	detail.Paper.ListenReviewDuration = mockPaper.ListenReviewDuration

	segments, err := paperPrepareSegmentsByLevelID(ctx, mockPaper.LevelId)
	if err != nil {
		return nil, nil, err
	}
	return &detail, segments, nil
}

func (s *sPaper) PaperPrepareSegments(ctx context.Context, mockPaperID int64) ([]exambo.PaperPrepareSegment, error) {
	rkey := paperPrepareRedisKey(mockPaperID)
	if cached := redisGetPrepareJSON(ctx, rkey); cached != "" {
		var out []exambo.PaperPrepareSegment
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			return out, nil
		}
	}
	v, err, _ := paperPrepareSF.Do(rkey, func() (interface{}, error) {
		if cached := redisGetPrepareJSON(ctx, rkey); cached != "" {
			var out []exambo.PaperPrepareSegment
			if err := json.Unmarshal([]byte(cached), &out); err == nil {
				return out, nil
			}
		}
		out, err := paperPrepareSegmentsFromDB(ctx, mockPaperID)
		if err != nil {
			return nil, err
		}
		if b, mErr := json.Marshal(out); mErr == nil {
			redisSetPrepareJSON(ctx, rkey, string(b))
		}
		return out, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]exambo.PaperPrepareSegment), nil
}

func paperPrepareSegmentsFromDB(ctx context.Context, mockPaperID int64) ([]exambo.PaperPrepareSegment, error) {
	mockPaper, err := loadMockPaperByID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	return paperPrepareSegmentsByLevelID(ctx, mockPaper.LevelId)
}

func paperPrepareSegmentsByLevelID(ctx context.Context, levelID int64) ([]exambo.PaperPrepareSegment, error) {
	segments, err := loadSegmentsByLevelID(ctx, levelID)
	if err != nil {
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

	parts, err := loadPartsByLevelID(ctx, levelID, segmentIDs)
	if err != nil {
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

func loadMockPaperByID(ctx context.Context, mockPaperID int64) (mockentity.MockExaminationPaper, error) {
	if entry, ok := mockPaperMetaCache.Load(mockPaperID); ok {
		e := entry.(*cachedMockPaperEntry)
		if time.Since(e.cachedAt) < mockMetaCacheTTL {
			return e.data, nil
		}
		mockPaperMetaCache.Delete(mockPaperID)
	}
	sfKey := fmt.Sprintf("mock-paper:%d", mockPaperID)
	v, err, _ := mockPaperMetaSF.Do(sfKey, func() (interface{}, error) {
		if entry, ok := mockPaperMetaCache.Load(mockPaperID); ok {
			e := entry.(*cachedMockPaperEntry)
			if time.Since(e.cachedAt) < mockMetaCacheTTL {
				return &e.data, nil
			}
			mockPaperMetaCache.Delete(mockPaperID)
		}
		var mockPaper mockentity.MockExaminationPaper
		if err := dao.MockExaminationPaper.Ctx(ctx).
			Where("id", mockPaperID).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&mockPaper); err != nil {
			return mockentity.MockExaminationPaper{}, err
		}
		if mockPaper.Id == 0 {
			return mockentity.MockExaminationPaper{}, gerror.NewCode(consts.CodeMockExamPaperNotFound)
		}
		mockPaperMetaCache.Store(mockPaperID, &cachedMockPaperEntry{data: mockPaper, cachedAt: time.Now()})
		return &mockPaper, nil
	})
	if err != nil {
		return mockentity.MockExaminationPaper{}, err
	}
	return *v.(*mockentity.MockExaminationPaper), nil
}

func loadSegmentsByLevelID(ctx context.Context, levelID int64) ([]mockentity.MockExaminationSegment, error) {
	if entry, ok := mockSegmentsCache.Load(levelID); ok {
		e := entry.(*cachedSegmentsEntry)
		if time.Since(e.cachedAt) < mockMetaCacheTTL {
			return e.data, nil
		}
		mockSegmentsCache.Delete(levelID)
	}
	sfKey := fmt.Sprintf("mock-segments:%d", levelID)
	v, err, _ := mockSegmentsSF.Do(sfKey, func() (interface{}, error) {
		if entry, ok := mockSegmentsCache.Load(levelID); ok {
			e := entry.(*cachedSegmentsEntry)
			if time.Since(e.cachedAt) < mockMetaCacheTTL {
				return e.data, nil
			}
			mockSegmentsCache.Delete(levelID)
		}
		segments := make([]mockentity.MockExaminationSegment, 0)
		if err := dao.MockExaminationSegment.Ctx(ctx).
			Where("level_id", levelID).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("seq,id").
			Scan(&segments); err != nil {
			return nil, err
		}
		mockSegmentsCache.Store(levelID, &cachedSegmentsEntry{data: segments, cachedAt: time.Now()})
		return segments, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]mockentity.MockExaminationSegment), nil
}

func loadPartsByLevelID(ctx context.Context, levelID int64, segmentIDs []int64) ([]mockentity.MockExaminationPart, error) {
	if entry, ok := mockPartsCache.Load(levelID); ok {
		e := entry.(*cachedPartsEntry)
		if time.Since(e.cachedAt) < mockMetaCacheTTL {
			return e.data, nil
		}
		mockPartsCache.Delete(levelID)
	}
	sfKey := fmt.Sprintf("mock-parts:%d", levelID)
	v, err, _ := mockPartsSF.Do(sfKey, func() (interface{}, error) {
		if entry, ok := mockPartsCache.Load(levelID); ok {
			e := entry.(*cachedPartsEntry)
			if time.Since(e.cachedAt) < mockMetaCacheTTL {
				return e.data, nil
			}
			mockPartsCache.Delete(levelID)
		}
		parts := make([]mockentity.MockExaminationPart, 0)
		if err := dao.MockExaminationPart.Ctx(ctx).
			WhereIn("segment_id", segmentIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("segment_id,code,id").
			Scan(&parts); err != nil {
			return nil, err
		}
		mockPartsCache.Store(levelID, &cachedPartsEntry{data: parts, cachedAt: time.Now()})
		return parts, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]mockentity.MockExaminationPart), nil
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
