package consts

// 通用状态：用户、角色、菜单等
const (
	StatusNormal   = 0 // 正常
	StatusDisabled = 1 // 停用
)

// 逻辑删除标识
const (
	DeleteFlagNotDeleted = 0 // 未删除
	DeleteFlagDeleted    = 1 // 已删除
)
