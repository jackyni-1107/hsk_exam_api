package batch

import (
	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
	"exam/internal/model/bo"
)

func normalizeExamBatchPolicy(p bo.ExamBatchPolicyInput) (bo.ExamBatchPolicyInput, error) {
	if p.BatchKind != consts.ExamBatchKindFormal && p.BatchKind != consts.ExamBatchKindPractice {
		return p, gerror.NewCode(consts.CodeInvalidParams)
	}
	if p.AllowMultipleAttempts != 0 {
		p.AllowMultipleAttempts = 1
	}
	if p.SkipScoring != 0 {
		p.SkipScoring = 1
	}
	if p.AutoSubmitOnDeadline != 0 {
		p.AutoSubmitOnDeadline = 1
	}
	if p.AllowMultipleAttempts == 0 {
		p.MaxAttemptsPerMember = 0
	}
	if p.MaxAttemptsPerMember < 0 {
		return p, gerror.NewCode(consts.CodeInvalidParams)
	}
	return p, nil
}
