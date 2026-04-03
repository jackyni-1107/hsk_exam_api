package user

import (
	"exam/api/admin/user"
)

type ControllerV1 struct{}

func NewV1() user.IUser {
	return &ControllerV1{}
}
