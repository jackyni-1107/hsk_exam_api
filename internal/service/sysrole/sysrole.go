// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package sysrole

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysRole interface {
		RoleList(ctx context.Context, page int, size int, name string, status int) ([]sysentity.SysRole, int, error)
		RoleMenuIds(ctx context.Context, roleId int64) ([]int64, error)
		RoleCreate(ctx context.Context, name string, code string, remark string, creator string, status int, sort int, typ int) (int64, error)
		RoleUpdate(ctx context.Context, id int64, name string, code string, remark string, updater string, status int, sort int, typ int) error
		RoleDelete(ctx context.Context, id int64, updater string) error
		RoleMenuAssign(ctx context.Context, roleId int64, menuIds []int64, creator string) error
	}
)

var (
	localSysRole ISysRole
)

func SysRole() ISysRole {
	if localSysRole == nil {
		panic("implement not found for interface ISysRole, forgot register?")
	}
	return localSysRole
}

func RegisterSysRole(i ISysRole) {
	localSysRole = i
}
