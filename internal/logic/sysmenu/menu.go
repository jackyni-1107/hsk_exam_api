package sysmenu

import (
	"context"
	"sort"
	"strings"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/bo"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	rolesvc "exam/internal/service/sysrole"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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

func (s *sSysMenu) MenuTreeNodes(ctx context.Context) ([]*bo.MenuTreeNode, error) {
	allMenus, err := s.MenuTree(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTreeNodes(allMenus, 0), nil
}

func (s *sSysMenu) VisibleMenuTreeForUser(ctx context.Context, userId int64) ([]*bo.MenuTreeNode, error) {
	visibleMenus, err := s.VisibleMenusForUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return buildMenuTreeNodes(visibleMenus, 0), nil
}

func (s *sSysMenu) MenuCreate(ctx context.Context, name, permission, path, icon, component, componentName, creator string, typ, sort int, parentId int64, status, visible, keepAlive, alwaysShow int) (int64, error) {
	state := newMenuCreateState(name, permission, path, icon, component, componentName, typ, sort, parentId, status, visible, keepAlive, alwaysShow)
	if err := validateMenuState(ctx, 0, state); err != nil {
		return 0, err
	}
	id, err := dao.SystemMenu.Ctx(ctx).InsertAndGetId(sysdo.SysMenu{
		Name:          state.Name,
		Permission:    state.Permission,
		Path:          state.Path,
		Icon:          state.Icon,
		Component:     state.Component,
		ComponentName: state.ComponentName,
		Type:          state.Type,
		Sort:          state.Sort,
		ParentId:      state.ParentID,
		Status:        state.Status,
		Visible:       boolToInt(state.Visible),
		KeepAlive:     boolToInt(state.KeepAlive),
		AlwaysShow:    boolToInt(state.AlwaysShow),
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

func (s *sSysMenu) MenuUpdate(ctx context.Context, id int64, input bo.MenuUpdateInput, updater string) error {
	before, err := loadMenuByID(ctx, id)
	if err != nil {
		return err
	}
	beforeState := menuStateFromEntity(*before)
	afterState := resolveMenuState(before, input)
	if err := validateMenuState(ctx, id, afterState); err != nil {
		return err
	}
	affectedUserIDs, err := affectedUserIDsByMenuIDs(ctx, []int64{id})
	if err != nil {
		return err
	}
	_, err = dao.SystemMenu.Ctx(ctx).Where("id", id).Data(buildMenuUpdateData(beforeState, afterState, updater)).Update()
	if err != nil {
		return err
	}
	bestEffortClearUserPermissionCaches(ctx, affectedUserIDs)
	var after sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemMenu.Table(), id, before, &after)
	}
	return nil
}

func (s *sSysMenu) MenuDelete(ctx context.Context, id int64, updater string) error {
	before, err := loadMenuByID(ctx, id)
	if err != nil {
		return err
	}
	menuIDs, err := menuSubtreeIDs(ctx, id)
	if err != nil {
		return err
	}
	affectedUserIDs, err := affectedUserIDsByMenuIDs(ctx, menuIDs)
	if err != nil {
		return err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemMenu.Table()).Ctx(ctx).
			WhereIn("id", menuIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysMenu{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    updater,
			}).Update(); err != nil {
			return err
		}
		_, err := tx.Model(dao.SystemRoleMenu.Table()).Ctx(ctx).
			WhereIn("menu_id", menuIDs).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysRoleMenu{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    updater,
			}).Update()
		return err
	})
	if err != nil {
		return err
	}
	bestEffortClearUserPermissionCaches(ctx, affectedUserIDs)
	var after sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemMenu.Table(), id, before, &after)
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

type menuState struct {
	Name          string
	Permission    string
	Type          int
	Sort          int
	Status        int
	ParentID      int64
	Path          string
	Icon          string
	Component     string
	ComponentName string
	Visible       bool
	KeepAlive     bool
	AlwaysShow    bool
}

func newMenuCreateState(name, permission, path, icon, component, componentName string, typ, sort int, parentID int64, status, visible, keepAlive, alwaysShow int) menuState {
	return normalizeCreateMenuState(menuState{
		Name:          name,
		Permission:    permission,
		Type:          typ,
		Sort:          sort,
		Status:        status,
		ParentID:      parentID,
		Path:          path,
		Icon:          icon,
		Component:     component,
		ComponentName: componentName,
		Visible:       visible != 0,
		KeepAlive:     keepAlive != 0,
		AlwaysShow:    alwaysShow != 0,
	})
}

func menuStateFromEntity(menu sysentity.SysMenu) menuState {
	return menuState{
		Name:          menu.Name,
		Permission:    menu.Permission,
		Type:          menu.Type,
		Sort:          menu.Sort,
		Status:        menu.Status,
		ParentID:      menu.ParentId,
		Path:          menu.Path,
		Icon:          menu.Icon,
		Component:     menu.Component,
		ComponentName: menu.ComponentName,
		Visible:       menu.Visible,
		KeepAlive:     menu.KeepAlive,
		AlwaysShow:    menu.AlwaysShow,
	}
}

func normalizeCreateMenuState(state menuState) menuState {
	state.Name = normalizeMenuString(state.Name)
	state.Permission = normalizeMenuString(state.Permission)
	state.Path = normalizeMenuString(state.Path)
	state.Icon = normalizeMenuString(state.Icon)
	state.Component = normalizeMenuString(state.Component)
	state.ComponentName = normalizeMenuString(state.ComponentName)
	return sanitizeMenuState(state)
}

func normalizeMenuString(value string) string {
	return strings.TrimSpace(value)
}

func sanitizeMenuState(state menuState) menuState {
	if state.Type == consts.MenuTypeButton {
		state.Path = ""
		state.Icon = ""
		state.Component = ""
		state.ComponentName = ""
		state.KeepAlive = false
		state.AlwaysShow = false
	}
	return state
}

func resolveMenuState(before *sysentity.SysMenu, input bo.MenuUpdateInput) menuState {
	state := menuStateFromEntity(*before)
	if input.Name != nil {
		state.Name = normalizeMenuString(*input.Name)
	}
	if input.Permission != nil {
		state.Permission = normalizeMenuString(*input.Permission)
	}
	if input.Type != nil {
		state.Type = *input.Type
	}
	if input.Sort != nil {
		state.Sort = *input.Sort
	}
	if input.Status != nil {
		state.Status = *input.Status
	}
	if input.ParentID != nil {
		state.ParentID = *input.ParentID
	}
	if input.Path != nil {
		state.Path = normalizeMenuString(*input.Path)
	}
	if input.Icon != nil {
		state.Icon = normalizeMenuString(*input.Icon)
	}
	if input.Component != nil {
		state.Component = normalizeMenuString(*input.Component)
	}
	if input.ComponentName != nil {
		state.ComponentName = normalizeMenuString(*input.ComponentName)
	}
	if input.Visible != nil {
		state.Visible = *input.Visible
	}
	if input.KeepAlive != nil {
		state.KeepAlive = *input.KeepAlive
	}
	if input.AlwaysShow != nil {
		state.AlwaysShow = *input.AlwaysShow
	}
	return sanitizeMenuState(state)
}

func validateMenuState(ctx context.Context, currentMenuID int64, state menuState) error {
	var menus []sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&menus); err != nil {
		return err
	}
	return validateMenuStructure(menus, currentMenuID, state)
}

func validateMenuStructure(menus []sysentity.SysMenu, currentMenuID int64, state menuState) error {
	name := normalizeMenuString(state.Name)
	permission := normalizeMenuString(state.Permission)
	path := normalizeMenuString(state.Path)
	component := normalizeMenuString(state.Component)
	componentName := normalizeMenuString(state.ComponentName)

	if !isValidMenuType(state.Type) || !isValidMenuStatus(state.Status) || state.ParentID < 0 || state.Sort < 0 || name == "" {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	switch state.Type {
	case consts.MenuTypeDir:
		if permission != "" {
			return gerror.NewCode(consts.CodeInvalidParams)
		}
	case consts.MenuTypeMenu:
		if path == "" || component == "" {
			return gerror.NewCode(consts.CodeInvalidParams)
		}
		if state.KeepAlive && componentName == "" {
			return gerror.NewCode(consts.CodeInvalidParams)
		}
	case consts.MenuTypeButton:
		if permission == "" || state.ParentID == 0 {
			return gerror.NewCode(consts.CodeInvalidParams)
		}
	}
	if currentMenuID > 0 && state.ParentID == currentMenuID {
		return gerror.NewCode(consts.CodeInvalidParams)
	}
	if state.ParentID == 0 {
		if currentMenuID > 0 && state.Type == consts.MenuTypeButton {
			subtreeIDs := collectMenuSubtreeIDs(menus, currentMenuID)
			if len(subtreeIDs) > 1 {
				return gerror.NewCode(consts.CodeInvalidParams)
			}
		}
		return nil
	}

	menuByID := make(map[int64]sysentity.SysMenu, len(menus))
	for _, menu := range menus {
		menuByID[menu.Id] = menu
	}
	parent, ok := menuByID[state.ParentID]
	if !ok {
		return gerror.NewCode(consts.CodeMenuNotFound)
	}
	if parent.Type == consts.MenuTypeButton {
		return gerror.NewCode(consts.CodeInvalidParams)
	}

	if currentMenuID > 0 {
		subtreeIDs := collectMenuSubtreeIDs(menus, currentMenuID)
		if state.Type == consts.MenuTypeButton && len(subtreeIDs) > 1 {
			return gerror.NewCode(consts.CodeInvalidParams)
		}
		for _, subtreeID := range subtreeIDs {
			if subtreeID == state.ParentID {
				return gerror.NewCode(consts.CodeInvalidParams)
			}
		}
	}
	return nil
}

func buildMenuUpdateData(before, after menuState, updater string) sysdo.SysMenu {
	data := sysdo.SysMenu{
		Updater: updater,
	}
	if after.Name != before.Name {
		data.Name = after.Name
	}
	if after.Permission != before.Permission {
		data.Permission = after.Permission
	}
	if after.Type != before.Type {
		data.Type = after.Type
	}
	if after.Sort != before.Sort {
		data.Sort = after.Sort
	}
	if after.ParentID != before.ParentID {
		data.ParentId = after.ParentID
	}
	if after.Path != before.Path {
		data.Path = after.Path
	}
	if after.Icon != before.Icon {
		data.Icon = after.Icon
	}
	if after.Component != before.Component {
		data.Component = after.Component
	}
	if after.ComponentName != before.ComponentName {
		data.ComponentName = after.ComponentName
	}
	if after.Status != before.Status {
		data.Status = after.Status
	}
	if after.Visible != before.Visible {
		data.Visible = boolToInt(after.Visible)
	}
	if after.KeepAlive != before.KeepAlive {
		data.KeepAlive = boolToInt(after.KeepAlive)
	}
	if after.AlwaysShow != before.AlwaysShow {
		data.AlwaysShow = boolToInt(after.AlwaysShow)
	}
	return data
}

func loadMenuByID(ctx context.Context, menuID int64) (*sysentity.SysMenu, error) {
	var menu sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).
		Where("id", menuID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&menu); err != nil {
		return nil, err
	}
	if menu.Id == 0 {
		return nil, gerror.NewCode(consts.CodeMenuNotFound)
	}
	return &menu, nil
}

func menuSubtreeIDs(ctx context.Context, rootMenuID int64) ([]int64, error) {
	var menus []sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&menus); err != nil {
		return nil, err
	}
	return collectMenuSubtreeIDs(menus, rootMenuID), nil
}

func collectMenuSubtreeIDs(menus []sysentity.SysMenu, rootMenuID int64) []int64 {
	if rootMenuID <= 0 {
		return nil
	}
	children := make(map[int64][]int64)
	rootExists := false
	for _, menu := range menus {
		if menu.Id == rootMenuID {
			rootExists = true
		}
		children[menu.ParentId] = append(children[menu.ParentId], menu.Id)
	}
	if !rootExists {
		return nil
	}

	out := make([]int64, 0, len(menus))
	stack := []int64{rootMenuID}
	seen := make(map[int64]struct{}, len(menus))
	for len(stack) > 0 {
		last := len(stack) - 1
		menuID := stack[last]
		stack = stack[:last]
		if _, ok := seen[menuID]; ok {
			continue
		}
		seen[menuID] = struct{}{}
		out = append(out, menuID)
		stack = append(stack, children[menuID]...)
	}
	return out
}

func affectedUserIDsByMenuIDs(ctx context.Context, menuIDs []int64) ([]int64, error) {
	menuIDs = normalizePositiveMenuIDs(menuIDs)
	if len(menuIDs) == 0 {
		return nil, nil
	}

	var roleMenus []sysentity.SysRoleMenu
	if err := dao.SystemRoleMenu.Ctx(ctx).
		WhereIn("menu_id", menuIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&roleMenus); err != nil {
		return nil, err
	}
	if len(roleMenus) == 0 {
		return nil, nil
	}

	roleIDs := make([]int64, 0, len(roleMenus))
	seenRoleIDs := make(map[int64]struct{}, len(roleMenus))
	for _, roleMenu := range roleMenus {
		if _, ok := seenRoleIDs[roleMenu.RoleId]; ok {
			continue
		}
		seenRoleIDs[roleMenu.RoleId] = struct{}{}
		roleIDs = append(roleIDs, roleMenu.RoleId)
	}

	var userRoles []sysentity.SysUserRole
	if err := dao.SystemUserRole.Ctx(ctx).
		WhereIn("role_id", roleIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&userRoles); err != nil {
		return nil, err
	}
	if len(userRoles) == 0 {
		return nil, nil
	}

	userIDs := make([]int64, 0, len(userRoles))
	seenUserIDs := make(map[int64]struct{}, len(userRoles))
	for _, userRole := range userRoles {
		if _, ok := seenUserIDs[userRole.UserId]; ok {
			continue
		}
		seenUserIDs[userRole.UserId] = struct{}{}
		userIDs = append(userIDs, userRole.UserId)
	}
	return userIDs, nil
}

func normalizePositiveMenuIDs(menuIDs []int64) []int64 {
	if len(menuIDs) == 0 {
		return nil
	}
	out := make([]int64, 0, len(menuIDs))
	seen := make(map[int64]struct{}, len(menuIDs))
	for _, menuID := range menuIDs {
		if menuID <= 0 {
			continue
		}
		if _, ok := seen[menuID]; ok {
			continue
		}
		seen[menuID] = struct{}{}
		out = append(out, menuID)
	}
	return out
}

func isValidMenuType(menuType int) bool {
	return menuType == consts.MenuTypeDir || menuType == consts.MenuTypeMenu || menuType == consts.MenuTypeButton
}

func isValidMenuStatus(status int) bool {
	return status == consts.StatusNormal || status == consts.StatusDisabled
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func bestEffortClearUserPermissionCaches(ctx context.Context, userIDs []int64) {
	for _, userID := range normalizePositiveMenuIDs(userIDs) {
		if _, err := g.Redis().Del(ctx, consts.PermCacheKeyPrefix+gconv.String(userID)); err != nil {
			g.Log().Warningf(ctx, "clear user permission cache failed user_id=%d: %v", userID, err)
		}
	}
}
