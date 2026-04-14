package member

import "exam/internal/service/member"

type sMember struct{}

func init() {
	member.RegisterMember(New())
}

func New() *sMember {
	return &sMember{}
}
