package exam

import "exam/api/admin/exam"

type ControllerV1 struct{}

func NewV1() exam.IExam {
	return &ControllerV1{}
}
