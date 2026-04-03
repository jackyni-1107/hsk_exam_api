package exam

import "exam/api/client/exam"

type ControllerV1 struct{}

func NewV1() exam.IExam {
	return &ControllerV1{}
}
