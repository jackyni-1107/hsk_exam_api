package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
)

type IRole interface {
	RoleList(ctx context.Context, req *v1.RoleListReq) (res *v1.RoleListRes, err error)
	RoleCreate(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error)
	RoleUpdate(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error)
	RoleDelete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error)
	RoleMenuAssign(ctx context.Context, req *v1.RoleMenuAssignReq) (res *v1.RoleMenuAssignRes, err error)
}
