package auth

import (
	"exam/api/admin/auth"
)

type ControllerV1 struct{}

func NewV1() auth.IAuth {
	return &ControllerV1{}
}
