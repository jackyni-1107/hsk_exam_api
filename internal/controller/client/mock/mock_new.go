package mock

import "exam/api/client/mock"

type ControllerV1 struct{}

func NewV1() mock.IMock {
	return &ControllerV1{}
}
