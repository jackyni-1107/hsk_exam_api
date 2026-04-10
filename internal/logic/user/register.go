package user

import usersvc "exam/internal/service/user"

func init() {
	usersvc.RegisterUser(new(sUser))
}
