package exam

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"exam/internal/consts"
	examdao "exam/internal/dao/exam"
	"exam/internal/exampaper"
	exambo "exam/internal/model/bo/exam"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gtime"
)

// ImportFromIndex 拉取或解析 index.json，写入 exam_* 表。
func (s *sExam) ImportFromIndex(ctx context.Context, p exambo.ImportParams) (*exambo.ImportResult, error) {
	res := &exambo.ImportResult{}
	if err := exampaper.EnsureMockExaminationPaper(ctx, p.MockExaminationPaperId); err != nil {
		return nil, err
	}
	indexStr, baseURL, level, paperID, err := resolveIndexPayload(ctx, p)
	if err != nil {
		return nil, err
	}
	if p.SourceBaseURL != "" {
		baseURL = strings.TrimRight(p.SourceBaseURL, "/") + "/"
	}
	audioHlsPrefix := strings.Trim(p.AudioHlsPrefix, "/")
	if p.Level != "" {
		level = p.Level
	}
	if p.PaperID != "" {
		paperID = p.PaperID
	}

	mode := p.ConflictMode
	if mode == "" {
		mode = consts.ExamImportConflictFail
	}

	idx := gjson.New(indexStr)
	title := idx.Get("title").String()
	prepareTitle := idx.Get("prepare.title").String()
	if prepareTitle == "" {
		prepareTitle = idx.Get("prepare_title").String()
	}
	prepareInstr := idx.Get("prepare.instruction").String()
	if prepareInstr == "" {
		prepareInstr = idx.Get("prepare_instruction").String()
	}
	prepareAudio := idx.Get("prepare.audio_file").String()
	if prepareAudio == "" {
		prepareAudio = idx.Get("prepare_audio_file").String()
	}

	mockID := p.MockExaminationPaperId
	var exist examentity.ExamPaper
	if err := examdao.ExamPaper.Ctx(ctx).
		Where("mock_examination_paper_id", mockID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&exist); err != nil {
		if !isNoRowsErr(err) {
			return nil, err
		}
	}
	if exist.Id > 0 {
		switch mode {
		case consts.ExamImportConflictFail:
			res.Conflict = true
			res.ExistingExaminationPaperID = mockID
			return res, nil
		case consts.ExamImportConflictOverwrite:
			// 子树伪删除 + 原地更新 exam_paper，在下方事务中执行
		case consts.ExamImportConflictNewCopy:
			if p.NewPaperID == "" {
				return nil, gerror.NewCode(consts.CodeExamNewPaperIdRequired)
			}
			paperID = p.NewPaperID
			var dup examentity.ExamPaper
			if err := examdao.ExamPaper.Ctx(ctx).
				Where("level", level).
				Where("paper_id", paperID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				WhereNot("id", exist.Id).
				Scan(&dup); err != nil {
				if !isNoRowsErr(err) {
					return nil, err
				}
			}
			if dup.Id > 0 {
				return nil, gerror.NewCode(consts.CodeExamNewPaperIdExists)
			}
		default:
			return nil, gerror.NewCode(consts.CodeExamConflictModeInvalid)
		}
	} else {
		var pathDup examentity.ExamPaper
		if err := examdao.ExamPaper.Ctx(ctx).
			Where("level", level).
			Where("paper_id", paperID).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&pathDup); err != nil {
			if !isNoRowsErr(err) {
				return nil, err
			}
		}
		if pathDup.Id > 0 {
			return nil, gerror.NewCode(consts.CodeExamNewPaperIdExists)
		}
	}

	indexSnapshot := gjson.New(indexStr).MustToJsonString()

	overwritePaperId := int64(0)
	if exist.Id > 0 && mode == consts.ExamImportConflictOverwrite {
		overwritePaperId = exist.Id
	}
	// 覆盖且未传 HLS 前缀时保留库内原值（与「伪删除保留历史」一致）
	audioHlsForPaper := audioHlsPrefix
	if overwritePaperId > 0 && strings.TrimSpace(p.AudioHlsPrefix) == "" {
		audioHlsForPaper = strings.Trim(exist.AudioHlsPrefix, "/")
	}

	var invalidatePaperPK int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var pid int64
		if overwritePaperId > 0 {
			if err := softDeletePaperTreeTx(ctx, tx, overwritePaperId, p.Creator); err != nil {
				return err
			}
			_, err := tx.Model(examdao.ExamPaper.Table()).Ctx(ctx).Where("id", overwritePaperId).Data(examdo.ExamPaper{
				Level:                   level,
				PaperId:                 paperID,
				MockExaminationPaperId:  mockID,
				Title:                   title,
				PrepareTitle:            prepareTitle,
				PrepareInstruction:      prepareInstr,
				PrepareAudioFile:        prepareAudio,
				SourceBaseUrl:           baseURL,
				AudioHlsPrefix:          audioHlsForPaper,
				AudioHlsSegmentCount:    exist.AudioHlsSegmentCount,
				AudioHlsSegmentPattern:  exist.AudioHlsSegmentPattern,
				AudioHlsKeyObject:       exist.AudioHlsKeyObject,
				AudioHlsIvHex:           exist.AudioHlsIvHex,
				AudioHlsSegmentDuration: exist.AudioHlsSegmentDuration,
				IndexJson:               indexSnapshot,
				DurationSeconds:         exist.DurationSeconds,
				Updater:                 p.Creator,
				UpdateTime:              gtime.Now(),
				DeleteFlag:              consts.DeleteFlagNotDeleted,
			}).Update()
			if err != nil {
				return err
			}
			pid = overwritePaperId
		} else {
			paperDO := examdo.ExamPaper{
				Level:                  level,
				PaperId:                paperID,
				MockExaminationPaperId: mockID,
				Title:                  title,
				PrepareTitle:           prepareTitle,
				PrepareInstruction:     prepareInstr,
				PrepareAudioFile:       prepareAudio,
				SourceBaseUrl:          baseURL,
				AudioHlsPrefix:         audioHlsForPaper,
				IndexJson:              indexSnapshot,
				Creator:                p.Creator,
				Updater:                p.Creator,
				DeleteFlag:             consts.DeleteFlagNotDeleted,
				CreateTime:             gtime.Now(),
				UpdateTime:             gtime.Now(),
			}
			inserted, err := tx.Model(examdao.ExamPaper.Table()).Ctx(ctx).InsertAndGetId(paperDO)
			if err != nil {
				return err
			}
			pid = inserted
		}
		res.ExaminationPaperID = mockID

		items := idx.Get("items").Array()
		secCount := 0
		qCount := 0
		for i, it := range items {
			item := gjson.New(it)
			topicFile := item.Get("topic_items").String()
			if topicFile == "" {
				continue
			}
			topicURL := baseURL + topicFile
			body, err := fetchRemote(ctx, topicURL)
			if err != nil {
				return gerror.Wrapf(err, "fetch topic %s", topicFile)
			}
			topicSnap := gjson.New(body).MustToJsonString()

			secDO := examdo.ExamSection{
				ExamPaperId:            pid,
				MockExaminationPaperId: mockID,
				SortOrder:              i,
				TopicTitle:             item.Get("topic_title").String(),
				TopicSubtitle:          item.Get("topic_subtitle").String(),
				TopicType:              item.Get("topic_type").String(),
				PartCode:               item.Get("part_code").Int(),
				SegmentCode:            item.Get("segment_code").String(),
				TopicItemsFile:         topicFile,
				TopicJson:              topicSnap,
				Creator:                p.Creator,
				Updater:                p.Creator,
				DeleteFlag:             consts.DeleteFlagNotDeleted,
				CreateTime:             gtime.Now(),
				UpdateTime:             gtime.Now(),
			}
			sid, err := tx.Model(examdao.ExamSection.Table()).Ctx(ctx).InsertAndGetId(secDO)
			if err != nil {
				return err
			}
			secCount++

			n, err := insertTopicContent(ctx, tx, pid, mockID, sid, body, p.Creator)
			if err != nil {
				return err
			}
			qCount += n
		}
		res.SectionCount = secCount
		res.QuestionCount = qCount
		invalidatePaperPK = pid
		return nil
	})
	if err != nil {
		return nil, err
	}
	if invalidatePaperPK > 0 {
		s.InvalidatePaperForExamCache(ctx, invalidatePaperPK)
	}
	return res, nil
}

func resolveIndexPayload(ctx context.Context, p exambo.ImportParams) (indexStr, baseURL, level, paperID string, err error) {
	if p.IndexJSON != "" {
		indexStr = p.IndexJSON
		if p.Level == "" || p.PaperID == "" {
			return "", "", "", "", gerror.NewCode(consts.CodeExamIndexMetaRequired)
		}
		baseURL = p.SourceBaseURL
		if baseURL == "" {
			return "", "", "", "", gerror.NewCode(consts.CodeExamSourceBaseRequired)
		}
		return indexStr, baseURL, p.Level, p.PaperID, nil
	}
	if p.IndexURL == "" {
		return "", "", "", "", gerror.NewCode(consts.CodeExamIndexUrlRequired)
	}
	baseURL, level, paperID, err = parseIndexURL(p.IndexURL)
	if err != nil {
		return "", "", "", "", err
	}
	indexStr, err = fetchRemote(ctx, p.IndexURL)
	if err != nil {
		return "", "", "", "", err
	}
	return indexStr, baseURL, level, paperID, nil
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

func isNoRowsErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "no rows in result set")
}

// deletePaperTree 独立事务删除（供管理端扩展等使用）。
func deletePaperTree(ctx context.Context, examPaperId int64) error {
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return deletePaperTreeTx(ctx, tx, examPaperId)
	})
	if err == nil {
		InvalidatePaperForExamCache(ctx, examPaperId)
	}
	return err
}

// softDeletePaperTreeTx 将试卷下未删除的大题/题块/试题/选项标记为逻辑删除（覆盖导入时用，保留历史数据）。
func softDeletePaperTreeTx(ctx context.Context, tx gdb.TX, examPaperId int64, updater string) error {
	now := gtime.Now()
	var qids []int64
	if err := tx.Model(examdao.ExamQuestion.Table()).Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Fields("id").
		Scan(&qids); err != nil {
		return err
	}
	if len(qids) > 0 {
		if _, err := tx.Model(examdao.ExamOption.Table()).Ctx(ctx).
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
	if _, err := tx.Model(examdao.ExamQuestion.Table()).Ctx(ctx).
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
	if err := tx.Model(examdao.ExamSection.Table()).Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Fields("id").
		Scan(&sids); err != nil {
		return err
	}
	if len(sids) > 0 {
		if _, err := tx.Model(examdao.ExamQuestionBlock.Table()).Ctx(ctx).
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
	if _, err := tx.Model(examdao.ExamSection.Table()).Ctx(ctx).
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
	if err := tx.Model(examdao.ExamQuestion.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Fields("id").Scan(&qids); err != nil {
		return err
	}
	if len(qids) > 0 {
		if _, err := tx.Model(examdao.ExamOption.Table()).Ctx(ctx).WhereIn("question_id", qids).Delete(); err != nil {
			return err
		}
	}
	if _, err := tx.Model(examdao.ExamQuestion.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Delete(); err != nil {
		return err
	}
	var sids []int64
	if err := tx.Model(examdao.ExamSection.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Fields("id").Scan(&sids); err != nil {
		return err
	}
	if len(sids) > 0 {
		if _, err := tx.Model(examdao.ExamQuestionBlock.Table()).Ctx(ctx).WhereIn("section_id", sids).Delete(); err != nil {
			return err
		}
	}
	if _, err := tx.Model(examdao.ExamSection.Table()).Ctx(ctx).Where("exam_paper_id", examPaperId).Delete(); err != nil {
		return err
	}
	if _, err := tx.Model(examdao.ExamPaper.Table()).Ctx(ctx).Where("id", examPaperId).Delete(); err != nil {
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
			bid, err := tx.Model(examdao.ExamQuestionBlock.Table()).Ctx(ctx).InsertAndGetId(blockDO)
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
		bid, err := tx.Model(examdao.ExamQuestionBlock.Table()).Ctx(ctx).InsertAndGetId(blockDO)
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

// classifySubjective 导入时标记主观题：无选项、无标准答案、或题干 type 为写作类（不以大题 topic_type 一刀切，避免写作部分内的客观题被误判）。
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
	qid, err := tx.Model(examdao.ExamQuestion.Table()).Ctx(ctx).InsertAndGetId(qDO)
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
		if _, err := tx.Model(examdao.ExamOption.Table()).Ctx(ctx).Insert(optDO); err != nil {
			return 0, err
		}
	}
	return 1, nil
}

// PaperList 分页试卷列表（管理端）
func (s *sExam) PaperList(ctx context.Context, page, size int, level string) (list []examentity.ExamPaper, total int, err error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	m := examdao.ExamPaper.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if level != "" {
		m = m.Where("level", level)
	}
	n, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	total = n
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
