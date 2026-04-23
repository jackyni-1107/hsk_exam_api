package middleware

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"exam/internal/consts"
	rolesvc "exam/internal/service/sysrole"
)

func RBAC(path string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
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
		permissionCode := r.GetServeHandler().GetMetaTag(consts.KeyPermission)

		if permissionCode != "" {
			// 获取用户权限码（建议通过 UserId 从缓存/数据库读取）
			userPerms, err := getUserPermissions(r.GetCtx(), ctxData.UserId)
			if err != nil {
				r.SetError(gerror.NewCode(consts.CodePermissionDenied))
				r.ExitAll()
				return
			}

			// 校验权限码
			if !hasPermission(userPerms, permissionCode) {
				r.SetError(gerror.NewCode(consts.CodePermissionDenied))
				r.ExitAll()
				return
			}
		}

		r.Middleware.Next()
	}
}

func getUserPermissions(ctx context.Context, userId int64) ([]string, error) {
	return rolesvc.SysRole().PermissionCodesByUser(ctx, userId)
}

func UserHasActiveRoleCode(ctx context.Context, userId int64, roleCode string) (bool, error) {
	return rolesvc.SysRole().HasActiveRoleCode(ctx, userId, roleCode)
}

func hasPermission(perms []string, required string) bool {
	for _, permission := range perms {
		if permission == required {
			return true
		}
	}
	return false
}

func RBACFromPath(r *ghttp.Request) {
	RBAC(r.URL.Path)(r)
}
