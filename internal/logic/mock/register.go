package mock

import svc "exam/internal/service/mock"

func init() {
	svc.RegisterMock(new(sMock))
}
