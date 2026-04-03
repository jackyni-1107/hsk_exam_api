package mock

import "exam/api/admin/mock"

type ControllerV1 struct{}

func NewV1() mock.IMock {
	return &ControllerV1{}
}
