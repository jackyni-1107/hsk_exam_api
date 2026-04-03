package role

import (
	"exam/api/admin/role"
)

type ControllerV1 struct{}

func NewV1() role.IRole {
	return &ControllerV1{}
}
