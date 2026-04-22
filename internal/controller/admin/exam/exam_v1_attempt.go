package exam

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	v1 "exam/api/admin/exam/v1"
	"exam/internal/consts"
	mockdao "exam/internal/dao/mock"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
	mockentity "exam/internal/model/entity/mock"
	attemptsvc "exam/internal/service/attempt"
	mocksvc "exam/internal/service/mock"
	"exam/internal/utility"
)

func (c *ControllerV1) AttemptList(ctx context.Context, req *v1.AttemptListReq) (res *v1.AttemptListRes, err error) {
	rows, total, err := attemptsvc.Attempt().AttemptAdminList(ctx, req.Page, req.Size, req.Level, req.ExaminationPaperId, req.ExamBatchId, req.Status, req.Username)
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
			SubjectiveGraded:   r.SubjectiveGraded,
			StartedAt:          utility.ToRFC3339UTC(r.StartedAt),
			SubmittedAt:        utility.ToRFC3339UTC(r.SubmittedAt),
			EndedAt:            utility.ToRFC3339UTC(r.EndedAt),
			CreateTime:         utility.ToRFC3339UTC(r.CreateTime),
		})
	}
	return &v1.AttemptListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) AttemptDetail(ctx context.Context, req *v1.AttemptDetailReq) (res *v1.AttemptDetailRes, err error) {
	d, err := attemptsvc.Attempt().AttemptAdminDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	a := d.Attempt

	mockPaperName := ""
	if d.Paper.MockExaminationPaperId > 0 {
		paperDetail, _ := mocksvc.Mock().ExaminationPaperDetail(ctx, d.Paper.MockExaminationPaperId)
		if paperDetail != nil {
			mockPaperName = paperDetail.Name
		}
	}

	levelDisplay := strings.TrimSpace(d.Paper.Level)
	if a.MockLevelId > 0 {
		// MockLevels lookup kept via DAO: no service method for single level by ID yet
		var lv mockentity.MockLevels
		_ = mockdao.MockLevels.Ctx(ctx).
			Where(mockdao.MockLevels.Columns().Id, a.MockLevelId).
			Where(mockdao.MockLevels.Columns().DeleteFlag, consts.DeleteFlagNotDeleted).
			Scan(&lv)
		if strings.TrimSpace(lv.LevelName) != "" {
			levelDisplay = strings.TrimSpace(lv.LevelName)
		}
	}

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
			StartedAt:          utility.ToRFC3339UTC(a.StartedAt),
			DeadlineAt:         utility.ToRFC3339UTC(a.DeadlineAt),
			SubmittedAt:        utility.ToRFC3339UTC(a.SubmittedAt),
			EndedAt:            utility.ToRFC3339UTC(a.EndedAt),
			CreateTime:         utility.ToRFC3339UTC(a.CreateTime),
		},
		User: v1.AttemptDetailUser{
			Id:       d.User.Id,
			Username: d.User.Username,
			Nickname: d.User.Nickname,
		},
		Paper: v1.AttemptDetailPaper{
			Id:            d.Paper.MockExaminationPaperId,
			Name:          mockPaperName,
			Level:         levelDisplay,
			PaperId:       d.Paper.PaperId,
			Title:         d.Paper.Title,
			ExamPaperId:   d.Paper.Id,
			SourceBaseUrl: d.Paper.SourceBaseUrl,
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
		var sectionID int64
		var sectionTitle string
		if row.Section != nil {
			sectionID = row.Section.Id
			sectionTitle = row.Section.TopicTitle
		}
		opts := mapAttemptDetailOptions(row.Options)
		out.Answers = append(out.Answers, v1.AttemptDetailAnswer{
			QuestionId:       row.Answer.ExamQuestionId,
			QuestionNo:       q.QuestionNo,
			StemText:         q.StemText,
			ScreenTextJson:   q.ScreenTextJson,
			IsExample:        q.IsExample,
			IsSubjective:     q.IsSubjective,
			Score:            q.Score,
			AnswerJson:       row.Answer.AnswerJson,
			AwardedScore:     awarded,
			ObjectiveCorrect: row.ObjectiveCorrect,
			SectionId:        sectionID,
			SectionTitle:     sectionTitle,
			AnalysisText:     pickAnalysisText(q.AnalysisJson),
			Options:          opts,
		})
	}
	return out, nil
}

func mapAttemptDetailOptions(src []examentity.ExamOption) []v1.AttemptDetailOption {
	if len(src) == 0 {
		return nil
	}
	cp := append([]examentity.ExamOption(nil), src...)
	sort.Slice(cp, func(i, j int) bool {
		if cp[i].SortOrder != cp[j].SortOrder {
			return cp[i].SortOrder < cp[j].SortOrder
		}
		return cp[i].Id < cp[j].Id
	})
	out := make([]v1.AttemptDetailOption, 0, len(cp))
	for _, o := range cp {
		out = append(out, v1.AttemptDetailOption{
			Id:         o.Id,
			Flag:       o.Flag,
			SortOrder:  o.SortOrder,
			IsCorrect:  o.IsCorrect,
			OptionType: o.OptionType,
			Content:    o.Content,
		})
	}
	return out
}

func pickAnalysisText(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return raw
	}
	for _, k := range []string{"zh", "zh-CN", "cn", "en"} {
		if s := stringifyAnalysisValue(m[k]); s != "" {
			return s
		}
	}
	for _, v := range m {
		if s := stringifyAnalysisValue(v); s != "" {
			return s
		}
	}
	return ""
}

func stringifyAnalysisValue(v interface{}) string {
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x)
	case []interface{}:
		var parts []string
		for _, it := range x {
			if s := stringifyAnalysisValue(it); s != "" {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, " ")
	default:
		return ""
	}
}

func (c *ControllerV1) AttemptSubjectiveScores(ctx context.Context, req *v1.AttemptSubjectiveScoresReq) (res *v1.AttemptSubjectiveScoresRes, err error) {
	items := make([]bo.SubjectiveScoreItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, bo.SubjectiveScoreItem{QuestionID: it.QuestionId, Score: it.Score})
	}
	subSum, total, err := attemptsvc.Attempt().AttemptAdminSaveSubjectiveScores(ctx, req.Id, items)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptSubjectiveScoresRes{SubjectiveScore: subSum, TotalScore: total}, nil
}
