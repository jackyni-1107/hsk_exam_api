// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
package role

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	IRole interface {
		RoleList(ctx context.Context, page, size int, name string, status int) ([]sysentity.SysRole, int, error)
		RoleMenuIds(ctx context.Context, roleId int64) ([]int64, error)
		RoleCreate(ctx context.Context, name, code, remark, creator string, status, sort, typ int) (int64, error)
		RoleUpdate(ctx context.Context, id int64, name, code, remark, updater string, status, sort, typ int) error
		RoleDelete(ctx context.Context, id int64, updater string) error
		RoleMenuAssign(ctx context.Context, roleId int64, menuIds []int64, creator string) error
	}
)

var (
	localRole IRole
)

func Role() IRole {
	if localRole == nil {
		panic("implement not found for interface IRole, forgot register?")
	}
	return localRole
}

func RegisterRole(i IRole) {
	localRole = i
}
