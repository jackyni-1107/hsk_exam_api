// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
package user

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	IUser interface {
		UserList(ctx context.Context, page, size int, username string, status int) ([]sysentity.SysUser, int, error)
		UserRoleIds(ctx context.Context, userId int64) ([]int64, error)
		UserCreate(ctx context.Context, username, password, nickname, email, mobile, creator string, status int, roleIds []int64) (int64, error)
		UserUpdate(ctx context.Context, id int64, password, nickname, email, mobile, updater string, status int) error
		UserDelete(ctx context.Context, id int64, updater string) error
		UserRoleAssign(ctx context.Context, userId int64, roleIds []int64, creator string) error
		FindByUsername(ctx context.Context, username string) (*sysentity.SysUser, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
