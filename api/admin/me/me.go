package me

import (
	"context"

	"exam/api/admin/me/v1"
)

type IMe interface {
	IMenus
}

type IMenus interface {
	Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error)
}
