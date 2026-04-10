// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
package menu

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	IMenu interface {
		MenuTree(ctx context.Context) ([]sysentity.SysMenu, error)
		MenuCreate(ctx context.Context, name, permission, path, icon, component, componentName, creator string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) (int64, error)
		MenuUpdate(ctx context.Context, id int64, name, permission, path, icon, component, componentName, updater string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) error
		MenuDelete(ctx context.Context, id int64, updater string) error
		MenuIdsForUser(ctx context.Context, userId int64) (map[int64]struct{}, error)
	}
)

var (
	localMenu IMenu
)

func Menu() IMenu {
	if localMenu == nil {
		panic("implement not found for interface IMenu, forgot register?")
	}
	return localMenu
}

func RegisterMenu(i IMenu) {
	localMenu = i
}
