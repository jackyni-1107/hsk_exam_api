package member

import membersvc "exam/internal/service/member"

func init() {
	membersvc.RegisterMember(new(sMember))
}
