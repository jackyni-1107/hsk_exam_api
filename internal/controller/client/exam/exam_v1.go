package exam

import (
	"context"
	v1 "exam/api/client/exam/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	exambo "exam/internal/model/bo/exam"
	attemptsvc "exam/internal/service/attempt"
	papersvc "exam/internal/service/paper"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) PaperForExam(ctx context.Context, req *v1.PaperForExamReq) (res *v1.PaperForExamRes, err error) {
	data := middleware.GetCtxData(ctx)
	if data == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}

	d, segments, mock, err := papersvc.Paper().PaperBootstrapForExam(ctx, req.PaperId)
	if err != nil {
		return nil, err
	}

	prepareSegments := make([]v1.PaperForExamSegment, 0, len(segments))
	for _, seg := range segments {
		partItems := make([]v1.PaperForExamSegmentPart, 0, len(seg.Parts))
		for _, part := range seg.Parts {
			partItems = append(partItems, v1.PaperForExamSegmentPart{
				PartCode:                part.PartCode,
				PartName:                part.PartName,
				PartNameTrans:           part.PartNameTrans,
				PartRate:                part.PartRate,
				PartScore:               part.PartScore,
				QuestionCount:           part.QuestionCount,
				ObjectiveQuestionCount:  part.ObjectiveQuestionCount,
				SubjectiveQuestionCount: part.SubjectiveQuestionCount,
				PartAnswerTime:          part.PartAnswerTime,
				ScoreTotal:              part.ScoreTotal,
				CorrectCount:            part.CorrectCount,
				CorrectRate:             part.CorrectRate,
				Practiced:               part.Practiced,
				QuestionType:            part.QuestionType,
			})
		}
		prepareSegments = append(prepareSegments, v1.PaperForExamSegment{
			SegmentCode:   seg.SegmentCode,
			TotalScore:    seg.TotalScore,
			QuestionCount: seg.QuestionCount,
			Duration:      seg.Duration,
			Seq:           seg.Seq,
			SegmentDesc:   seg.SegmentDesc,
			Parts:         partItems,
		})
	}
	playURL, _, _ := papersvc.Paper().IssuePaperHlsPlay(ctx, data.UserId, d.Paper.Id)
	res = &v1.PaperForExamRes{
		Id:                     d.Paper.Id,
		MockExaminationPaperId: mock.Id,
		Level:                  d.Paper.Level,
		PaperId:              d.Paper.PaperId,
		Title:                d.Paper.Title,
		SourceBaseUrl:        d.Paper.SourceBaseUrl,
		AudioUrl:             playURL,
		DurationSeconds:      d.Paper.DurationSeconds,
		ListenReviewDuration: d.Paper.ListenReviewDuration,
		Prepare: v1.PaperForExamPrepare{
			Instruction: d.Paper.PrepareInstruction,
			AudioFile:   d.Paper.PrepareAudioFile,
			Title:       d.Paper.PrepareTitle,
			Name:        mock.Name,
			Segments:    prepareSegments,
		},
		Items: make([]v1.PaperForExamItemInit, 0, len(d.Sections)),
	}
	for _, sec := range d.Sections {
		item := v1.PaperForExamItemInit{
			Id:            sec.Id,
			SortOrder:     sec.SortOrder,
			TopicTitle:    sec.TopicTitle,
			TopicSubtitle: sec.TopicSubtitle,
			TopicType:     sec.TopicType,
			PartCode:      sec.PartCode,
			SegmentCode:   sec.SegmentCode,
			TopicItems:    sec.TopicItemsFile,
		}
		res.Items = append(res.Items, item)
	}
	return res, nil
}

func (c *ControllerV1) PaperSectionForExam(ctx context.Context, req *v1.PaperSectionForExamReq) (res *exambo.SectionTopic, err error) {
	return papersvc.Paper().PaperSectionTopicForExam(ctx, req.PaperId, req.SectionId)
}

func (c *ControllerV1) AttemptCreateByBatch(ctx context.Context, req *v1.AttemptCreateByBatchReq) (res *v1.AttemptCreateRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	//_, err = exam.Exam().ExamBatchMemberDetail(ctx, req.BatchId, ctxData.UserId, req.MockExaminationPaperId)
	//if err != nil {
	//	return nil, err
	//}
	var id int64
	id, err = attemptsvc.Attempt().CreateAttemptForBatch(ctx, ctxData.UserId, req.BatchId, req.ExamPaperId)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptCreateRes{AttemptId: id}, nil
}

func (c *ControllerV1) AttemptStart(ctx context.Context, req *v1.AttemptStartReq) (res *v1.AttemptStartRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	err = attemptsvc.Attempt().StartAttempt(ctx, ctxData.UserId, req.Id, req.DurationSeconds)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptStartRes{}, nil
}

func (c *ControllerV1) AttemptGet(ctx context.Context, req *v1.AttemptGetReq) (res *v1.AttemptGetRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	v, err := attemptsvc.Attempt().GetAttempt(ctx, ctxData.UserId, req.Id)
	if err != nil {
		return nil, err
	}
	a := v.Attempt
	out := &v1.AttemptGetRes{
		Id:                 a.Id,
		ExamPaperId:        a.ExamPaperId,
		ExaminationPaperId: a.MockExaminationPaperId,
		Status:             a.Status,
		//DurationSeconds:    a.DurationSeconds,
		//ObjectiveScore:     a.ObjectiveScore,
		//SubjectiveScore:    a.SubjectiveScore,
		//TotalScore:         a.TotalScore,
		//HasSubjective:      a.HasSubjective,
		ServerTime: v.ServerTime,
		//DeadlineReached:    v.DeadlineReached,
		SegmentCode:      v.SegmentCode,
		RemainingSeconds: v.RemainingSeconds,
	}
	out.StartedAt = utility.ToRFC3339UTC(a.StartedAt)
	out.DeadlineAt = utility.ToRFC3339UTC(a.DeadlineAt)
	//out.SubmittedAt = utility.ToRFC3339UTC(a.SubmittedAt)
	//out.EndedAt = utility.ToRFC3339UTC(a.EndedAt)
	return out, nil
}

func (c *ControllerV1) AttemptAnswersGet(ctx context.Context, req *v1.AttemptAnswersGetReq) (res *v1.AttemptAnswersGetRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	rows, err := attemptsvc.Attempt().GetAttemptAnswers(ctx, ctxData.UserId, req.Id)
	if err != nil {
		return nil, err
	}
	items := make([]v1.AttemptAnswerItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.AttemptAnswerItem{
			QuestionId: r.QuestionID,
			OptionID:   r.OptionID,
			Text:       r.Text,
		})
	}
	return &v1.AttemptAnswersGetRes{Items: items}, nil
}

func (c *ControllerV1) AttemptSaveAnswers(ctx context.Context, req *v1.AttemptSaveAnswersReq) (res *v1.AttemptSaveAnswersRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	items := make([]bo.SaveAnswerItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, bo.SaveAnswerItem{
			QuestionID: it.QuestionId,
			OptionID:   it.OptionID,
			Text:       it.Text,
		})
	}
	err = attemptsvc.Attempt().SaveAnswers(ctx, ctxData.UserId, req.Id, req.SegmentCode, items)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptSaveAnswersRes{}, nil
}

func (c *ControllerV1) AttemptSubmit(ctx context.Context, req *v1.AttemptSubmitReq) (res *v1.AttemptSubmitRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	err = attemptsvc.Attempt().SubmitAttempt(ctx, ctxData.UserId, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptSubmitRes{}, nil
}

func (c *ControllerV1) AttemptRandomAnswers(ctx context.Context, req *v1.AttemptRandomAnswersReq) (res *v1.AttemptRandomAnswersRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	drafts, err := papersvc.Paper().RandomFillAnswersForTest(ctx, ctxData.UserId, req.PaperId, req.AttemptId)
	if err != nil {
		return nil, err
	}
	items := make([]v1.AttemptAnswerItem, 0, len(drafts))
	for _, d := range drafts {
		items = append(items, v1.AttemptAnswerItem{
			QuestionId: d.QuestionID,
			OptionID:   d.OptionID,
			Text:       d.Text,
		})
	}
	return &v1.AttemptRandomAnswersRes{
		GeneratedCount: len(items),
		SubmitJSON: v1.AttemptSaveAnswersBody{
			Items: items,
		},
	}, nil
}

//
//func (c *ControllerV1) AudioHlsPlayIssue(ctx context.Context, req *v1.AudioHlsPlayIssueReq) (res *v1.AudioHlsPlayIssueRes, err error) {
//	ctxData := middleware.GetCtxData(ctx)
//	if ctxData == nil {
//		return nil, gerror.NewCode(consts.CodeTokenRequired)
//	}
//	playURL, exp, err := exam.Exam().IssueAudioHlsPlay(ctx, ctxData.UserId, req.Id, req.QuestionId)
//	if err != nil {
//		return nil, err
//	}
//	return &v1.AudioHlsPlayIssueRes{PlayUrl: playURL, ExpiresAt: exp}, nil
//}
