package member

import (
	"exam/api/admin/member"
)

type ControllerV1 struct{}

func NewV1() member.IMember {
	return &ControllerV1{}
}
