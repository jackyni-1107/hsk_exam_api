package sysmenu

import (
	"context"
	"sort"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	rolesvc "exam/internal/service/sysrole"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysMenu) MenuTree(ctx context.Context) ([]sysentity.SysMenu, error) {
	var list []sysentity.SysMenu
	err := dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort").
		OrderAsc("id").
		Scan(&list)
	return list, err
}

func (s *sSysMenu) VisibleMenusForUser(ctx context.Context, userId int64) ([]sysentity.SysMenu, error) {
	allMenus, err := s.MenuTree(ctx)
	if err != nil {
		return nil, err
	}
	activeMenus := filterActiveMenus(allMenus)
	if len(activeMenus) == 0 {
		return []sysentity.SysMenu{}, nil
	}

	allowedMenuIDs, err := s.MenuIdsForUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return filterVisibleMenus(activeMenus, allowedMenuIDs), nil
}

func (s *sSysMenu) VisibleMenuTreeForUser(ctx context.Context, userId int64) ([]*bo.MenuTreeNode, error) {
	visibleMenus, err := s.VisibleMenusForUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return buildMenuTreeNodes(visibleMenus, 0), nil
}

func (s *sSysMenu) MenuCreate(ctx context.Context, name, permission, path, icon, component, componentName, creator string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) (int64, error) {
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
	if err != nil {
		return 0, err
	}
	var after sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SystemMenu.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sSysMenu) MenuUpdate(ctx context.Context, id int64, name, permission, path, icon, component, componentName, updater string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) error {
	var before sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return err
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeMenuNotFound)
	}
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
	if err != nil {
		return err
	}
	var after sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemMenu.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysMenu) MenuDelete(ctx context.Context, id int64, updater string) error {
	var before sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return err
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeMenuNotFound)
	}
	_, err := dao.SystemMenu.Ctx(ctx).Where("id", id).Data(sysdo.SysMenu{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return err
	}
	var after sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemMenu.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysMenu) MenuIdsForUser(ctx context.Context, userId int64) (map[int64]struct{}, error) {
	menuIDs, err := rolesvc.SysRole().MenuIDsByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]struct{})
	for _, menuID := range menuIDs {
		result[menuID] = struct{}{}
	}
	return result, nil
}

func filterActiveMenus(menus []sysentity.SysMenu) []sysentity.SysMenu {
	active := make([]sysentity.SysMenu, 0, len(menus))
	for _, menu := range menus {
		if menu.Status == consts.StatusNormal {
			active = append(active, menu)
		}
	}
	return active
}

func filterVisibleMenus(activeMenus []sysentity.SysMenu, allowedMenuIDs map[int64]struct{}) []sysentity.SysMenu {
	if len(activeMenus) == 0 || len(allowedMenuIDs) == 0 {
		return []sysentity.SysMenu{}
	}

	byID := make(map[int64]sysentity.SysMenu, len(activeMenus))
	for _, menu := range activeMenus {
		byID[menu.Id] = menu
	}

	visibleMenuIDs := make(map[int64]struct{}, len(allowedMenuIDs))
	for menuID := range allowedMenuIDs {
		for menuID != 0 {
			if _, ok := visibleMenuIDs[menuID]; ok {
				break
			}
			visibleMenuIDs[menuID] = struct{}{}
			menu, ok := byID[menuID]
			if !ok {
				break
			}
			menuID = menu.ParentId
		}
	}

	filtered := make([]sysentity.SysMenu, 0, len(visibleMenuIDs))
	for _, menu := range activeMenus {
		if _, ok := visibleMenuIDs[menu.Id]; ok {
			filtered = append(filtered, menu)
		}
	}
	return filtered
}

func buildMenuTreeNodes(menus []sysentity.SysMenu, parentID int64) []*bo.MenuTreeNode {
	children := make(map[int64][]sysentity.SysMenu)
	for _, menu := range menus {
		children[menu.ParentId] = append(children[menu.ParentId], menu)
	}
	for pid := range children {
		sort.Slice(children[pid], func(i, j int) bool {
			a, b := children[pid][i], children[pid][j]
			if a.Sort != b.Sort {
				return a.Sort < b.Sort
			}
			return a.Id < b.Id
		})
	}

	var walk func(pid int64) []*bo.MenuTreeNode
	walk = func(pid int64) []*bo.MenuTreeNode {
		slice := children[pid]
		out := make([]*bo.MenuTreeNode, 0, len(slice))
		for _, menu := range slice {
			node := &bo.MenuTreeNode{
				Id:            menu.Id,
				Name:          menu.Name,
				Permission:    menu.Permission,
				Type:          menu.Type,
				Sort:          menu.Sort,
				ParentId:      menu.ParentId,
				Path:          menu.Path,
				Icon:          menu.Icon,
				Component:     menu.Component,
				ComponentName: menu.ComponentName,
				Status:        menu.Status,
				Visible:       menu.Visible,
				KeepAlive:     menu.KeepAlive,
				AlwaysShow:    menu.AlwaysShow,
			}
			node.Children = walk(menu.Id)
			out = append(out, node)
		}
		return out
	}

	return walk(parentID)
}
