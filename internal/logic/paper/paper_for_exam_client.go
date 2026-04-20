package paper

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
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

// sectionTopicMemTTL 进程内 SectionTopic 对象缓存的 TTL。
// 该缓存叠加在 Redis 之上，目的是消除热路径上"Redis GET + 大对象 json.Unmarshal"这两大开销。
// 写入方（DB miss 路径）和显式 invalidate 调用会同步更新/清理；即便有极端遗漏，30s 后也会自然过期。
const sectionTopicMemTTL = 30 * time.Second

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

type sectionTopicMemKey struct {
	paperID   int64
	sectionID int64
}

type cachedSectionTopicEntry struct {
	data     *exambo.SectionTopic
	cachedAt time.Time
}

var (
	mockPaperMetaCache sync.Map
	mockSegmentsCache  sync.Map
	mockPartsCache     sync.Map

	mockPaperMetaSF singleflight.Group
	mockSegmentsSF  singleflight.Group
	mockPartsSF     singleflight.Group

	sectionTopicMemCache sync.Map // sectionTopicMemKey -> *cachedSectionTopicEntry
)

func loadSectionTopicFromMem(paperID, sectionID int64) *exambo.SectionTopic {
	if v, ok := sectionTopicMemCache.Load(sectionTopicMemKey{paperID: paperID, sectionID: sectionID}); ok {
		e := v.(*cachedSectionTopicEntry)
		if time.Since(e.cachedAt) < sectionTopicMemTTL {
			return e.data
		}
		sectionTopicMemCache.Delete(sectionTopicMemKey{paperID: paperID, sectionID: sectionID})
	}
	return nil
}

func storeSectionTopicToMem(paperID, sectionID int64, topic *exambo.SectionTopic) {
	if topic == nil {
		return
	}
	sectionTopicMemCache.Store(
		sectionTopicMemKey{paperID: paperID, sectionID: sectionID},
		&cachedSectionTopicEntry{data: topic, cachedAt: time.Now()},
	)
}

func invalidateSectionTopicMemCache(paperID, sectionID int64) {
	sectionTopicMemCache.Delete(sectionTopicMemKey{paperID: paperID, sectionID: sectionID})
}

func invalidateSectionTopicMemCacheByPaper(paperID int64) {
	sectionTopicMemCache.Range(func(k, _ any) bool {
		if key, ok := k.(sectionTopicMemKey); ok && key.paperID == paperID {
			sectionTopicMemCache.Delete(key)
		}
		return true
	})
}

func (s *sPaper) PaperSectionTopicForExam(ctx context.Context, mockPaperID int64, sectionId int64) (*exambo.SectionTopic, error) {
	paper, err := exampaper.ByMockID(ctx, mockPaperID)
	if err != nil {
		return nil, err
	}
	// L1: 进程内对象缓存。命中即直接返回，绕过 Redis 往返和大对象 JSON 反序列化。
	if t := loadSectionTopicFromMem(paper.Id, sectionId); t != nil {
		return t, nil
	}
	rkey := paperSectionTopicRedisKey(paper.Id, sectionId)
	// L2: Redis 缓存。写入前已经做过 YCT 剥离和 EOID 补全，这里无需任何后处理。
	if cached := redisGetSectionTopicJSON(ctx, rkey); cached != "" {
		var t exambo.SectionTopic
		if err := json.Unmarshal([]byte(cached), &t); err == nil {
			storeSectionTopicToMem(paper.Id, sectionId, &t)
			return &t, nil
		}
	}
	v, err, _ := paperSectionTopicSF.Do(rkey, func() (interface{}, error) {
		// 单飞组里再检查一次 mem + redis，避免并发请求重复回源。
		if t := loadSectionTopicFromMem(paper.Id, sectionId); t != nil {
			return t, nil
		}
		if cached := redisGetSectionTopicJSON(ctx, rkey); cached != "" {
			var t exambo.SectionTopic
			if err := json.Unmarshal([]byte(cached), &t); err == nil {
				storeSectionTopicToMem(paper.Id, sectionId, &t)
				return &t, nil
			}
		}
		// 真正的 DB miss 路径。只有这里才需要查一次 level 判断 YCT，读热路径完全不碰 DB。
		isYCTPaper, yErr := isYCTMockPaper(ctx, mockPaperID)
		if yErr != nil {
			return nil, yErr
		}
		t, err := paperSectionTopicFromDB(ctx, paper.Id, sectionId)
		if err != nil {
			return nil, err
		}
		if isYCTPaper {
			stripYCTItemRenderFields(t)
		}
		if b, mErr := json.Marshal(t); mErr == nil {
			redisSetSectionTopicJSON(ctx, rkey, string(b))
		}
		storeSectionTopicToMem(paper.Id, sectionId, t)
		return t, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*exambo.SectionTopic), nil
}

func isYCTMockPaper(ctx context.Context, mockPaperID int64) (bool, error) {
	mockPaper, err := loadMockPaperByID(ctx, mockPaperID)
	if err != nil {
		return false, err
	}
	var level mockentity.MockLevels
	if err := dao.MockLevels.Ctx(ctx).
		Fields("id", "type_name").
		Where("id", mockPaper.LevelId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&level); err != nil {
		return false, err
	}
	return level.Id > 0 && level.TypeName == "YCT", nil
}

func stripYCTItemRenderFields(topic *exambo.SectionTopic) {
	for i := range topic.Items {
		if topic.Items[i].Extra == nil {
			continue
		}
		delete(topic.Items[i].Extra, "_converter")
		delete(topic.Items[i].Extra, "_element")
	}
}

func paperSectionTopicFromDB(ctx context.Context, examPaperId int64, sectionId int64) (*exambo.SectionTopic, error) {
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
	var t exambo.SectionTopic
	if err := json.Unmarshal([]byte(sec.TopicJson), &t); err != nil {
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
		stripSensitiveExamFields(&t)
		return &t, nil
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
		stripSensitiveExamFields(&t)
		return &t, nil
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
		Fields("id", "question_id", "sort_order", "flag").
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
	enrichTopicWithExamIDs(&t, blocks, questionsByBlock, optionsByQuestion)
	stripSensitiveExamFields(&t)
	return &t, nil
}

func enrichTopicWithExamIDs(
	topic *exambo.SectionTopic,
	blocks []examentity.ExamQuestionBlock,
	questionsByBlock map[int64][]examentity.ExamQuestion,
	optionsByQuestion map[int64][]examentity.ExamOption,
) {
	for i := range topic.Items {
		if i >= len(blocks) {
			return
		}
		item := &topic.Items[i]
		blockQs := questionsByBlock[blocks[i].Id]
		if len(blockQs) == 0 {
			continue
		}
		if len(item.Questions) > 0 {
			for qi := range item.Questions {
				if qi >= len(blockQs) {
					break
				}
				enrichQuestionWithExamIDs(&item.Questions[qi], blockQs[qi], optionsByQuestion[blockQs[qi].Id])
			}
			continue
		}
		item.Eqid = &blockQs[0].Id
		enrichAnswersWithExamIDs(item.Answers, optionsByQuestion[blockQs[0].Id])
	}
}

func enrichQuestionWithExamIDs(question *exambo.TopicQuestion, q examentity.ExamQuestion, options []examentity.ExamOption) {
	question.Eqid = &q.Id
	enrichAnswersWithExamIDs(question.Answers, options)
}

func enrichAnswersWithExamIDs(answers []exambo.TopicAnswer, options []examentity.ExamOption) {
	optionIDBySortOrder := make(map[int]int64, len(options))
	optionIDByFlag := make(map[string]int64, len(options))
	optionIDByIDString := make(map[string]int64, len(options))
	for _, opt := range options {
		optionIDBySortOrder[opt.SortOrder] = opt.Id
		if opt.Flag != "" {
			optionIDByFlag[strings.ToUpper(strings.TrimSpace(opt.Flag))] = opt.Id
		}
		optionIDByIDString[strconv.FormatInt(opt.Id, 10)] = opt.Id
	}
	for ai := range answers {
		answer := &answers[ai]
		var eoid int64
		matched := false
		if answer.Flag != "" {
			if v, ok := optionIDByFlag[strings.ToUpper(strings.TrimSpace(answer.Flag))]; ok {
				eoid = v
				matched = true
			}
		}
		if !matched {
			if rawID, ok := exambo.RawString(answer.Id); ok && rawID != "" {
				if v, ok := optionIDByIDString[rawID]; ok {
					eoid = v
					matched = true
				}
			}
		}
		if !matched {
			sortOrder := ai
			if n, ok := exambo.RawInt(answer.Index); ok {
				sortOrder = n
			}
			if v, ok := optionIDBySortOrder[sortOrder]; ok {
				eoid = v
				matched = true
			} else if v, ok := optionIDBySortOrder[sortOrder+1]; ok {
				// 兼容 index 起点与 DB sort_order 的 0/1 差异。
				eoid = v
				matched = true
			}
		}
		if matched {
			answer.EOID = &eoid
		}
	}
}

func topicHasStaleEoid(topic *exambo.SectionTopic) bool {
	if topic == nil {
		return false
	}
	for i := range topic.Items {
		item := &topic.Items[i]
		if answersHaveMixedEOID(item.Answers) {
			return true
		}
		for qi := range item.Questions {
			if answersHaveMixedEOID(item.Questions[qi].Answers) {
				return true
			}
		}
	}
	return false
}

func answersHaveMixedEOID(answers []exambo.TopicAnswer) bool {
	if len(answers) <= 1 {
		return false
	}
	hasEOID := false
	hasMissingEOID := false
	for i := range answers {
		if answers[i].EOID != nil {
			hasEOID = true
		} else {
			hasMissingEOID = true
		}
		if hasEOID && hasMissingEOID {
			return true
		}
	}
	return false
}

func stripSensitiveExamFields(topic *exambo.SectionTopic) {
	rootIsExample := exambo.RawTruthy(topic.IsExample)
	stripSensitiveFieldsOnExtra(topic.Extra, rootIsExample)
	for i := range topic.Items {
		item := &topic.Items[i]
		itemIsExample := rootIsExample || exambo.RawTruthy(item.IsExample)
		stripSensitiveFieldsOnExtra(item.Extra, itemIsExample)
		for qi := range item.Questions {
			q := &item.Questions[qi]
			qIsExample := itemIsExample || exambo.RawTruthy(q.IsExample)
			stripSensitiveFieldsOnExtra(q.Extra, qIsExample)
			for ai := range q.Answers {
				stripSensitiveFieldsOnExtra(q.Answers[ai].Extra, qIsExample)
			}
		}
		for ai := range item.Answers {
			stripSensitiveFieldsOnExtra(item.Answers[ai].Extra, itemIsExample)
		}
	}
}

func stripSensitiveFieldsOnExtra(v map[string]json.RawMessage, isExample bool) {
	if v == nil {
		return
	}
	if !isExample {
		delete(v, "correct_answer")
		delete(v, "correct")
		delete(v, "answer")
	}
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
