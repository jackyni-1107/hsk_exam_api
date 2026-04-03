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
	"exam/internal/logic/clientexam"
	"exam/internal/logic/examresult"
	"exam/internal/model/do"
	"exam/internal/model/entity"
)

// AttemptAdminListRow 管理端列表行（与 Raw 列别名一致，供 Scan）。
type AttemptAdminListRow struct {
	Id                 int64       `json:"id"`
	MemberId           int64       `json:"member_id"`
	ExaminationPaperId int64       `json:"examination_paper_id"`
	Status             int         `json:"status"`
	ObjectiveScore     float64     `json:"objective_score"`
	SubjectiveScore    float64     `json:"subjective_score"`
	TotalScore         float64     `json:"total_score"`
	HasSubjective      int         `json:"has_subjective"`
	StartedAt          *gtime.Time `json:"started_at"`
	SubmittedAt        *gtime.Time `json:"submitted_at"`
	EndedAt            *gtime.Time `json:"ended_at"`
	CreateTime         *gtime.Time `json:"create_time"`
	Username           string      `json:"username"`
	Nickname           string      `json:"nickname"`
	PaperTitle         string      `json:"paper_title"`
	PaperLevel         string      `json:"paper_level"`
	RemotePaperId      string      `json:"remote_paper_id"`
}

// AttemptAdminDetailView 管理端会话详情（服务层聚合）。
type AttemptAdminDetailView struct {
	Attempt entity.ExamAttempt
	User    entity.ClientUser
	Paper   entity.ExamPaper
	Answers []AttemptAdminAnswerRow
}

// AttemptAdminAnswerRow 单题答题展示行。
type AttemptAdminAnswerRow struct {
	Answer           entity.ExamAttemptAnswer
	Question         entity.ExamQuestion
	ObjectiveCorrect *bool
}

// AttemptAdminList 分页查询答题会话（联表学员、试卷）。
func AttemptAdminList(ctx context.Context, page, size int, level string, examinationPaperId int64, status int, username string) ([]AttemptAdminListRow, int, error) {
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
WHERE ` + w
	joinArgs := []interface{}{consts.DeleteFlagNotDeleted, consts.DeleteFlagNotDeleted, consts.DeleteFlagNotDeleted}

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
	listSQL := `SELECT r.attempt_id AS id, r.member_id, IFNULL(p.mock_examination_paper_id,0) AS examination_paper_id, r.status,
r.objective_score, r.subjective_score, r.total_score, r.has_subjective,
r.started_at, r.submitted_at, r.ended_at, r.create_time,
IFNULL(u.username,'') AS username, IFNULL(u.nickname,'') AS nickname,
IFNULL(p.title,'') AS paper_title, IFNULL(p.level,'') AS paper_level, IFNULL(p.paper_id,'') AS remote_paper_id` +
		from + ` ORDER BY r.attempt_id DESC LIMIT ? OFFSET ?`
	listArgs := append(append([]interface{}{}, joinArgs...), args...)
	listArgs = append(listArgs, size, offset)

	var rows []AttemptAdminListRow
	if err := g.DB().Ctx(ctx).Raw(listSQL, listArgs...).Scan(&rows); err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// AttemptAdminDetail 按 id 加载会话、学员、试卷及答题明细（含客观题是否选对）。
func AttemptAdminDetail(ctx context.Context, attemptID int64) (*AttemptAdminDetailView, error) {
	if attemptID <= 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	var att entity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att); err != nil {
		return nil, err
	}
	if att.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	var user entity.ClientUser
	_ = dao.SysMember.Ctx(ctx).
		Where("id", att.MemberId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&user)
	var paper entity.ExamPaper
	_ = dao.ExamPaper.Ctx(ctx).
		Where("id", att.ExamPaperId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&paper)

	var ansRows []entity.ExamAttemptAnswer
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
	qByID := make(map[int64]entity.ExamQuestion)
	if len(qIDs) > 0 {
		var qs []entity.ExamQuestion
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
	correctByQ := loadCorrectOptionIDsByQuestion(ctx, qIDs)

	out := &AttemptAdminDetailView{
		Attempt: att,
		User:    user,
		Paper:   paper,
		Answers: make([]AttemptAdminAnswerRow, 0, len(ansRows)),
	}
	for _, ar := range ansRows {
		q := qByID[ar.ExamQuestionId]
		row := AttemptAdminAnswerRow{Answer: ar, Question: q}
		if q.Id != 0 && q.IsExample == 0 && q.IsSubjective == 0 {
			payload := clientexam.ParseAnswerPayload(ar.AnswerJson)
			ok := clientexam.ObjectiveAnswerCorrect(correctByQ[q.Id], payload.SelectedOptionIDs)
			row.ObjectiveCorrect = boolPtr(ok)
		}
		out.Answers = append(out.Answers, row)
	}
	return out, nil
}

// SubjectiveScoreItem 管理端提交的主观题得分项。
type SubjectiveScoreItem struct {
	QuestionID int64
	Score      float64
}

// AttemptAdminSaveSubjectiveScores 写入主观题人工分并汇总 subjective_score、total_score（允许部分题目已评）。
func AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error) {
	if attemptID <= 0 {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	if len(items) == 0 {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
	}
	var att entity.ExamAttempt
	if err := dao.ExamAttempt.Ctx(ctx).
		Where("id", attemptID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&att); err != nil {
		return 0, 0, err
	}
	if att.Id == 0 {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
	}
	if att.Status != consts.ExamAttemptEnded {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_ended")
	}
	if att.HasSubjective != 1 {
		return 0, 0, gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_no_subjective")
	}

	byQ := make(map[int64]SubjectiveScoreItem, len(items))
	for _, it := range items {
		byQ[it.QuestionID] = it
	}
	uniq := make([]SubjectiveScoreItem, 0, len(byQ))
	for _, it := range byQ {
		uniq = append(uniq, it)
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var attTx entity.ExamAttempt
		if err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Scan(&attTx); err != nil {
			return err
		}
		if attTx.Id == 0 {
			return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_found")
		}
		if attTx.Status != consts.ExamAttemptEnded {
			return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_not_ended")
		}
		if attTx.HasSubjective != 1 {
			return gerror.NewCode(consts.CodeInvalidParams, "err.exam_attempt_no_subjective")
		}
		paperID := attTx.ExamPaperId

		for _, it := range uniq {
			var q entity.ExamQuestion
			if err := tx.Model(dao.ExamQuestion.Table()).Ctx(ctx).
				Where("id", it.QuestionID).
				Where("exam_paper_id", paperID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&q); err != nil {
				return err
			}
			if q.Id == 0 || q.IsSubjective != 1 || q.IsExample != 0 {
				return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
			}
			if it.Score < 0 || it.Score > q.Score {
				return gerror.NewCode(consts.CodeInvalidParams, "err.invalid_params")
			}

			var row entity.ExamAttemptAnswer
			_ = tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).
				Where("attempt_id", attemptID).
				Where("exam_question_id", it.QuestionID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&row)
			score := it.Score
			if row.Id == 0 {
				if _, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Insert(do.ExamAttemptAnswer{
					AttemptId:      attemptID,
					ExamQuestionId: it.QuestionID,
					AnswerJson:     "{}",
					AwardedScore:   score,
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
				if _, err := tx.Model(dao.ExamAttemptAnswer.Table()).Ctx(ctx).Where("id", row.Id).Update(do.ExamAttemptAnswer{
					AwardedScore: score,
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
		if _, err := tx.Model(dao.ExamAttempt.Table()).Ctx(ctx).Where("id", attemptID).Update(do.ExamAttempt{
			SubjectiveScore: sum,
			TotalScore:      tot,
			Updater:         "admin",
			UpdateTime:      gtime.Now(),
		}); err != nil {
			return err
		}
		if err := examresult.UpsertFromAttemptTx(ctx, tx, attemptID); err != nil {
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
	var opts []entity.ExamOption
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
