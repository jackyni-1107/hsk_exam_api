package consts

// 系统预置 ID
const (
	// AnonymousUserId 匿名用户 ID（未登录时的占位值）
	AnonymousUserId = 0
	// SuperAdminUserId 超级管理员用户 ID（003_seed_admin 中初始化的 admin 用户，通常为 1）
	SuperAdminUserId = 1
	// RoleCodeSuperAdmin 系统角色 code（sys_role.code），仅该角色可执行试卷物理删除等高危操作
	RoleCodeSuperAdmin = "super_admin"
)
