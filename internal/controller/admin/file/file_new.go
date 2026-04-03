package file

import "exam/api/admin/file"

type ControllerV1 struct{}

func NewV1() file.IFile {
	return &ControllerV1{}
}
