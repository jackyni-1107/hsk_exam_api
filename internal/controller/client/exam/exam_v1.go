package exam

import (
	"context"
	"encoding/json"
	v1 "exam/api/client/exam/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	"exam/internal/service/exam"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) PaperForExam(ctx context.Context, req *v1.PaperForExamReq) (res *v1.PaperForExamRes, err error) {
	d, err := exam.Exam().PaperDetailForExamInit(ctx, req.PaperId)
	if err != nil {
		return nil, err
	}
	res = &v1.PaperForExamRes{
		Id:                   d.Paper.Id,
		Level:                d.Paper.Level,
		PaperId:              d.Paper.PaperId,
		Title:                d.Paper.Title,
		SourceBaseUrl:        d.Paper.SourceBaseUrl,
		ListeningAudioPrefix: d.Paper.AudioHlsPrefix,
		DurationSeconds:      d.Paper.DurationSeconds,
		Prepare: v1.PaperForExamPrepare{
			Instruction: d.Paper.PrepareInstruction,
			AudioFile:   d.Paper.PrepareAudioFile,
			Title:       d.Paper.PrepareTitle,
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
			//TopicJson:     sec.TopicJson,
			//Blocks:        make([]v1.PaperForExamBlockInit, 0, len(sec.Blocks)),
		}
		//for _, blk := range sec.Blocks {
		//	item.Blocks = append(item.Blocks, v1.PaperForExamBlockInit{
		//		Id:                      blk.Id,
		//		BlockOrder:              blk.BlockOrder,
		//		GroupIndex:              blk.GroupIndex,
		//		QuestionDescriptionJson: blk.QuestionDescriptionJson,
		//		QuestionCount:           blk.QuestionCount,
		//	})
		//}
		res.Items = append(res.Items, item)
	}
	return res, nil
}

func (c *ControllerV1) PaperSectionForExam(ctx context.Context, req *v1.PaperSectionForExamReq) (res map[string]interface{}, err error) {
	return exam.Exam().PaperSectionTopicForExam(ctx, req.PaperId, req.SectionId)
}

//func (c *ControllerV1) AttemptCreate(ctx context.Context, req *v1.AttemptCreateReq) (res *v1.AttemptCreateRes, err error) {
//	ctxData := middleware.GetCtxData(ctx)
//	if ctxData == nil {
//		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
//	}
//	id, err := exam.Exam().CreateAttempt(ctx, ctxData.UserId, req.PaperId)
//	if err != nil {
//		return nil, err
//	}
//	return &v1.AttemptCreateRes{AttemptId: id}, nil
//}

func (c *ControllerV1) AttemptCreateByBatch(ctx context.Context, req *v1.AttemptCreateByBatchReq) (res *v1.AttemptCreateRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	//
	batch, err := exam.Exam().ExamBatchMemberDetail(ctx, req.BatchId, ctxData.UserId)
	if err != nil {
		return nil, err
	}
	id, err := exam.Exam().CreateAttemptForBatch(ctx, ctxData.UserId, batch.BatchId, batch.MemberId)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptCreateRes{AttemptId: id}, nil
}

func (c *ControllerV1) AttemptStart(ctx context.Context, req *v1.AttemptStartReq) (res *v1.AttemptStartRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	err = exam.Exam().StartAttempt(ctx, ctxData.UserId, req.Id, req.DurationSeconds)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptStartRes{}, nil
}

func (c *ControllerV1) AttemptGet(ctx context.Context, req *v1.AttemptGetReq) (res *v1.AttemptGetRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	v, err := exam.Exam().GetAttempt(ctx, ctxData.UserId, req.Id)
	if err != nil {
		return nil, err
	}
	a := v.Attempt
	out := &v1.AttemptGetRes{
		Id:                 a.Id,
		ExaminationPaperId: a.MockExaminationPaperId,
		Status:             a.Status,
		DurationSeconds:    a.DurationSeconds,
		ObjectiveScore:     a.ObjectiveScore,
		SubjectiveScore:    a.SubjectiveScore,
		TotalScore:         a.TotalScore,
		HasSubjective:      a.HasSubjective,
		ServerTime:         v.ServerTime,
		DeadlineReached:    v.DeadlineReached,
	}
	out.StartedAt = util.ToRFC3339UTCPtr(a.StartedAt)
	out.DeadlineAt = util.ToRFC3339UTCPtr(a.DeadlineAt)
	out.SubmittedAt = util.ToRFC3339UTCPtr(a.SubmittedAt)
	out.EndedAt = util.ToRFC3339UTCPtr(a.EndedAt)
	return out, nil
}

func (c *ControllerV1) AttemptSaveAnswers(ctx context.Context, req *v1.AttemptSaveAnswersReq) (res *v1.AttemptSaveAnswersRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	items := make([]bo.SaveAnswerItem, 0, len(req.Items))
	for _, it := range req.Items {
		payload := make(map[string]interface{}, 1)
		switch v := it.Answer.(type) {
		case float64:
			if v > 0 {
				id := int64(v)
				payload["option_id"] = id
			}
		case json.Number:
			if id, err := v.Int64(); err == nil && id > 0 {
				payload["option_id"] = id
			}
		case int:
			if v > 0 {
				id := int64(v)
				payload["option_id"] = id
			}
		case int64:
			if v > 0 {
				id := v
				payload["option_id"] = id
			}
		case string:
			payload["text"] = v
		}
		raw, _ := json.Marshal(payload)
		items = append(items, bo.SaveAnswerItem{
			QuestionID: it.QuestionId,
			AnswerJSON: string(raw),
		})
	}
	err = exam.Exam().SaveAnswers(ctx, ctxData.UserId, req.Id, items)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptSaveAnswersRes{}, nil
}

func (c *ControllerV1) AttemptSubmit(ctx context.Context, req *v1.AttemptSubmitReq) (res *v1.AttemptSubmitRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	err = exam.Exam().SubmitAttempt(ctx, ctxData.UserId, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.AttemptSubmitRes{}, nil
}

func (c *ControllerV1) AttemptRandomAnswers(ctx context.Context, req *v1.AttemptRandomAnswersReq) (res *v1.AttemptRandomAnswersRes, err error) {
	ctxData := middleware.GetCtxData(ctx)
	if ctxData == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
	}
	drafts, err := exam.Exam().RandomFillAnswersForTest(ctx, ctxData.UserId, req.PaperId, req.AttemptId)
	if err != nil {
		return nil, err
	}
	items := make([]v1.AttemptAnswerItem, 0, len(drafts))
	for _, d := range drafts {
		items = append(items, v1.AttemptAnswerItem{
			QuestionId: d.QuestionID,
			Answer:     d.Answer,
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
//		return nil, gerror.NewCode(consts.CodeTokenRequired, "")
//	}
//	playURL, exp, err := exam.Exam().IssueAudioHlsPlay(ctx, ctxData.UserId, req.Id, req.QuestionId)
//	if err != nil {
//		return nil, err
//	}
//	return &v1.AudioHlsPlayIssueRes{PlayUrl: playURL, ExpiresAt: exp}, nil
//}
