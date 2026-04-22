package attempt

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	examdo "exam/internal/model/do/exam"
	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/utility/examutil"
)

// AttemptAdminList 分页查询答题会话（联表学员、试卷）。
func (s *sAttempt) AttemptAdminList(ctx context.Context, page, size int, level string, examinationPaperId int64, examBatchId int64, status int, username string, subjectivePending int) ([]bo.AttemptAdminListRow, int, error) {
	page, size = s.getPageSize(page, size)
	q := AttemptAdminListQuery{
		Level:              level,
		ExaminationPaperId: examinationPaperId,
		ExamBatchId:        examBatchId,
		Status:             status,
		Username:           username,
		SubjectivePending:  subjectivePending,
	}
	from, joinArgs, whereArgs := q.buildAttemptAdminListFrom()
	countSQL := "SELECT COUNT(1) AS total" + from
	countArgs := attemptAdminListCountArgs(joinArgs, whereArgs)
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
	// 列顺序须与 model/bo.AttemptAdminListRow 字段一致，否则 Raw().Scan 错位导致 subjective_graded 等乱值
	listSQL := `SELECT r.attempt_id AS id, r.member_id, IFNULL(p.mock_examination_paper_id,0) AS examination_paper_id,
IFNULL(r.exam_batch_id,0) AS exam_batch_id, IFNULL(r.mock_level_id,0) AS mock_level_id, r.status,
r.objective_score, r.subjective_score, r.total_score, r.has_subjective,
IFNULL(subj_gr.has_subjective_graded, 0) AS subjective_graded,
a.started_at, a.submitted_at, a.ended_at, a.create_time,
IFNULL(u.username,'') AS username, IFNULL(u.nickname,'') AS nickname,
	IFNULL(TRIM(IFNULL(m.name,'')), '') AS paper_title,
	COALESCE(NULLIF(TRIM(IFNULL(ml.level_name,'')), ''), IFNULL(p.level,'')) AS paper_level,
	IFNULL(p.paper_id,'') AS remote_paper_id` +
		from + ` ORDER BY r.attempt_id DESC LIMIT ? OFFSET ?`
	listArgs := attemptAdminListCountArgs(joinArgs, whereArgs)
	listArgs = append(listArgs, size, offset)

	var rows []bo.AttemptAdminListRow
	if err := g.DB().Ctx(ctx).Raw(listSQL, listArgs...).Scan(&rows); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// AttemptAdminDetail 按 id 加载会话、学员、试卷及答题明细（含客观题是否选对）。
func (s *sAttempt) AttemptAdminDetail(ctx context.Context, attemptID int64) (*bo.AttemptAdminDetailView, error) {
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
			ok := examutil.ObjectiveAnswerCorrect(correctByQ[q.Id], payload.OptionID)
			row.ObjectiveCorrect = boolPtr(ok)
		}
		out.Answers = append(out.Answers, row)
	}
	return out, nil
}

// hasAnySubjectiveAwarded 该会话下是否已有主观题写入人工分（每会话仅允许保存一次）。
func hasAnySubjectiveAwarded(ctx context.Context, attemptID int64, paperID int64) (bool, error) {
	if attemptID <= 0 || paperID <= 0 {
		return false, nil
	}
	var cnt int
	if err := g.DB().Ctx(ctx).Raw(
		`SELECT COUNT(1) AS c FROM exam_attempt_answer eaa
INNER JOIN exam_question eq ON eq.id = eaa.exam_question_id
  AND eq.exam_paper_id = ? AND eq.is_subjective = 1 AND eq.is_example = 0
  AND eq.delete_flag = ?
WHERE eaa.attempt_id = ? AND eaa.delete_flag = ? AND eaa.awarded_score IS NOT NULL`,
		paperID, consts.DeleteFlagNotDeleted, attemptID, consts.DeleteFlagNotDeleted,
	).Scan(&cnt); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// AttemptAdminSaveSubjectiveScores 写入主观题人工分并汇总 subjective_score、total_score。每会话仅允许首次成功保存。
func (s *sAttempt) AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []bo.SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error) {
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
	graded, err := hasAnySubjectiveAwarded(ctx, att.Id, att.ExamPaperId)
	if err != nil {
		return 0, 0, err
	}
	if graded {
		return 0, 0, gerror.NewCode(consts.CodeExamSubjectiveAlreadyGraded)
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
					Creator:        updaterAdmin,
					Updater:        updaterAdmin,
					DeleteFlag:     consts.DeleteFlagNotDeleted,
					CreateTime:     gtime.Now(),
					UpdateTime:     gtime.Now(),
				}); err != nil {
					return err
				}
			} else {
				if _, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Where("id", row.Id).Update(examdo.ExamAttemptAnswer{
					AwardedScore: &scoreVal,
					Updater:      updaterAdmin,
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
			Updater:         updaterAdmin,
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
	var afterAttempt examentity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).Where("id", attemptID).Scan(&afterAttempt); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.ExamAttempt.Table(), attemptID, &att, &afterAttempt)
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
