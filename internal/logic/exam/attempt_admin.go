package exam

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/examutil"
	"exam/internal/model/bo"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"
)

// AttemptAdminList 分页查询答题会话（联表学员、试卷）。
func (s *sExam) AttemptAdminList(ctx context.Context, page, size int, level string, examinationPaperId int64, examBatchId int64, status int, username string) ([]bo.AttemptAdminListRow, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var where strings.Builder
	where.WriteString("r.delete_flag = ?")
	args := []interface{}{consts.DeleteFlagNotDeleted}
	if level != "" {
		where.WriteString(" AND p.level = ?")
		args = append(args, level)
	}
	if examinationPaperId > 0 {
		where.WriteString(" AND p.mock_examination_paper_id = ?")
		args = append(args, examinationPaperId)
	}
	if examBatchId > 0 {
		where.WriteString(" AND r.exam_batch_id = ?")
		args = append(args, examBatchId)
	}
	if status > 0 {
		where.WriteString(" AND r.status = ?")
		args = append(args, status)
	}
	if username != "" {
		where.WriteString(" AND u.username LIKE ?")
		args = append(args, "%"+username+"%")
	}
	w := where.String()
	from := ` FROM exam_result r
INNER JOIN exam_attempt a ON a.id = r.attempt_id AND a.delete_flag = ?
LEFT JOIN sys_member u ON u.id = r.member_id AND u.delete_flag = ?
LEFT JOIN exam_paper p ON p.id = r.exam_paper_id AND p.delete_flag = ?
LEFT JOIN mock_examination_paper m ON m.id = p.mock_examination_paper_id AND m.delete_flag = ?
LEFT JOIN mock_levels ml ON ml.id = r.mock_level_id AND ml.delete_flag = ?
WHERE ` + w
	joinArgs := []interface{}{
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
		consts.DeleteFlagNotDeleted,
	}

	countSQL := "SELECT COUNT(1) AS total" + from
	countArgs := append(append([]interface{}{}, joinArgs...), args...)
	var cnt struct {
		Total int `json:"total"`
	}
	if err := g.DB().Ctx(ctx).Raw(countSQL, countArgs...).Scan(&cnt); err != nil {
		return nil, 0, err
	}
	total := cnt.Total
	if total == 0 {
		return nil, 0, nil
	}
	offset := (page - 1) * size
	listSQL := `SELECT r.attempt_id AS id, r.member_id, IFNULL(p.mock_examination_paper_id,0) AS examination_paper_id,
IFNULL(r.exam_batch_id,0) AS exam_batch_id, IFNULL(r.mock_level_id,0) AS mock_level_id, r.status,
r.objective_score, r.subjective_score, r.total_score, r.has_subjective,
a.started_at, a.submitted_at, a.ended_at, a.create_time,
IFNULL(u.username,'') AS username, IFNULL(u.nickname,'') AS nickname,
	IFNULL(TRIM(IFNULL(m.name,'')), '') AS paper_title,
	COALESCE(NULLIF(TRIM(IFNULL(ml.level_name,'')), ''), IFNULL(p.level,'')) AS paper_level,
	IFNULL(p.paper_id,'') AS remote_paper_id` +
		from + ` ORDER BY r.attempt_id DESC LIMIT ? OFFSET ?`
	listArgs := append(append([]interface{}{}, joinArgs...), args...)
	listArgs = append(listArgs, size, offset)

	var rows []bo.AttemptAdminListRow
	if err := g.DB().Ctx(ctx).Raw(listSQL, listArgs...).Scan(&rows); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// AttemptAdminDetail 按 id 加载会话、学员、试卷及答题明细（含客观题是否选对）。
func (s *sExam) AttemptAdminDetail(ctx context.Context, attemptID int64) (*bo.AttemptAdminDetailView, error) {
	if attemptID <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	var att examentity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att); err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.NewCode(consts.CodeExamAttemptNotFound)
	}
	var user sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).
		Where("id", att.MemberId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&user); err != nil {
		return nil, err
	}
	var paper examentity.ExamPaper
	if err := dao.ExamPaper.Ctx(ctx).
		Where("id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper); err != nil {
		return nil, err
	}

	var ansRows []examentity.ExamAttemptAnswer
	if err := dao.ExamAttemptAnswer.Ctx(ctx).
		Where("attempt_id", att.Id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&ansRows); err != nil {
		return nil, err
	}
	qIDs := make([]interface{}, 0, len(ansRows))
	for _, a := range ansRows {
		qIDs = append(qIDs, a.ExamQuestionId)
	}
	qByID := make(map[int64]examentity.ExamQuestion)
	if len(qIDs) > 0 {
		var qs []examentity.ExamQuestion
		if err := dao.ExamQuestion.Ctx(ctx).
			Where("exam_paper_id", att.ExamPaperId).
			WhereIn("id", qIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&qs); err != nil {
			return nil, err
		}
		for _, q := range qs {
			qByID[q.Id] = q
		}
	}

	// 预加载块与 section 信息，便于按 section 展示
	blockIDs := make([]interface{}, 0, len(qByID))
	for _, q := range qByID {
		blockIDs = append(blockIDs, q.BlockId)
	}
	blockByID := loadBlocksByID(ctx, blockIDs)
	sectionByID := loadSectionsByID(ctx, att.ExamPaperId, blockByID)

	correctByQ := loadCorrectOptionIDsByQuestion(ctx, qIDs)
	optionsByQ := loadOptionsByQuestion(ctx, qIDs)

	out := &bo.AttemptAdminDetailView{
		Attempt: att,
		User:    user,
		Paper:   paper,
		Answers: make([]bo.AttemptAdminAnswerRow, 0, len(ansRows)),
	}
	for _, ar := range ansRows {
		q := qByID[ar.ExamQuestionId]
		var secPtr *examentity.ExamSection
		if blk, ok := blockByID[q.BlockId]; ok {
			if sec, ok2 := sectionByID[blk.SectionId]; ok2 {
				// 复制一份，避免后续误修改原始 map
				sv := sec
				secPtr = &sv
			}
		}
		row := bo.AttemptAdminAnswerRow{
			Answer:   ar,
			Question: q,
			Section:  secPtr,
			Options:  optionsByQ[ar.ExamQuestionId],
		}
		if q.Id != 0 && q.IsExample == 0 && q.IsSubjective == 0 {
			payload := examutil.ParseAnswerPayload(ar.AnswerJson)
			ok := examutil.ObjectiveAnswerCorrect(correctByQ[q.Id], payload.SelectedOptionIDs)
			row.ObjectiveCorrect = boolPtr(ok)
		}
		out.Answers = append(out.Answers, row)
	}
	return out, nil
}

// AttemptAdminSaveSubjectiveScores 写入主观题人工分并汇总 subjective_score、total_score（允许部分题目已评）。
func (s *sExam) AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []bo.SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error) {
	if attemptID <= 0 {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	if len(items) == 0 {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	var att examentity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att); err != nil {
		return 0, 0, err
	}
	if att.Id == 0 {
		return 0, 0, gerror.NewCode(consts.CodeExamAttemptNotFound)
	}
	if att.Status != consts.ExamAttemptEnded {
		return 0, 0, gerror.NewCode(consts.CodeExamAttemptNotEnded)
	}
	if att.HasSubjective != 1 {
		return 0, 0, gerror.NewCode(consts.CodeExamAttemptNoSubjective)
	}

	byQ := make(map[int64]bo.SubjectiveScoreItem, len(items))
	for _, it := range items {
		byQ[it.QuestionID] = it
	}
	uniq := make([]bo.SubjectiveScoreItem, 0, len(byQ))
	for _, it := range byQ {
		uniq = append(uniq, it)
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var attTx examentity.ExamAttempt
		if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&attTx); err != nil {
			return err
		}
		if attTx.Id == 0 {
			return gerror.NewCode(consts.CodeExamAttemptNotFound)
		}
		if attTx.Status != consts.ExamAttemptEnded {
			return gerror.NewCode(consts.CodeExamAttemptNotEnded)
		}
		if attTx.HasSubjective != 1 {
			return gerror.NewCode(consts.CodeExamAttemptNoSubjective)
		}
		paperID := attTx.ExamPaperId

		for _, it := range uniq {
			var q examentity.ExamQuestion
			if err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).
				Where("id", it.QuestionID).
				Where("exam_paper_id", paperID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&q); err != nil {
				return err
			}
			if q.Id == 0 || q.IsSubjective != 1 || q.IsExample != 0 {
				return gerror.NewCode(consts.CodeInvalidParams)
			}
			if it.Score < 0 || it.Score > q.Score {
				return gerror.NewCode(consts.CodeInvalidParams)
			}

			var row examentity.ExamAttemptAnswer
			if err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
				Where("attempt_id", attemptID).
				Where("exam_question_id", it.QuestionID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&row); err != nil {
				return err
			}
			scoreVal := it.Score
			if row.Id == 0 {
				if _, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Insert(examdo.ExamAttemptAnswer{
					AttemptId:      attemptID,
					ExamQuestionId: it.QuestionID,
					AnswerJson:     "{}",
					AwardedScore:   &scoreVal,
					Version:        0,
					Creator:        "admin",
					Updater:        "admin",
					DeleteFlag:     consts.DeleteFlagNotDeleted,
					CreateTime:     gtime.Now(),
					UpdateTime:     gtime.Now(),
				}); err != nil {
					return err
				}
			} else {
				if _, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Where("id", row.Id).Update(examdo.ExamAttemptAnswer{
					AwardedScore: &scoreVal,
					Updater:      "admin",
					UpdateTime:   gtime.Now(),
				}); err != nil {
					return err
				}
			}
		}

		sum, err := sumSubjectiveAwardedTx(ctx, tx, attemptID, paperID)
		if err != nil {
			return err
		}
		obj := attTx.ObjectiveScore
		tot := obj + sum
		if _, err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Update(examdo.ExamAttempt{
			SubjectiveScore: sum,
			TotalScore:      tot,
			Updater:         "admin",
			UpdateTime:      gtime.Now(),
		}); err != nil {
			return err
		}
		if err := examutil.UpsertFromAttemptTx(ctx, tx, attemptID); err != nil {
			return err
		}
		subjectiveSum = sum
		totalScore = tot
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return subjectiveSum, totalScore, nil
}

func sumSubjectiveAwardedTx(ctx context.Context, tx gdb.TX, attemptID int64, paperID int64) (float64, error) {
	type qrow struct {
		Id int64 `json:"id"`
	}
	var qrows []qrow
	if err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).
		Fields("id").
		Where("exam_paper_id", paperID).
		Where("is_subjective", 1).
		Where("is_example", 0).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&qrows); err != nil {
		return 0, err
	}
	if len(qrows) == 0 {
		return 0, nil
	}
	subjSet := make(map[int64]struct{}, len(qrows))
	for _, q := range qrows {
		subjSet[q.Id] = struct{}{}
	}
	type ansRow struct {
		ExamQuestionId int64    `json:"exam_question_id"`
		AwardedScore   *float64 `json:"awarded_score"`
	}
	var ansRows []ansRow
	if err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
		Fields("exam_question_id", "awarded_score").
		Where("attempt_id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&ansRows); err != nil {
		return 0, err
	}
	var sum float64
	for _, a := range ansRows {
		if _, ok := subjSet[a.ExamQuestionId]; !ok {
			continue
		}
		if a.AwardedScore == nil {
			continue
		}
		sum += *a.AwardedScore
	}
	return sum, nil
}

func boolPtr(b bool) *bool {
	return &b
}

func loadCorrectOptionIDsByQuestion(ctx context.Context, qIDs []interface{}) map[int64][]int64 {
	out := make(map[int64][]int64)
	if len(qIDs) == 0 {
		return out
	}
	var opts []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		WhereIn("question_id", qIDs).
		Where("is_correct", 1).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&opts); err != nil {
		g.Log().Warningf(ctx, "loadCorrectOptionIDsByQuestion: %v", err)
		return out
	}
	for _, o := range opts {
		out[o.QuestionId] = append(out[o.QuestionId], o.Id)
	}
	return out
}

func loadBlocksByID(ctx context.Context, blockIDs []interface{}) map[int64]examentity.ExamQuestionBlock {
	out := make(map[int64]examentity.ExamQuestionBlock)
	if len(blockIDs) == 0 {
		return out
	}
	var blocks []examentity.ExamQuestionBlock
	if err := dao.ExamQuestionBlock.Ctx(ctx).
		WhereIn("id", blockIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&blocks); err != nil {
		g.Log().Warningf(ctx, "loadBlocksByID: %v", err)
		return out
	}
	for _, b := range blocks {
		out[b.Id] = b
	}
	return out
}

func loadSectionsByID(ctx context.Context, examPaperId int64, blocks map[int64]examentity.ExamQuestionBlock) map[int64]examentity.ExamSection {
	out := make(map[int64]examentity.ExamSection)
	if len(blocks) == 0 {
		return out
	}
	sectionIDs := make([]interface{}, 0, len(blocks))
	seen := make(map[int64]struct{})
	for _, b := range blocks {
		if _, ok := seen[b.SectionId]; ok {
			continue
		}
		seen[b.SectionId] = struct{}{}
		sectionIDs = append(sectionIDs, b.SectionId)
	}
	if len(sectionIDs) == 0 {
		return out
	}
	var sections []examentity.ExamSection
	if err := dao.ExamSection.Ctx(ctx).
		Where("exam_paper_id", examPaperId).
		WhereIn("id", sectionIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&sections); err != nil {
		g.Log().Warningf(ctx, "loadSectionsByID: %v", err)
		return out
	}
	for _, s := range sections {
		out[s.Id] = s
	}
	return out
}

func loadOptionsByQuestion(ctx context.Context, qIDs []interface{}) map[int64][]examentity.ExamOption {
	out := make(map[int64][]examentity.ExamOption)
	if len(qIDs) == 0 {
		return out
	}
	var opts []examentity.ExamOption
	if err := dao.ExamOption.Ctx(ctx).
		WhereIn("question_id", qIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort_order").
		Scan(&opts); err != nil {
		g.Log().Warningf(ctx, "loadOptionsByQuestion: %v", err)
		return out
	}
	for _, o := range opts {
		qid := o.QuestionId
		out[qid] = append(out[qid], o)
	}
	return out
}
