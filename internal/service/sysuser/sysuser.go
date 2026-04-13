// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package sysuser

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysUser interface {
		UserList(ctx context.Context, page int, size int, username string, status int) ([]sysentity.SysUser, int, error)
		UserRoleIds(ctx context.Context, userId int64) ([]int64, error)
		UserCreate(ctx context.Context, username string, password string, nickname string, email string, mobile string, creator string, status int, roleIds []int64) (int64, error)
		UserUpdate(ctx context.Context, id int64, password string, nickname string, email string, mobile string, updater string, status int) error
		UserDelete(ctx context.Context, id int64, updater string) error
		UserRoleAssign(ctx context.Context, userId int64, roleIds []int64, creator string) error
		FindByUsername(ctx context.Context, username string) (*sysentity.SysUser, error)
	}
)

var (
	localSysUser ISysUser
)

func SysUser() ISysUser {
	if localSysUser == nil {
		panic("implement not found for interface ISysUser, forgot register?")
	}
	return localSysUser
}

func RegisterSysUser(i ISysUser) {
	localSysUser = i
}
