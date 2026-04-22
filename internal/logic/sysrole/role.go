package sysrole

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func (s *sSysRole) RoleList(ctx context.Context, page, size int, name string, status int) ([]sysentity.SysRole, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	name = strings.TrimSpace(name)
	model := dao.SystemRole.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if name != "" {
		model = model.WhereLike("name", "%"+name+"%")
	}
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		model = model.Where("status", status)
	}
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []sysentity.SysRole
	err = model.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *sSysRole) RoleMenuIds(ctx context.Context, roleId int64) ([]int64, error) {
	roleMenusByRole, err := s.RoleMenuIDsByRoleIDs(ctx, []int64{roleId})
	if err != nil {
		return nil, err
	}
	return roleMenusByRole[roleId], nil
}

func (s *sSysRole) RoleMenuIDsByRoleIDs(ctx context.Context, roleIDs []int64) (map[int64][]int64, error) {
	roleIDs = normalizePositiveIDs(roleIDs)
	if len(roleIDs) == 0 {
		return map[int64][]int64{}, nil
	}
	var rows []sysentity.SysRoleMenu
	err := dao.SystemRoleMenu.Ctx(ctx).
		WhereIn("role_id", roleIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	menuIDsByRole := make(map[int64][]int64, len(roleIDs))
	for _, roleID := range roleIDs {
		menuIDsByRole[roleID] = []int64{}
	}
	for _, row := range rows {
		menuIDsByRole[row.RoleId] = append(menuIDsByRole[row.RoleId], row.MenuId)
	}
	return menuIDsByRole, nil
}

func (s *sSysRole) RoleCreate(ctx context.Context, name, code, remark, creator string, status, sort, typ int) (int64, error) {
	name, code, remark = normalizeRoleInput(name, code, remark)
	if name == "" || code == "" {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	cnt, err := dao.SystemRole.Ctx(ctx).
		Where("code", code).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeRoleCodeExists)
	}

	status = normalizeRoleStatus(status)

	id, err := dao.SystemRole.Ctx(ctx).InsertAndGetId(sysdo.SysRole{
		Name:    name,
		Code:    code,
		Remark:  remark,
		Sort:    sort,
		Status:  status,
		Type:    typ,
		Creator: creator,
		Updater: creator,
	})
	if err != nil {
		return 0, err
	}
	var after sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SystemRole.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sSysRole) PermissionCodesByUser(ctx context.Context, userId int64) ([]string, error) {
	cacheKey := consts.PermCacheKeyPrefix + gconv.String(userId)
	val, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !val.IsEmpty() {
		var perms []string
		if err := json.Unmarshal([]byte(val.String()), &perms); err == nil {
			return perms, nil
		}
	}

	menus, err := activeMenusByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return []string{}, nil
	}

	permSet := make(map[string]struct{}, len(menus))
	for _, menu := range menus {
		if menu.Permission != "" {
			permSet[menu.Permission] = struct{}{}
		}
	}
	perms := make([]string, 0, len(permSet))
	for permission := range permSet {
		perms = append(perms, permission)
	}
	sort.Strings(perms)

	if len(perms) > 0 {
		if data, err := json.Marshal(perms); err == nil {
			_ = g.Redis().SetEX(ctx, cacheKey, string(data), consts.PermCacheTTLSeconds)
		}
	}
	return perms, nil
}

func (s *sSysRole) HasActiveRoleCode(ctx context.Context, userId int64, roleCode string) (bool, error) {
	roleCode = strings.TrimSpace(roleCode)
	if userId <= 0 || roleCode == "" {
		return false, nil
	}
	if userId == consts.SuperAdminUserId && roleCode == consts.RoleCodeSuperAdmin {
		return true, nil
	}

	roleIDs, err := activeRoleIDsByUser(ctx, userId)
	if err != nil {
		return false, err
	}
	if len(roleIDs) == 0 {
		return false, nil
	}

	cnt, err := dao.SystemRole.Ctx(ctx).
		WhereIn("id", roleIDs).
		Where("code", roleCode).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		Count()
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *sSysRole) MenuIDsByUser(ctx context.Context, userId int64) ([]int64, error) {
	menus, err := activeMenusByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(menus) == 0 {
		return []int64{}, nil
	}

	menuIDs := make([]int64, 0, len(menus))
	for _, menu := range menus {
		menuIDs = append(menuIDs, menu.Id)
	}
	sort.Slice(menuIDs, func(i, j int) bool {
		return menuIDs[i] < menuIDs[j]
	})
	return menuIDs, nil
}

func (s *sSysRole) RoleUpdate(ctx context.Context, id int64, name, code, remark, updater string, status, sort, typ int) error {
	before, err := loadRoleByID(ctx, id)
	if err != nil {
		return err
	}
	name, code, remark = normalizeRoleInput(name, code, remark)
	if name == "" {
		name = before.Name
	}
	if code == "" {
		code = before.Code
	}
	if code != "" {
		cnt, err := dao.SystemRole.Ctx(ctx).
			Where("code", code).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			WhereNot("id", id).
			Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			return gerror.NewCode(consts.CodeRoleCodeExists)
		}
	}
	data := sysdo.SysRole{
		Name:    name,
		Code:    code,
		Remark:  remark,
		Sort:    sort,
		Type:    typ,
		Updater: updater,
	}
	shouldClearPermissionCache := false
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		data.Status = status
		shouldClearPermissionCache = before.Status != status
	}
	_, err = dao.SystemRole.Ctx(ctx).Where("id", id).Data(data).Update()
	if err != nil {
		return err
	}
	if shouldClearPermissionCache {
		userIDs, err := userIDsByRoleIDs(ctx, []int64{id})
		if err != nil {
			return err
		}
		bestEffortClearUserPermissionCaches(ctx, userIDs)
	}
	var after sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemRole.Table(), id, before, &after)
	}
	return nil
}

func (s *sSysRole) RoleDelete(ctx context.Context, id int64, updater string) error {
	before, err := loadRoleByID(ctx, id)
	if err != nil {
		return err
	}
	userIDs, err := userIDsByRoleIDs(ctx, []int64{id})
	if err != nil {
		return err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemRole.Table()).Ctx(ctx).Where("id", id).Data(sysdo.SysRole{
			DeleteFlag: consts.DeleteFlagDeleted,
			Updater:    updater,
		}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.SystemRoleMenu.Table()).Ctx(ctx).
			Where("role_id", id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysRoleMenu{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    updater,
			}).Update(); err != nil {
			return err
		}
		_, err := tx.Model(dao.SystemUserRole.Table()).Ctx(ctx).
			Where("role_id", id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysUserRole{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    updater,
			}).Update()
		return err
	})
	if err != nil {
		return err
	}
	bestEffortClearUserPermissionCaches(ctx, userIDs)
	var after sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemRole.Table(), id, before, &after)
	}
	return nil
}

func (s *sSysRole) RoleMenuAssign(ctx context.Context, roleId int64, menuIds []int64, creator string) error {
	if _, err := loadRoleByID(ctx, roleId); err != nil {
		return err
	}
	menuIds = normalizePositiveIDs(menuIds)
	if err := validateMenuIDs(ctx, menuIds); err != nil {
		return err
	}
	beforeIDs, err := s.RoleMenuIds(ctx, roleId)
	if err != nil {
		return err
	}
	beforeStr := utility.JoinSortedInt64IDs(beforeIDs)
	userIDs, err := userIDsByRoleIDs(ctx, []int64{roleId})
	if err != nil {
		return err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemRoleMenu.Table()).Ctx(ctx).
			Where("role_id", roleId).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysRoleMenu{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    creator,
			}).Update(); err != nil {
			return err
		}
		if len(menuIds) == 0 {
			return nil
		}
		batch := buildRoleMenuBatch(roleId, menuIds, creator)
		_, err := tx.Model(dao.SystemRoleMenu.Table()).Ctx(ctx).Data(batch).Insert()
		return err
	})
	if err != nil {
		return err
	}
	bestEffortClearUserPermissionCaches(ctx, userIDs)
	afterStr := utility.JoinSortedInt64IDs(menuIds)
	auditutil.RecordMapDiff(ctx, dao.SystemRoleMenu.Table(), roleId,
		map[string]interface{}{"menu_ids": beforeStr},
		map[string]interface{}{"menu_ids": afterStr})
	return nil
}

func normalizeRoleInput(name, code, remark string) (string, string, string) {
	return strings.TrimSpace(name), strings.TrimSpace(code), strings.TrimSpace(remark)
}

func normalizeRoleStatus(status int) int {
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		return status
	}
	return consts.StatusNormal
}

func normalizePositiveIDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func validateMenuIDs(ctx context.Context, menuIDs []int64) error {
	if len(menuIDs) == 0 {
		return nil
	}
	cnt, err := dao.SystemMenu.Ctx(ctx).
		WhereIn("id", menuIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		Count()
	if err != nil {
		return err
	}
	if cnt != len(menuIDs) {
		return gerror.NewCode(consts.CodeMenuNotFound)
	}
	return nil
}

func loadRoleByID(ctx context.Context, roleID int64) (*sysentity.SysRole, error) {
	var role sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).
		Where("id", roleID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&role); err != nil {
		return nil, err
	}
	if role.Id == 0 {
		return nil, gerror.NewCode(consts.CodeRoleNotFound)
	}
	return &role, nil
}

func buildRoleMenuBatch(roleID int64, menuIDs []int64, actor string) []sysdo.SysRoleMenu {
	batch := make([]sysdo.SysRoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		batch = append(batch, sysdo.SysRoleMenu{
			RoleId:  roleID,
			MenuId:  menuID,
			Creator: actor,
			Updater: actor,
		})
	}
	return batch
}

func userIDsByRoleIDs(ctx context.Context, roleIDs []int64) ([]int64, error) {
	roleIDs = normalizePositiveIDs(roleIDs)
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var rows []sysentity.SysUserRole
	err := dao.SystemUserRole.Ctx(ctx).
		WhereIn("role_id", roleIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, 0, len(rows))
	seen := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.UserId]; ok {
			continue
		}
		seen[row.UserId] = struct{}{}
		userIDs = append(userIDs, row.UserId)
	}
	return userIDs, nil
}

func bestEffortClearUserPermissionCaches(ctx context.Context, userIDs []int64) {
	for _, userID := range normalizePositiveIDs(userIDs) {
		if _, err := g.Redis().Del(ctx, consts.PermCacheKeyPrefix+gconv.String(userID)); err != nil {
			g.Log().Warningf(ctx, "clear user permission cache failed user_id=%d: %v", userID, err)
		}
	}
}

func activeMenusByUser(ctx context.Context, userId int64) ([]sysentity.SysMenu, error) {
	if userId == consts.SuperAdminUserId {
		return activeMenusForSuperAdmin(ctx)
	}

	roleIDs, err := activeRoleIDsByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []sysentity.SysMenu{}, nil
	}

	var roleMenus []sysentity.SysRoleMenu
	if err := dao.SystemRoleMenu.Ctx(ctx).
		WhereIn("role_id", roleIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&roleMenus); err != nil {
		return nil, err
	}
	if len(roleMenus) == 0 {
		return []sysentity.SysMenu{}, nil
	}

	menuIDs := make([]int64, 0, len(roleMenus))
	seenMenuIDs := make(map[int64]struct{}, len(roleMenus))
	for _, roleMenu := range roleMenus {
		if _, ok := seenMenuIDs[roleMenu.MenuId]; ok {
			continue
		}
		seenMenuIDs[roleMenu.MenuId] = struct{}{}
		menuIDs = append(menuIDs, roleMenu.MenuId)
	}

	var menus []sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).
		WhereIn("id", menuIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		Scan(&menus); err != nil {
		return nil, err
	}
	return menus, nil
}

func activeMenusForSuperAdmin(ctx context.Context) ([]sysentity.SysMenu, error) {
	var menus []sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		OrderAsc("sort").
		OrderAsc("id").
		Scan(&menus); err != nil {
		return nil, err
	}
	return menus, nil
}

func activeRoleIDsByUser(ctx context.Context, userId int64) ([]int64, error) {
	var userRoles []sysentity.SysUserRole
	if err := dao.SystemUserRole.Ctx(ctx).
		Where("user_id", userId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&userRoles); err != nil {
		return nil, err
	}
	if len(userRoles) == 0 {
		return []int64{}, nil
	}

	roleIDs := make([]int64, 0, len(userRoles))
	seenRoleIDs := make(map[int64]struct{}, len(userRoles))
	for _, userRole := range userRoles {
		if _, ok := seenRoleIDs[userRole.RoleId]; ok {
			continue
		}
		seenRoleIDs[userRole.RoleId] = struct{}{}
		roleIDs = append(roleIDs, userRole.RoleId)
	}

	var roles []sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).
		WhereIn("id", roleIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		Scan(&roles); err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return []int64{}, nil
	}

	activeRoleIDs := make([]int64, 0, len(roles))
	for _, role := range roles {
		activeRoleIDs = append(activeRoleIDs, role.Id)
	}
	sort.Slice(activeRoleIDs, func(i, j int) bool {
		return activeRoleIDs[i] < activeRoleIDs[j]
	})
	return activeRoleIDs, nil
}
