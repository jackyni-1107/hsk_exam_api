package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
)

type IUser interface {
	UserList(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error)
	UserCreate(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error)
	UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error)
	UserDelete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error)
	UserRoleAssign(ctx context.Context, req *v1.UserRoleAssignReq) (res *v1.UserRoleAssignRes, err error)
	UserKickSessions(ctx context.Context, req *v1.UserKickSessionsReq) (res *v1.UserKickSessionsRes, err error)
}
