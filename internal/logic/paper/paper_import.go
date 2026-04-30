package paper

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"exam/internal/consts"
	"exam/internal/dao"
	exambo "exam/internal/model/bo/exam"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gtime"
)

type importIndexSource struct {
	mockPaper      mockentity.MockExaminationPaper
	mockID         int64
	baseURL        string
	level          string
	paperID        string
	indexJSON      *gjson.Json
	indexSnapshot  string
	title          string
	prepareTitle   string
	prepareInstr   string
	prepareAudio   string
	audioHlsPrefix string
}

type importConflictPlan struct {
	existing         []examentity.ExamPaper
	overwritePaperID int64
	conflict         bool
}

// ImportFromIndex 根据 mock_examination_paper.resource_url 推导 index.json 并导入 exam_* 表。
func (s *sPaper) ImportFromIndex(ctx context.Context, p exambo.ImportParams) (*exambo.ImportResult, error) {
	res := &exambo.ImportResult{}

	source, err := loadImportIndexSource(ctx, p)
	if err != nil {
		return nil, err
	}

	plan, err := buildImportConflictPlan(ctx, p.ConflictMode, source.mockID, p.OverwriteExamPaperId)
	if err != nil {
		return nil, err
	}
	if plan.conflict {
		res.Conflict = true
		res.ExistingMockExaminationPaperID = source.mockID
		return res, nil
	}

	var importedExamPaperID int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		paperID, err := upsertImportedPaper(ctx, tx, p.Creator, source, plan)
		if err != nil {
			return err
		}

		sectionCount, questionCount, err := importPaperSections(ctx, tx, paperID, source, p.Creator)
		if err != nil {
			return err
		}

		res.MockExaminationPaperID = source.mockID
		res.SectionCount = sectionCount
		res.QuestionCount = questionCount
		importedExamPaperID = paperID
		return nil
	})
	if err != nil {
		return nil, err
	}

	invalidatePaperCaches(ctx, importedExamPaperID, source.mockID)
	return res, nil
}

func loadImportIndexSource(ctx context.Context, p exambo.ImportParams) (importIndexSource, error) {
	var source importIndexSource
	if err := dao.MockExaminationPaper.Ctx(ctx).
		Where(dao.MockExaminationPaper.Columns().Id, p.MockExaminationPaperId).
		Where(dao.MockExaminationPaper.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
		Scan(&source.mockPaper); err != nil {
		return source, err
	}
	if source.mockPaper.Id == 0 {
		return source, gerror.NewCode(consts.CodeMockExamPaperNotFound)
	}

	indexURL, err := indexJSONURLFromMockResourceURL(source.mockPaper.ResourceUrl)
	if err != nil {
		return source, err
	}
	indexBody, err := fetchRemote(ctx, indexURL)
	if err != nil {
		return source, err
	}
	baseURL, level, paperID, err := parseIndexURL(indexURL)
	if err != nil {
		return source, err
	}

	idx := gjson.New(indexBody)
	source = importIndexSource{
		mockPaper:      source.mockPaper,
		mockID:         p.MockExaminationPaperId,
		baseURL:        baseURL,
		level:          level,
		paperID:        paperID,
		indexJSON:      idx,
		indexSnapshot:  idx.MustToJsonString(),
		audioHlsPrefix: strings.Trim(p.AudioHlsPrefix, "/"),
		title: firstNonEmpty(
			strings.TrimSpace(p.Title),
			strings.TrimSpace(source.mockPaper.Name),
			strings.TrimSpace(idx.Get("title").String()),
		),
		prepareTitle: firstNonEmpty(idx.Get("prepare.title").String(), idx.Get("prepare_title").String()),
		prepareInstr: firstNonEmpty(idx.Get("prepare.instruction").String(), idx.Get("prepare_instruction").String()),
		prepareAudio: firstNonEmpty(idx.Get("prepare.audio_file").String(), idx.Get("prepare_audio_file").String()),
	}
	return source, nil
}

func buildImportConflictPlan(ctx context.Context, conflictMode string, mockID, overwriteExamPaperID int64) (importConflictPlan, error) {
	var plan importConflictPlan
	mode, err := normalizeExamImportConflictMode(conflictMode)
	if err != nil {
		return plan, err
	}

	if err := dao.ExamPaper.Ctx(ctx).
		Where("mock_examination_paper_id", mockID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&plan.existing); err != nil {
		return plan, err
	}
	if len(plan.existing) == 0 {
		return plan, nil
	}
	if mode == consts.ExamImportConflictFail {
		plan.conflict = true
		return plan, nil
	}
	if mode == consts.ExamImportConflictNew {
		return plan, nil
	}
	if overwriteExamPaperID <= 0 {
		return plan, gerror.NewCode(consts.CodeInvalidParams)
	}
	for _, paper := range plan.existing {
		if paper.Id == overwriteExamPaperID {
			plan.overwritePaperID = overwriteExamPaperID
			return plan, nil
		}
	}
	return plan, gerror.NewCode(consts.CodeExamPaperNotFound)
}

func getOverwriteExisting(plan importConflictPlan) examentity.ExamPaper {
	for _, paper := range plan.existing {
		if paper.Id == plan.overwritePaperID {
			return paper
		}
	}
	return examentity.ExamPaper{}
}

func upsertImportedPaper(ctx context.Context, tx gdb.TX, creator string, source importIndexSource, plan importConflictPlan) (int64, error) {
	audioHlsPrefix := resolveImportAudioHlsPrefix(source.audioHlsPrefix, plan)

	if plan.overwritePaperID > 0 {
		existing := getOverwriteExisting(plan)
		if existing.Id == 0 {
			return 0, gerror.NewCode(consts.CodeExamPaperNotFound)
		}
		if err := softDeletePaperTreeTx(ctx, tx, plan.overwritePaperID, creator); err != nil {
			return 0, err
		}
		_, err := tx.Model(dao.ExamPaper.Table()).Ctx(ctx).Where("id", plan.overwritePaperID).Data(examdo.ExamPaper{
			Level:                   source.level,
			PaperId:                 source.paperID,
			MockExaminationPaperId:  source.mockID,
			Title:                   source.title,
			PrepareTitle:            source.prepareTitle,
			PrepareInstruction:      source.prepareInstr,
			PrepareAudioFile:        source.prepareAudio,
			SourceBaseUrl:           source.baseURL,
			AudioHlsPrefix:          audioHlsPrefix,
			AudioHlsSegmentCount:    existing.AudioHlsSegmentCount,
			AudioHlsSegmentPattern:  existing.AudioHlsSegmentPattern,
			AudioHlsKeyObject:       existing.AudioHlsKeyObject,
			AudioHlsIvHex:           existing.AudioHlsIvHex,
			AudioHlsSegmentDuration: existing.AudioHlsSegmentDuration,
			IndexJson:               source.indexSnapshot,
			DurationSeconds:         existing.DurationSeconds,
			Updater:                 creator,
			UpdateTime:              gtime.Now(),
			DeleteFlag:              consts.DeleteFlagNotDeleted,
		}).Update()
		if err != nil {
			return 0, err
		}
		return plan.overwritePaperID, nil
	}

	inserted, err := tx.Model(dao.ExamPaper.Table()).Ctx(ctx).InsertAndGetId(examdo.ExamPaper{
		Level:                  source.level,
		PaperId:                source.paperID,
		MockExaminationPaperId: source.mockID,
		Title:                  source.title,
		PrepareTitle:           source.prepareTitle,
		PrepareInstruction:     source.prepareInstr,
		PrepareAudioFile:       source.prepareAudio,
		SourceBaseUrl:          source.baseURL,
		AudioHlsPrefix:         audioHlsPrefix,
		IndexJson:              source.indexSnapshot,
		Creator:                creator,
		Updater:                creator,
		DeleteFlag:             consts.DeleteFlagNotDeleted,
		CreateTime:             gtime.Now(),
		UpdateTime:             gtime.Now(),
	})
	if err != nil {
		return 0, err
	}
	return inserted, nil
}

func importPaperSections(ctx context.Context, tx gdb.TX, examPaperID int64, source importIndexSource, creator string) (sectionCount int, questionCount int, err error) {
	items := source.indexJSON.Get("items").Array()
	for i, it := range items {
		item := gjson.New(it)
		topicFile := item.Get("topic_items").String()
		if topicFile == "" {
			continue
		}

		topicBody, err := fetchRemote(ctx, source.baseURL+topicFile)
		if err != nil {
			return 0, 0, gerror.Wrapf(err, "fetch topic %s", topicFile)
		}

		sectionID, err := tx.Model(dao.ExamSection.Table()).Ctx(ctx).InsertAndGetId(examdo.ExamSection{
			ExamPaperId:            examPaperID,
			MockExaminationPaperId: source.mockID,
			SortOrder:              i,
			TopicTitle:             item.Get("topic_title").String(),
			TopicSubtitle:          item.Get("topic_subtitle").String(),
			TopicType:              item.Get("topic_type").String(),
			PartCode:               item.Get("part_code").Int(),
			SegmentCode:            item.Get("segment_code").String(),
			TopicItemsFile:         topicFile,
			TopicJson:              gjson.New(topicBody).MustToJsonString(),
			Creator:                creator,
			Updater:                creator,
			DeleteFlag:             consts.DeleteFlagNotDeleted,
			CreateTime:             gtime.Now(),
			UpdateTime:             gtime.Now(),
		})
		if err != nil {
			return 0, 0, err
		}
		sectionCount++

		n, err := insertTopicContent(ctx, tx, examPaperID, source.mockID, sectionID, topicBody, creator)
		if err != nil {
			return 0, 0, err
		}
		questionCount += n
	}
	return sectionCount, questionCount, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func resolveImportAudioHlsPrefix(requestPrefix string, plan importConflictPlan) string {
	prefix := strings.Trim(requestPrefix, "/")
	if plan.overwritePaperID > 0 && prefix == "" {
		existing := getOverwriteExisting(plan)
		return strings.Trim(existing.AudioHlsPrefix, "/")
	}
	return prefix
}

func normalizeExamImportConflictMode(mode string) (string, error) {
	m := strings.ToLower(strings.TrimSpace(mode))
	if m == "" {
		return consts.ExamImportConflictFail, nil
	}
	if m == "new_copy" {
		m = consts.ExamImportConflictNew
	}
	switch m {
	case consts.ExamImportConflictFail, consts.ExamImportConflictOverwrite, consts.ExamImportConflictNew:
		return m, nil
	default:
		return "", gerror.NewCode(consts.CodeExamConflictModeInvalid)
	}
}

func indexJSONURLFromMockResourceURL(resource string) (string, error) {
	s := strings.TrimSpace(resource)
	if s == "" {
		return "", gerror.NewCode(consts.CodeInvalidParams)
	}
	parsed, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}
	path := parsed.Path
	lp := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lp, ".zip"):
		parsed.Path = path[:len(path)-len(".zip")] + "/index.json"
	case strings.HasSuffix(lp, "/index.json"):
		// keep
	default:
		parsed.Path = strings.TrimSuffix(path, "/") + "/index.json"
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	out := parsed.String()
	if out == "" {
		return "", gerror.NewCode(consts.CodeInvalidParams)
	}
	return out, nil
}

func parseIndexURL(raw string) (baseURL, level, paperID string, err error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", "", err
	}
	path := strings.Trim(parsed.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", "", gerror.NewCode(consts.CodeExamIndexUrlInvalid)
	}
	if parts[len(parts)-1] != "index.json" {
		return "", "", "", gerror.NewCode(consts.CodeExamIndexUrlInvalid)
	}
	basePath := strings.Join(parts[:len(parts)-1], "/")
	if len(parts) < 4 {
		return "", "", "", gerror.NewCode(consts.CodeExamIndexUrlInvalid)
	}
	level = parts[len(parts)-3]
	paperID = parts[len(parts)-2]
	scheme := parsed.Scheme
	if scheme == "" {
		scheme = "https"
	}
	baseURL = fmt.Sprintf("%s://%s/%s/", scheme, parsed.Host, basePath)
	return baseURL, level, paperID, nil
}

func fetchRemote(ctx context.Context, u string) (string, error) {
	c := gclient.New()
	c.SetTimeout(90 * time.Second)
	r, err := c.Get(ctx, u)
	if err != nil {
		return "", err
	}
	defer r.Close()
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return "", gerror.Newf("http %d for %s", r.StatusCode, u)
	}
	return r.ReadAllString(), nil
}

func deletePaperTree(ctx context.Context, examPaperId int64) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return deletePaperTreeTx(ctx, tx, examPaperId)
	})
	if err == nil {
		invalidatePaperCaches(ctx, examPaperId, 0)
	}
	return err
}

func softDeletePaperTreeTx(ctx context.Context, tx gdb.TX, examPaperId int64, updater string) error {
	now := gtime.Now()
	var qids []int64
	if err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Fields("id").
		Scan(&qids); err != nil {
		return err
	}
	if len(qids) > 0 {
		if _, err := tx.Model(dao.ExamOption.Table()).Ctx(ctx).
			WhereIn("question_id", qids).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(g.Map{
				"delete_flag": consts.DeleteFlagDeleted,
				"updater":     updater,
				"update_time": now,
			}).Update(); err != nil {
			return err
		}
	}
	if _, err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(g.Map{
			"delete_flag": consts.DeleteFlagDeleted,
			"updater":     updater,
			"update_time": now,
		}).Update(); err != nil {
		return err
	}
	var sids []int64
	if err := tx.Model(dao.ExamSection.Table()).Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Fields("id").
		Scan(&sids); err != nil {
		return err
	}
	if len(sids) > 0 {
		if _, err := tx.Model(dao.ExamQuestionBlock.Table()).Ctx(ctx).
			WhereIn("section_id", sids).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(g.Map{
				"delete_flag": consts.DeleteFlagDeleted,
				"updater":     updater,
				"update_time": now,
			}).Update(); err != nil {
			return err
		}
	}
	if _, err := tx.Model(dao.ExamSection.Table()).Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(g.Map{
			"delete_flag": consts.DeleteFlagDeleted,
			"updater":     updater,
			"update_time": now,
		}).Update(); err != nil {
		return err
	}
	return nil
}

func deletePaperTreeTx(ctx context.Context, tx gdb.TX, examPaperId int64) error {
	var qids []int64
	if err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Fields("id").Scan(&qids); err != nil {
		return err
	}
	if len(qids) > 0 {
		if _, err := tx.Model(dao.ExamOption.Table()).Ctx(ctx).WhereIn("question_id", qids).Delete(); err != nil {
			return err
		}
	}
	if _, err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Delete(); err != nil {
		return err
	}
	var sids []int64
	if err := tx.Model(dao.ExamSection.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Fields("id").Scan(&sids); err != nil {
		return err
	}
	if len(sids) > 0 {
		if _, err := tx.Model(dao.ExamQuestionBlock.Table()).Ctx(ctx).WhereIn("section_id", sids).Delete(); err != nil {
			return err
		}
	}
	if _, err := tx.Model(dao.ExamSection.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Delete(); err != nil {
		return err
	}
	if _, err := tx.Model(dao.ExamPaper.Table()).Ctx(ctx).Where("id", examPaperId).Delete(); err != nil {
		return err
	}
	return nil
}

func insertTopicContent(ctx context.Context, tx gdb.TX, examPaperId, mockPaperID, sectionId int64, topicBody string, creator string) (questionCount int, err error) {
	tj := gjson.New(topicBody)
	arr := tj.Get("items").Array()
	for blockIdx, raw := range arr {
		blk := gjson.New(raw)
		qs := blk.Get("questions").Array()
		if len(qs) > 0 {
			blockDO := examdo.ExamQuestionBlock{
				SectionId:               sectionId,
				BlockOrder:              blockIdx,
				GroupIndex:              nilIfZero(blk.Get("index").Int()),
				QuestionDescriptionJson: jsonNodeToDB(blk.GetJson("question_description_obj")),
				Creator:                 creator,
				Updater:                 creator,
				DeleteFlag:              consts.DeleteFlagNotDeleted,
				CreateTime:              gtime.Now(),
				UpdateTime:              gtime.Now(),
			}
			bid, err := tx.Model(dao.ExamQuestionBlock.Table()).Ctx(ctx).InsertAndGetId(blockDO)
			if err != nil {
				return 0, err
			}
			for qi, qv := range qs {
				n, err := insertOneQuestion(ctx, tx, examPaperId, mockPaperID, bid, qi, gjson.New(qv), creator)
				if err != nil {
					return 0, err
				}
				questionCount += n
			}
			continue
		}
		blockDO := examdo.ExamQuestionBlock{
			SectionId:               sectionId,
			BlockOrder:              blockIdx,
			GroupIndex:              nil,
			QuestionDescriptionJson: nil,
			Creator:                 creator,
			Updater:                 creator,
			DeleteFlag:              consts.DeleteFlagNotDeleted,
			CreateTime:              gtime.Now(),
			UpdateTime:              gtime.Now(),
		}
		bid, err := tx.Model(dao.ExamQuestionBlock.Table()).Ctx(ctx).InsertAndGetId(blockDO)
		if err != nil {
			return 0, err
		}
		n, err := insertOneQuestion(ctx, tx, examPaperId, mockPaperID, bid, 0, blk, creator)
		if err != nil {
			return 0, err
		}
		questionCount += n
	}
	return questionCount, nil
}

func nilIfZero(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func jsonNodeToDB(j *gjson.Json) interface{} {
	if j == nil {
		return nil
	}
	s := j.MustToJsonString()
	if s == "" || s == "null" {
		return nil
	}
	return s
}

func classifySubjective(q *gjson.Json, hasCorrectOption bool) int {
	if q.Get("is_example").Bool() {
		return 0
	}
	if len(q.Get("answers").Array()) == 0 {
		return 1
	}
	if !hasCorrectOption {
		return 1
	}
	ct := strings.ToLower(strings.TrimSpace(q.Get("type").String()))
	if ct == "writing" || ct == "essay" || ct == "composition" {
		return 1
	}
	return 0
}

func insertOneQuestion(ctx context.Context, tx gdb.TX, examPaperId, mockPaperID, blockId int64, sortInBlock int, q *gjson.Json, creator string) (int, error) {
	raw := q.MustToJsonString()
	score := q.Get("score").Float64()
	qno := q.Get("index").Int()
	isEx := 0
	if q.Get("is_example").Bool() {
		isEx = 1
	}
	hasCorrectOption := false
	for _, av := range q.Get("answers").Array() {
		if gjson.New(av).Get("correct").Bool() {
			hasCorrectOption = true
			break
		}
	}
	isSubjective := classifySubjective(q, hasCorrectOption)
	qDO := examdo.ExamQuestion{
		ExamPaperId:             examPaperId,
		MockExaminationPaperId:  mockPaperID,
		BlockId:                 blockId,
		SortInBlock:             sortInBlock,
		QuestionNo:              qno,
		Score:                   score,
		IsExample:               isEx,
		IsSubjective:            isSubjective,
		ContentType:             q.Get("type").String(),
		AudioFile:               q.Get("content").String(),
		StemText:                q.Get("content_sentence").String(),
		ScreenTextJson:          jsonNodeToDB(q.GetJson("screen_text")),
		AnalysisJson:            jsonNodeToDB(q.GetJson("analysis")),
		QuestionDescriptionJson: jsonNodeToDB(q.GetJson("question_description_obj")),
		RawJson:                 raw,
		Creator:                 creator,
		Updater:                 creator,
		DeleteFlag:              consts.DeleteFlagNotDeleted,
		CreateTime:              gtime.Now(),
		UpdateTime:              gtime.Now(),
	}
	qid, err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).InsertAndGetId(qDO)
	if err != nil {
		return 0, err
	}
	ans := q.Get("answers").Array()
	for _, av := range ans {
		a := gjson.New(av)
		correct := 0
		if a.Get("correct").Bool() {
			correct = 1
		}
		optDO := examdo.ExamOption{
			QuestionId: qid,
			Flag:       a.Get("flag").String(),
			SortOrder:  a.Get("index").Int(),
			IsCorrect:  correct,
			OptionType: a.Get("type").String(),
			Content:    a.Get("content").String(),
			Creator:    creator,
			Updater:    creator,
			DeleteFlag: consts.DeleteFlagNotDeleted,
			CreateTime: gtime.Now(),
			UpdateTime: gtime.Now(),
		}
		if _, err := tx.Model(dao.ExamOption.Table()).Ctx(ctx).Insert(optDO); err != nil {
			return 0, err
		}
	}
	return 1, nil
}
