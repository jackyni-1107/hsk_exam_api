package exam

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/singleflight"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/entity"
	"exam/internal/util"
)

const (
	paperForExamCacheTTLSeconds = 3600
	paperForExamMaxStringBytes  = 256 * 1024
)

func paperForExamInitRedisKey(examPaperId int64) string {
	return fmt.Sprintf("exam:paper:init:%d", examPaperId)
}

// paperForExamLegacyRedisKey 历史整卷考前缓存（大 value），失效时一并删除以免长期占用 Redis。
func paperForExamLegacyRedisKey(examPaperId int64) string {
	return fmt.Sprintf("exam:paper:%d", examPaperId)
}

func paperForExamSectionRedisKey(examPaperId, sectionId int64) string {
	return fmt.Sprintf("exam:paper:section:%d:%d", examPaperId, sectionId)
}

var paperForExamInitSF singleflight.Group
var paperForExamSectionSF singleflight.Group

// InvalidatePaperForExamCache 试卷树变更后删除考前相关缓存（初始化 + 各 section 详情 + 历史整卷 key）。
func InvalidatePaperForExamCache(ctx context.Context, examPaperId int64) {
	if examPaperId <= 0 {
		return
	}
	initKey := paperForExamInitRedisKey(examPaperId)
	legacyKey := paperForExamLegacyRedisKey(examPaperId)
	if _, err := g.Redis().Del(ctx, initKey, legacyKey); err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis del init/legacy %s %s: %v", initKey, legacyKey, err)
	}
	invalidatePaperSectionExamCachesByPaper(ctx, examPaperId)
}

// InvalidatePaperSectionForExamCache 删除单个 section 的考前详情缓存（精确 key）。
func InvalidatePaperSectionForExamCache(ctx context.Context, examPaperId, sectionId int64) {
	if examPaperId <= 0 || sectionId <= 0 {
		return
	}
	key := paperForExamSectionRedisKey(examPaperId, sectionId)
	if _, err := g.Redis().Del(ctx, key); err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis del section %s: %v", key, err)
	}
}

func invalidatePaperSectionExamCachesByPaper(ctx context.Context, examPaperId int64) {
	pattern := fmt.Sprintf("exam:paper:section:%d:*", examPaperId)
	v, err := g.Redis().Do(ctx, "KEYS", pattern)
	if err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis keys %s: %v", pattern, err)
		return
	}
	keys := v.Strings()
	if len(keys) == 0 {
		return
	}
	if _, err := g.Redis().Del(ctx, keys...); err != nil {
		g.Log().Warningf(ctx, "paper for-exam redis del section keys pattern %s: %v", pattern, err)
	}
}

func redisGetPaperForExamInit(ctx context.Context, rkey string) *PaperDetailForExamInitTree {
	val, err := g.Redis().Get(ctx, rkey)
	if err != nil {
		g.Log().Warningf(ctx, "paper for-exam init redis get %s: %v", rkey, err)
		return nil
	}
	if val.IsEmpty() {
		return nil
	}
	var out PaperDetailForExamInitTree
	if err := json.Unmarshal([]byte(val.String()), &out); err != nil {
		g.Log().Warningf(ctx, "paper for-exam init redis unmarshal %s: %v", rkey, err)
		return nil
	}
	return &out
}

// PaperDetailForExamInit 客户端考前初始化：仅试卷结构（paper + section 概要 + block 概要 + 题量），不含题目与选项。
// 流程：Redis → singleflight → DB；TTL 1h。
func PaperDetailForExamInit(ctx context.Context, examPaperId int64) (*PaperDetailForExamInitTree, error) {
	rkey := paperForExamInitRedisKey(examPaperId)
	if hit := redisGetPaperForExamInit(ctx, rkey); hit != nil {
		return hit, nil
	}
	v, err, _ := paperForExamInitSF.Do(rkey, func() (interface{}, error) {
		if hit := redisGetPaperForExamInit(ctx, rkey); hit != nil {
			return hit, nil
		}
		tree, err := paperDetailForExamInitFromDB(ctx, examPaperId)
		if err != nil {
			return nil, err
		}
		b, mErr := json.Marshal(tree)
		if mErr != nil {
			g.Log().Warningf(ctx, "paper for-exam init json marshal for redis: %v", mErr)
			return tree, nil
		}
		if err := g.Redis().SetEX(ctx, rkey, string(b), paperForExamCacheTTLSeconds); err != nil {
			g.Log().Warningf(ctx, "paper for-exam init redis setex %s: %v", rkey, err)
		}
		return tree, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*PaperDetailForExamInitTree), nil
}

func redisGetPaperSectionForExam(ctx context.Context, rkey string) *SectionDetailForExamView {
	val, err := g.Redis().Get(ctx, rkey)
	if err != nil {
		g.Log().Warningf(ctx, "paper for-exam section redis get %s: %v", rkey, err)
		return nil
	}
	if val.IsEmpty() {
		return nil
	}
	var out SectionDetailForExamView
	if err := json.Unmarshal([]byte(val.String()), &out); err != nil {
		g.Log().Warningf(ctx, "paper for-exam section redis unmarshal %s: %v", rkey, err)
		return nil
	}
	return &out
}

// PaperSectionDetailForExam 按 section 拉取完整题目树（blocks + questions + options），不含选项正误。
func PaperSectionDetailForExam(ctx context.Context, examPaperId, sectionId int64) (*SectionDetailForExamView, error) {
	rkey := paperForExamSectionRedisKey(examPaperId, sectionId)
	if hit := redisGetPaperSectionForExam(ctx, rkey); hit != nil {
		return hit, nil
	}
	v, err, _ := paperForExamSectionSF.Do(rkey, func() (interface{}, error) {
		if hit := redisGetPaperSectionForExam(ctx, rkey); hit != nil {
			return hit, nil
		}
		secView, err := paperSectionDetailForExamFromDB(ctx, examPaperId, sectionId)
		if err != nil {
			return nil, err
		}
		b, mErr := json.Marshal(secView)
		if mErr != nil {
			g.Log().Warningf(ctx, "paper for-exam section json marshal for redis: %v", mErr)
			return secView, nil
		}
		if err := g.Redis().SetEX(ctx, rkey, string(b), paperForExamCacheTTLSeconds); err != nil {
			g.Log().Warningf(ctx, "paper for-exam section redis setex %s: %v", rkey, err)
		}
		return secView, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*SectionDetailForExamView), nil
}

func trimForExamLarge(s string) string {
	if len(s) > paperForExamMaxStringBytes {
		return ""
	}
	return s
}

func paperDetailForExamInitFromDB(ctx context.Context, examPaperId int64) (*PaperDetailForExamInitTree, error) {
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
		OrderAsc("id").
		Scan(&sections); err != nil {
		return nil, err
	}

	sectionIDs := make([]interface{}, 0, len(sections))
	for _, sec := range sections {
		sectionIDs = append(sectionIDs, sec.Id)
	}

	var allBlocks []entity.ExamQuestionBlock
	if len(sectionIDs) > 0 {
		if err := dao.ExamQuestionBlock.Ctx(ctx).
			WhereIn("section_id", sectionIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("section_id").
			OrderAsc("block_order").
			OrderAsc("id").
			Scan(&allBlocks); err != nil {
			return nil, err
		}
	}

	blocksBySection := make(map[int64][]entity.ExamQuestionBlock, len(sections))
	blockIDs := make([]interface{}, 0, len(allBlocks))
	for _, b := range allBlocks {
		blocksBySection[b.SectionId] = append(blocksBySection[b.SectionId], b)
		blockIDs = append(blockIDs, b.Id)
	}
	for sid, bl := range blocksBySection {
		sort.Slice(bl, func(i, j int) bool {
			if bl[i].BlockOrder != bl[j].BlockOrder {
				return bl[i].BlockOrder < bl[j].BlockOrder
			}
			return bl[i].Id < bl[j].Id
		})
		blocksBySection[sid] = bl
	}

	questionCountByBlock := make(map[int64]int, len(allBlocks))
	if len(blockIDs) > 0 {
		type blockCountRow struct {
			BlockId int64 `json:"block_id"`
			Cnt     int64 `json:"cnt"`
		}
		var countRows []blockCountRow
		if err := dao.ExamQuestion.Ctx(ctx).
			Fields("block_id", "COUNT(1) AS cnt").
			Where("exam_paper_id", examPaperId).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			WhereIn("block_id", blockIDs).
			Group("block_id").
			Scan(&countRows); err != nil {
			return nil, err
		}
		for _, row := range countRows {
			questionCountByBlock[row.BlockId] = int(row.Cnt)
		}
	}

	out := &PaperDetailForExamInitTree{
		Paper: paperHeadForExam(paper),
	}
	for _, sec := range sections {
		sv := SectionOutlineForExamView{
			Id:             sec.Id,
			SortOrder:      sec.SortOrder,
			TopicTitle:     sec.TopicTitle,
			TopicSubtitle:  sec.TopicSubtitle,
			TopicType:      sec.TopicType,
			PartCode:       sec.PartCode,
			SegmentCode:    sec.SegmentCode,
			TopicItemsFile: sec.TopicItemsFile,
			//TopicJson:      trimForExamLarge(sec.TopicJson),
		}
		for _, blk := range blocksBySection[sec.Id] {
			sv.Blocks = append(sv.Blocks, BlockOutlineForExamView{
				Id:                      blk.Id,
				BlockOrder:              blk.BlockOrder,
				GroupIndex:              blk.GroupIndex,
				QuestionDescriptionJson: trimForExamLarge(blk.QuestionDescriptionJson),
				QuestionCount:           questionCountByBlock[blk.Id],
			})
		}
		out.Sections = append(out.Sections, sv)
	}
	return out, nil
}

func paperSectionDetailForExamFromDB(ctx context.Context, examPaperId, sectionId int64) (*SectionDetailForExamView, error) {
	var sec entity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("id", sectionId).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&sec); err != nil {
		return nil, err
	}
	if sec.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.exam_section_not_found")
	}

	var blocks []entity.ExamQuestionBlock
	if err := dao.ExamQuestionBlock.Ctx(ctx).
		Where("section_id", sectionId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("block_order").
		OrderAsc("id").
		Scan(&blocks); err != nil {
		return nil, err
	}

	blockIDs := make([]interface{}, 0, len(blocks))
	for _, b := range blocks {
		blockIDs = append(blockIDs, b.Id)
	}

	qcols := dao.ExamQuestion.Columns()
	var allQuestions []entity.ExamQuestion
	if len(blockIDs) > 0 {
		if err := dao.ExamQuestion.Ctx(ctx).
			Fields(
				qcols.Id,
				qcols.BlockId,
				qcols.SortInBlock,
				qcols.QuestionNo,
				qcols.Score,
				qcols.IsExample,
				qcols.IsSubjective,
				qcols.ContentType,
				qcols.AudioFile,
				qcols.StemText,
				qcols.ScreenTextJson,
				qcols.QuestionDescriptionJson,
			).
			Where("exam_paper_id", examPaperId).
			WhereIn("block_id", blockIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("block_id").
			OrderAsc("sort_in_block").
			OrderAsc("id").
			Scan(&allQuestions); err != nil {
			return nil, err
		}
	}

	qsByBlock := make(map[int64][]entity.ExamQuestion)
	for _, q := range allQuestions {
		qsByBlock[q.BlockId] = append(qsByBlock[q.BlockId], q)
	}
	for bid, qs := range qsByBlock {
		sort.Slice(qs, func(i, j int) bool {
			if qs[i].SortInBlock != qs[j].SortInBlock {
				return qs[i].SortInBlock < qs[j].SortInBlock
			}
			return qs[i].Id < qs[j].Id
		})
		qsByBlock[bid] = qs
	}

	qIDs := make([]interface{}, 0, len(allQuestions))
	for _, q := range allQuestions {
		qIDs = append(qIDs, q.Id)
	}

	var allOpts []entity.ExamOption
	if len(qIDs) > 0 {
		if err := dao.ExamOption.Ctx(ctx).
			WhereIn("question_id", qIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			OrderAsc("question_id").
			OrderAsc("sort_order").
			OrderAsc("id").
			Scan(&allOpts); err != nil {
			return nil, err
		}
	}
	optsByQ := make(map[int64][]entity.ExamOption)
	for _, o := range allOpts {
		optsByQ[o.QuestionId] = append(optsByQ[o.QuestionId], o)
	}
	for qid, opts := range optsByQ {
		sort.Slice(opts, func(i, j int) bool {
			if opts[i].SortOrder != opts[j].SortOrder {
				return opts[i].SortOrder < opts[j].SortOrder
			}
			return opts[i].Id < opts[j].Id
		})
		optsByQ[qid] = opts
	}

	out := &SectionDetailForExamView{
		Id:             sec.Id,
		SortOrder:      sec.SortOrder,
		TopicTitle:     sec.TopicTitle,
		TopicSubtitle:  sec.TopicSubtitle,
		TopicType:      sec.TopicType,
		PartCode:       sec.PartCode,
		SegmentCode:    sec.SegmentCode,
		TopicItemsFile: sec.TopicItemsFile,
		//TopicJson:      trimForExamLarge(sec.TopicJson),
	}
	for _, blk := range blocks {
		bv := BlockDetailForExamView{
			Id:                      blk.Id,
			BlockOrder:              blk.BlockOrder,
			GroupIndex:              blk.GroupIndex,
			QuestionDescriptionJson: trimForExamLarge(blk.QuestionDescriptionJson),
		}
		for _, q := range qsByBlock[blk.Id] {
			// TODO(产品确认): 考试中若不需要解析说明，可删除 analysis_json 字段并停止 SELECT 该列以进一步瘦身。
			qv := QuestionDetailForExamView{
				Id:                      q.Id,
				SortInBlock:             q.SortInBlock,
				QuestionNo:              q.QuestionNo,
				Score:                   q.Score,
				IsExample:               q.IsExample,
				IsSubjective:            q.IsSubjective,
				ContentType:             q.ContentType,
				AudioFile:               q.AudioFile,
				StemText:                trimForExamLarge(q.StemText),
				ScreenTextJson:          trimForExamLarge(q.ScreenTextJson),
				AnalysisJson:            "",
				QuestionDescriptionJson: trimForExamLarge(q.QuestionDescriptionJson),
				RawJson:                 "",
			}
			for _, o := range optsByQ[q.Id] {
				qv.Options = append(qv.Options, OptionDetailForExamView{
					Id:         o.Id,
					Flag:       o.Flag,
					SortOrder:  o.SortOrder,
					OptionType: o.OptionType,
					Content:    trimForExamLarge(o.Content),
				})
			}
			bv.Questions = append(bv.Questions, qv)
		}
		out.Blocks = append(out.Blocks, bv)
	}
	return out, nil
}

func paperHeadForExam(p entity.ExamPaper) PaperHeadForExamView {
	v := PaperHeadForExamView{
		Id:                 p.MockExaminationPaperId,
		Level:              p.Level,
		PaperId:            p.PaperId,
		Title:              p.Title,
		PrepareInstruction: p.PrepareInstruction,
		PrepareAudioFile:   p.PrepareAudioFile,
		SourceBaseUrl:      p.SourceBaseUrl,
		//IndexJson:          trimForExamLarge(p.IndexJson),
		DurationSeconds: p.DurationSeconds,
	}
	v.CreateTime = util.ToRFC3339UTC(p.CreateTime)
	return v
}
// PaperDetailForExamInitTree 考前初始化视图（轻量）。
type PaperDetailForExamInitTree struct {
	Paper    PaperHeadForExamView        `json:"paper"`
	Sections []SectionOutlineForExamView `json:"sections"`
}

// SectionOutlineForExamView section 概要 + block 概要（无题目）。
type SectionOutlineForExamView struct {
	Id             int64                     `json:"id"`
	SortOrder      int                       `json:"sort_order"`
	TopicTitle     string                    `json:"topic_title"`
	TopicSubtitle  string                    `json:"topic_subtitle"`
	TopicType      string                    `json:"topic_type"`
	PartCode       int                       `json:"part_code"`
	SegmentCode    string                    `json:"segment_code"`
	TopicItemsFile string                    `json:"topic_items_file"`
	TopicJson      string                    `json:"topic_json"`
	Blocks         []BlockOutlineForExamView `json:"blocks"`
}

// BlockOutlineForExamView block 概要（仅题量，无题目列表）。
type BlockOutlineForExamView struct {
	Id                      int64  `json:"id"`
	BlockOrder              int    `json:"block_order"`
	GroupIndex              int    `json:"group_index"`
	QuestionDescriptionJson string `json:"question_description_json"`
	QuestionCount           int    `json:"question_count"`
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

// SectionDetailForExamView 单个 section 完整考前树（脱敏）。
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

