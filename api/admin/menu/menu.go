package menu

import (
	"context"

	"exam/api/admin/menu/v1"
)

type IMenu interface {
	IMenuTree
	IMenuCreate
	IMenuUpdate
	IMenuDelete
}

type IMenuTree interface {
	MenuTree(ctx context.Context, req *v1.MenuTreeReq) (res *v1.MenuTreeRes, err error)
}

type IMenuCreate interface {
	MenuCreate(ctx context.Context, req *v1.MenuCreateReq) (res *v1.MenuCreateRes, err error)
}

type IMenuUpdate interface {
	MenuUpdate(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error)
}

type IMenuDelete interface {
	MenuDelete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error)
}
