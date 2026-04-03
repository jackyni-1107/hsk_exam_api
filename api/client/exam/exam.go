package exam

import (
	"context"

	v1 "exam/api/client/exam/v1"
)

type IExam interface {
	PaperForExam(ctx context.Context, req *v1.PaperForExamReq) (res *v1.PaperForExamRes, err error)
	// PaperSectionForExam 响应体与 topic JSON（如 pt.json）根结构一致，为 map 序列化结果。
	PaperSectionForExam(ctx context.Context, req *v1.PaperSectionForExamReq) (res map[string]interface{}, err error)
	AttemptCreate(ctx context.Context, req *v1.AttemptCreateReq) (res *v1.AttemptCreateRes, err error)
	AttemptStart(ctx context.Context, req *v1.AttemptStartReq) (res *v1.AttemptStartRes, err error)
	AttemptGet(ctx context.Context, req *v1.AttemptGetReq) (res *v1.AttemptGetRes, err error)
	AttemptSaveAnswers(ctx context.Context, req *v1.AttemptSaveAnswersReq) (res *v1.AttemptSaveAnswersRes, err error)
	AttemptSubmit(ctx context.Context, req *v1.AttemptSubmitReq) (res *v1.AttemptSubmitRes, err error)
	AttemptRandomAnswers(ctx context.Context, req *v1.AttemptRandomAnswersReq) (res *v1.AttemptRandomAnswersRes, err error)
	AudioHlsPlayIssue(ctx context.Context, req *v1.AudioHlsPlayIssueReq) (res *v1.AudioHlsPlayIssueRes, err error)
}
