package me

import "exam/api/client/me"

type ControllerV1 struct{}

func NewV1() me.IMe {
	return &ControllerV1{}
}
