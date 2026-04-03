package task

import "exam/api/admin/task"

type ControllerV1 struct{}

func NewV1() task.ITask {
	return &ControllerV1{}
}
