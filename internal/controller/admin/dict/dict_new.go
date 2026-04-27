package dict

import "exam/api/admin/dict"

type ControllerV1 struct{}

func NewV1() dict.IDict {
	return &ControllerV1{}
}
