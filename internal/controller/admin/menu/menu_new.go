package menu

import (
	"exam/api/admin/menu"
)

type ControllerV1 struct{}

func NewV1() menu.IMenu {
	return &ControllerV1{}
}
