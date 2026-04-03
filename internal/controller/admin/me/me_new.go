package me

import (
	"exam/api/admin/me"
)

type ControllerV1 struct{}

func NewV1() me.IMe {
	return &ControllerV1{}
}
