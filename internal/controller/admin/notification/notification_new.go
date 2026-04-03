package notification

import "exam/api/admin/notification"

type ControllerV1 struct{}

func NewV1() notification.INotification {
	return &ControllerV1{}
}
