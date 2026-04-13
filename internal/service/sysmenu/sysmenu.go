// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package sysmenu

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysMenu interface {
		MenuTree(ctx context.Context) ([]sysentity.SysMenu, error)
		MenuCreate(ctx context.Context, name string, permission string, path string, icon string, component string, componentName string, creator string, typ int, sort int, parentId int64, status int, visible int, keepAlive int, alwaysShow int) (int64, error)
		MenuUpdate(ctx context.Context, id int64, name string, permission string, path string, icon string, component string, componentName string, updater string, typ int, sort int, parentId int64, status int, visible int, keepAlive int, alwaysShow int) error
		MenuDelete(ctx context.Context, id int64, updater string) error
		MenuIdsForUser(ctx context.Context, userId int64) (map[int64]struct{}, error)
	}
)

var (
	localSysMenu ISysMenu
)

func SysMenu() ISysMenu {
	if localSysMenu == nil {
		panic("implement not found for interface ISysMenu, forgot register?")
	}
	return localSysMenu
}

func RegisterSysMenu(i ISysMenu) {
	localSysMenu = i
}
