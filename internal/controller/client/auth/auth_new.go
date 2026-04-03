package auth

import (
	"exam/api/client/auth"
)

type ControllerV1 struct{}

func NewV1() auth.IAuth {
	return &ControllerV1{}
}
