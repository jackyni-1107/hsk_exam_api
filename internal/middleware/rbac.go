package middleware

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"

	"exam/internal/consts"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/service/audit"
)

// RBAC 权限校验中间件：校验用户是否拥有指定权限
// permission 为空时跳过校验；非空时要求用户拥有该权限（system_menu.permission）
func RBAC(permission string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		if permission == "" {
			r.Middleware.Next()
			return
		}

		ctxData := GetCtxData(r.GetCtx())
		if ctxData == nil {
			r.SetError(gerror.NewCode(consts.CodeTokenRequired))
			r.ExitAll()
			return
		}

		// 前台用户无角色概念，不校验权限
		if ctxData.UserType == consts.UserTypeClient {
			r.Middleware.Next()
			return
		}

		// 超级管理员跳过权限校验
		if ctxData.UserId == consts.SuperAdminUserId {
			r.Middleware.Next()
			return
		}

		perms, err := getUserPermissions(r.GetCtx(), ctxData.UserId)
		if err != nil {
			r.SetError(gerror.Wrap(err, "get user permissions failed"))
			r.ExitAll()
			return
		}

		if !hasPermission(perms, permission) {
			audit.Audit().RecordSecurityEvent(r.GetCtx(), consts.SecurityEventPermissionDenied, ctxData.UserId, r.GetClientIp(), r.Header.Get("User-Agent"), "permission denied: "+permission, GetTraceId(r.GetCtx()))
			r.SetError(gerror.NewCode(consts.CodePermissionDenied))
			r.ExitAll()
			return
		}

		r.Middleware.Next()
	}
}

func getUserPermissions(ctx context.Context, userId int64) ([]string, error) {
	// 尝试从 Redis 缓存获取
	cacheKey := consts.PermCacheKeyPrefix + gconv.String(userId)
	val, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !val.IsEmpty() {
		var perms []string
		if err := json.Unmarshal([]byte(val.String()), &perms); err == nil {
			return perms, nil
		}
	}

	// 从 DB 查询：user -> roles -> menus -> permissions
	var userRoles []sysentity.SysUserRole
	if err := dao.SystemUserRole.Ctx(ctx).Where("user_id", userId).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&userRoles); err != nil {
		return nil, err
	}
	if len(userRoles) == 0 {
		return []string{}, nil
	}

	roleIds := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleId)
	}

	var roleMenus []sysentity.SysRoleMenu
	if err := dao.SystemRoleMenu.Ctx(ctx).WhereIn("role_id", roleIds).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&roleMenus); err != nil {
		return nil, err
	}
	if len(roleMenus) == 0 {
		return []string{}, nil
	}

	menuIds := make([]int64, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIds = append(menuIds, rm.MenuId)
	}

	var menus []sysentity.SysMenu
	if err := dao.SystemMenu.Ctx(ctx).WhereIn("id", menuIds).Where("delete_flag", consts.DeleteFlagNotDeleted).Where("status", consts.StatusNormal).Scan(&menus); err != nil {
		return nil, err
	}

	permSet := make(map[string]struct{})
	for _, m := range menus {
		if m.Permission != "" {
			permSet[m.Permission] = struct{}{}
		}
	}
	perms := make([]string, 0, len(permSet))
	for p := range permSet {
		perms = append(perms, p)
	}

	// 写入缓存
	if len(perms) > 0 {
		if data, err := json.Marshal(perms); err == nil {
			_ = g.Redis().SetEX(ctx, cacheKey, string(data), consts.PermCacheTTLSeconds)
		}
	}

	return perms, nil
}

// GetUserMenuIds 获取用户有权限的菜单 ID 列表（用于动态菜单）
// 超级管理员返回所有菜单 ID；普通用户从 user->roles->role_menu 查询
func GetUserMenuIds(ctx context.Context, userId int64) ([]int64, error) {
	// 超级管理员：返回所有菜单 ID
	if userId == consts.SuperAdminUserId {
		var menus []sysentity.SysMenu
		if err := dao.SystemMenu.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&menus); err != nil {
			return nil, err
		}
		ids := make([]int64, 0, len(menus))
		for _, m := range menus {
			ids = append(ids, m.Id)
		}
		return ids, nil
	}

	var userRoles []sysentity.SysUserRole
	if err := dao.SystemUserRole.Ctx(ctx).Where("user_id", userId).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&userRoles); err != nil {
		return nil, err
	}
	if len(userRoles) == 0 {
		return []int64{}, nil
	}

	roleIds := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleId)
	}

	var roleMenus []sysentity.SysRoleMenu
	if err := dao.SystemRoleMenu.Ctx(ctx).WhereIn("role_id", roleIds).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&roleMenus); err != nil {
		return nil, err
	}
	if len(roleMenus) == 0 {
		return []int64{}, nil
	}

	idSet := make(map[int64]struct{})
	for _, rm := range roleMenus {
		idSet[rm.MenuId] = struct{}{}
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	return ids, nil
}

func hasPermission(perms []string, required string) bool {
	for _, p := range perms {
		if p == required {
			return true
		}
		// 支持通配符：user:* 匹配 user:list, user:add 等
		if strings.HasSuffix(p, ":*") && strings.HasPrefix(required, strings.TrimSuffix(p, "*")) {
			return true
		}
	}
	return false
}

// RBACFromPath 根据路径和方法推断所需权限
// 例如：GET /api/admin/user/list 或 GET /admin/user/list（无 /api 前缀时） -> user:list
func RBACFromPath(r *ghttp.Request) {
	perm := inferPermission(r.URL.Path, r.Method)
	RBAC(perm)(r)
}

func inferPermission(path, method string) string {
	route, ok := parseAdminRoute(path)
	if !ok {
		return ""
	}
	resource := route.module()
	if resource == "me" {
		return ""
	}
	if perm, ok := matchPermissionRule(method, route); ok {
		return perm
	}
	return inferResourcePermission(resource, method)
}

func inferResourcePermission(resource, method string) string {
	switch method {
	case "GET":
		return resource + ":list"
	case "POST":
		return resource + ":create"
	case "PUT", "PATCH":
		return resource + ":update"
	case "DELETE":
		return resource + ":delete"
	default:
		return ""
	}
}
