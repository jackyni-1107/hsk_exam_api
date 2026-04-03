package exam

import svc "exam/internal/service/exam"

func init() {
	svc.RegisterExam(new(sExam))
}
