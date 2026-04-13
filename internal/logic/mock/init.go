package mock

import svc "exam/internal/service/mock"

type sMock struct{}

func init() {
	svc.RegisterMock(New())
}

func New() *sMock {
	return &sMock{}
}
