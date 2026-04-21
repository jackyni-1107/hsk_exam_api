package middleware

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
	"exam/internal/service/audit"
	menusvc "exam/internal/service/sysmenu"
	rolesvc "exam/internal/service/sysrole"
)

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

		if ctxData.UserType == consts.UserTypeClient {
			r.Middleware.Next()
			return
		}

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
	return rolesvc.SysRole().PermissionCodesByUser(ctx, userId)
}

func GetUserMenuIds(ctx context.Context, userId int64) ([]int64, error) {
	if userId == consts.SuperAdminUserId {
		menus, err := menusvc.SysMenu().MenuTree(ctx)
		if err != nil {
			return nil, err
		}
		ids := make([]int64, 0, len(menus))
		for _, menu := range menus {
			if menu.Status == consts.StatusNormal {
				ids = append(ids, menu.Id)
			}
		}
		return ids, nil
	}
	return rolesvc.SysRole().MenuIDsByUser(ctx, userId)
}

func hasPermission(perms []string, required string) bool {
	for _, permission := range perms {
		if permission == required {
			return true
		}
		if strings.HasSuffix(permission, ":*") && strings.HasPrefix(required, strings.TrimSuffix(permission, "*")) {
			return true
		}
	}
	return false
}

func RBACFromPath(r *ghttp.Request) {
	perm := inferPermission(r.URL.Path, r.Method)
	RBAC(perm)(r)
}

func inferPermission(path, method string) string {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] == "api" {
		parts = parts[1:]
	}
	if len(parts) < 2 || parts[0] != "admin" {
		return ""
	}
	resource := parts[1]
	if resource == "me" {
		return ""
	}
	if resource == "exam" && len(parts) >= 4 && parts[2] == "paper" && parts[3] == "import" {
		return "exam:import"
	}
	if resource == "exam" && len(parts) >= 3 && parts[2] == "attempt" && method == "GET" {
		return "exam:result:list"
	}
	if resource == "exam" && len(parts) >= 5 && parts[2] == "attempt" && parts[4] == "subjective-scores" && method == "PUT" {
		return "exam:result:grade"
	}
	if resource == "file" && len(parts) >= 4 && parts[3] == "upload" && method == "POST" {
		return "file:list"
	}
	if resource == "task" && len(parts) >= 3 {
		if parts[2] == "run" {
			return "task:run"
		}
		if parts[2] == "log" {
			return "task:log"
		}
	}
	switch method {
	case "GET":
		return resource + ":list"
	case "POST":
		if resource == "user" && len(parts) >= 4 && parts[3] == "kick-sessions" {
			return "user:update"
		}
		return resource + ":create"
	case "PUT", "PATCH":
		return resource + ":update"
	case "DELETE":
		return resource + ":delete"
	default:
		return ""
	}
}
