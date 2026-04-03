package exam

import (
	"context"

	"exam/api/admin/exam/v1"
	"exam/internal/service/exam"
	"exam/internal/util"
)

func (c *ControllerV1) AttemptList(ctx context.Context, req *v1.AttemptListReq) (res *v1.AttemptListRes, err error) {
	rows, total, err := exam.AttemptAdminList(ctx, req.Page, req.Size, req.Level, req.ExaminationPaperId, req.ExamBatchId, req.Status, req.Username)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.AttemptListItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, &v1.AttemptListItem{
			Id:                 r.Id,
			MemberId:           r.MemberId,
			Username:           r.Username,
			Nickname:           r.Nickname,
			ExaminationPaperId: r.ExaminationPaperId,
			ExamBatchId:        r.ExamBatchId,
			MockLevelId:        r.MockLevelId,
			PaperTitle:         r.PaperTitle,
			PaperLevel:         r.PaperLevel,
			RemotePaperId:      r.RemotePaperId,
			Status:             r.Status,
			ObjectiveScore:     r.ObjectiveScore,
			SubjectiveScore:    r.SubjectiveScore,
			TotalScore:         r.TotalScore,
			HasSubjective:      r.HasSubjective,
			StartedAt:          util.ToRFC3339UTC(r.StartedAt),
			SubmittedAt:        util.ToRFC3339UTC(r.SubmittedAt),
			EndedAt:            util.ToRFC3339UTC(r.EndedAt),
			CreateTime:         util.ToRFC3339UTC(r.CreateTime),
		})
	}
	return &v1.AttemptListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) AttemptDetail(ctx context.Context, req *v1.AttemptDetailReq) (res *v1.AttemptDetailRes, err error) {
	d, err := exam.AttemptAdminDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	a := d.Attempt
	out := &v1.AttemptDetailRes{
		Attempt: v1.AttemptDetailAttempt{
			Id:                 a.Id,
			MemberId:           a.MemberId,
			ExaminationPaperId: a.MockExaminationPaperId,
			Status:             a.Status,
			DurationSeconds:    a.DurationSeconds,
			ObjectiveScore:     a.ObjectiveScore,
			SubjectiveScore:    a.SubjectiveScore,
			TotalScore:         a.TotalScore,
			HasSubjective:      a.HasSubjective,
			StartedAt:          util.ToRFC3339UTC(a.StartedAt),
			DeadlineAt:         util.ToRFC3339UTC(a.DeadlineAt),
			SubmittedAt:        util.ToRFC3339UTC(a.SubmittedAt),
			EndedAt:            util.ToRFC3339UTC(a.EndedAt),
			CreateTime:         util.ToRFC3339UTC(a.CreateTime),
		},
		User: v1.AttemptDetailUser{
			Id:       d.User.Id,
			Username: d.User.Username,
			Nickname: d.User.Nickname,
		},
		Paper: v1.AttemptDetailPaper{
			Id:      d.Paper.MockExaminationPaperId,
			Level:   d.Paper.Level,
			PaperId: d.Paper.PaperId,
			Title:   d.Paper.Title,
		},
		Answers: make([]v1.AttemptDetailAnswer, 0, len(d.Answers)),
	}
	for _, row := range d.Answers {
		q := row.Question
		var awarded *float64
		if q.IsSubjective != 0 && q.IsExample == 0 {
			v := row.Answer.AwardedScore
			awarded = &v
		}
		out.Answers = append(out.Answers, v1.AttemptDetailAnswer{
			QuestionId:       row.Answer.ExamQuestionId,
			QuestionNo:       q.QuestionNo,
			StemText:         q.StemText,
			IsExample:        q.IsExample,
			IsSubjective:     q.IsSubjective,
			Score:            q.Score,
			AnswerJson:       row.Answer.AnswerJson,
			AwardedScore:     awarded,
			ObjectiveCorrect: row.ObjectiveCorrect,
		})
	}
	return out, nil
}

func (c *ControllerV1) AttemptSubjectiveScores(ctx context.Context, req *v1.AttemptSubjectiveScoresReq) (res *v1.AttemptSubjectiveScoresRes, err error) {
	items := make([]exam.SubjectiveScoreItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, exam.SubjectiveScoreItem{QuestionID: it.QuestionId, Score: it.Score})
	}
	subSum, total, err := exam.AttemptAdminSaveSubjectiveScores(ctx, req.Id, items)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptSubjectiveScoresRes{SubjectiveScore: subSum, TotalScore: total}, nil
}
