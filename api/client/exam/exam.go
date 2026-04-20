package exam

import (
	"context"

	v1 "exam/api/client/exam/v1"
	exambo "exam/internal/model/bo/exam"
)

type IExam interface {
	PaperForExam(ctx context.Context, req *v1.PaperForExamReq) (res *v1.PaperForExamRes, err error)
	// PaperSectionForExam 响应体与 topic JSON（如 pt.json）根结构一致。
	PaperSectionForExam(ctx context.Context, req *v1.PaperSectionForExamReq) (res *exambo.SectionTopic, err error)
	//AttemptCreate(ctx context.Context, req *v1.AttemptCreateReq) (res *v1.AttemptCreateRes, err error)
	AttemptCreateByBatch(ctx context.Context, req *v1.AttemptCreateByBatchReq) (res *v1.AttemptCreateRes, err error)
	AttemptStart(ctx context.Context, req *v1.AttemptStartReq) (res *v1.AttemptStartRes, err error)
	AttemptGet(ctx context.Context, req *v1.AttemptGetReq) (res *v1.AttemptGetRes, err error)
	AttemptAnswersGet(ctx context.Context, req *v1.AttemptAnswersGetReq) (res *v1.AttemptAnswersGetRes, err error)
	AttemptSaveAnswers(ctx context.Context, req *v1.AttemptSaveAnswersReq) (res *v1.AttemptSaveAnswersRes, err error)
	AttemptSubmit(ctx context.Context, req *v1.AttemptSubmitReq) (res *v1.AttemptSubmitRes, err error)
	AttemptRandomAnswers(ctx context.Context, req *v1.AttemptRandomAnswersReq) (res *v1.AttemptRandomAnswersRes, err error)
	//AudioHlsPlayIssue(ctx context.Context, req *v1.AudioHlsPlayIssueReq) (res *v1.AudioHlsPlayIssueRes, err error)
}
