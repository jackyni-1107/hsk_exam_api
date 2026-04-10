package menu

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
)

func (s *sMenu) MenuTree(ctx context.Context) ([]sysentity.SysMenu, error) {
	var list []sysentity.SysMenu
	err := dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort").
		OrderAsc("id").
		Scan(&list)
	return list, err
}

func (s *sMenu) MenuCreate(ctx context.Context, name, permission, path, icon, component, componentName, creator string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) (int64, error) {
	id, err := dao.SystemMenu.Ctx(ctx).InsertAndGetId(sysdo.SysMenu{
		Name:          name,
		Permission:    permission,
		Path:          path,
		Icon:          icon,
		Component:     component,
		ComponentName: componentName,
		Type:          typ,
		Sort:          sort,
		ParentId:      parentId,
		Status:        status,
		Visible:       visible,
		KeepAlive:     keepAlive,
		AlwaysShow:    alwaysShow,
		Creator:       creator,
		Updater:       creator,
	})
	return id, err
}

func (s *sMenu) MenuUpdate(ctx context.Context, id int64, name, permission, path, icon, component, componentName, updater string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) error {
	_, err := dao.SystemMenu.Ctx(ctx).Where("id", id).Data(sysdo.SysMenu{
		Name:          name,
		Permission:    permission,
		Path:          path,
		Icon:          icon,
		Component:     component,
		ComponentName: componentName,
		Type:          typ,
		Sort:          sort,
		ParentId:      parentId,
		Status:        status,
		Visible:       visible,
		KeepAlive:     keepAlive,
		AlwaysShow:    alwaysShow,
		Updater:       updater,
	}).Update()
	return err
}

func (s *sMenu) MenuDelete(ctx context.Context, id int64, updater string) error {
	_, err := dao.SystemMenu.Ctx(ctx).Where("id", id).Data(sysdo.SysMenu{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	return err
}

func (s *sMenu) MenuIdsForUser(ctx context.Context, userId int64) (map[int64]struct{}, error) {
	var userRoles []sysentity.SysUserRole
	err := dao.SystemUserRole.Ctx(ctx).
		Where("user_id", userId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&userRoles)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]struct{})
	for _, ur := range userRoles {
		var roleMenus []sysentity.SysRoleMenu
		err = dao.SystemRoleMenu.Ctx(ctx).
			Where("role_id", ur.RoleId).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Scan(&roleMenus)
		if err != nil {
			return nil, err
		}
		for _, rm := range roleMenus {
			result[rm.MenuId] = struct{}{}
		}
	}
	return result, nil
}
