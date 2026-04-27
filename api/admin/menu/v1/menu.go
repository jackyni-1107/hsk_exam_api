package v1

import "github.com/gogf/gf/v2/frame/g"

type MenuTreeReq struct {
	g.Meta `path:"/menu/tree" method:"get" tags:"菜单" summary:"菜单树" permission:"menu:tree"`
}

type MenuTreeRes struct {
	List []*MenuTreeNode `json:"list" dc:"菜单树列表"`
}

type MenuTreeNode struct {
	Id            int64           `json:"id" dc:"菜单ID"`
	Name          string          `json:"name" dc:"菜单名称"`
	Permission    string          `json:"permission" dc:"权限标识"`
	Type          int             `json:"type" dc:"菜单类型：1目录 2菜单 3按钮"`
	Sort          int             `json:"sort" dc:"排序值"`
	ParentId      int64           `json:"parent_id" dc:"父级ID"`
	Path          string          `json:"path" dc:"路由路径"`
	Icon          string          `json:"icon" dc:"图标"`
	Component     string          `json:"component" dc:"组件路径"`
	ComponentName string          `json:"component_name" dc:"组件名称"`
	Status        int             `json:"status" dc:"状态：0正常 1停用"`
	Visible       bool            `json:"visible" dc:"是否可见"`
	KeepAlive     bool            `json:"keep_alive" dc:"是否缓存"`
	AlwaysShow    bool            `json:"always_show" dc:"是否总是显示"`
	Children      []*MenuTreeNode `json:"children,omitempty" dc:"子菜单"`
}

type MenuCreateReq struct {
	g.Meta        `path:"/menu" method:"post" tags:"菜单" summary:"新增菜单" permission:"menu:create"`
	Name          string `json:"name" v:"required#err.invalid_params" dc:"菜单名称"`
	Permission    string `json:"permission" dc:"权限标识"`
	Type          int    `json:"type" dc:"菜单类型：1目录 2菜单 3按钮"`
	Sort          int    `json:"sort" dc:"排序值"`
	ParentId      int64  `json:"parent_id" dc:"父级ID"`
	Path          string `json:"path" dc:"路由路径"`
	Icon          string `json:"icon" dc:"图标"`
	Component     string `json:"component" dc:"组件路径"`
	ComponentName string `json:"component_name" dc:"组件名称"`
	Status        int    `json:"status" dc:"状态：0正常 1停用"`
	Visible       bool   `json:"visible" dc:"是否可见"`
	KeepAlive     bool   `json:"keep_alive" dc:"是否缓存"`
	AlwaysShow    bool   `json:"always_show" dc:"是否总是显示"`
}

type MenuCreateRes struct {
	Id int64 `json:"id" dc:"菜单ID"`
}

type MenuUpdateReq struct {
	g.Meta        `path:"/menu/{id}" method:"put" tags:"菜单" summary:"更新菜单" permission:"menu:update"`
	Id            int64   `json:"id" in:"path" v:"required#err.invalid_params" dc:"菜单ID"`
	Name          *string `json:"name" dc:"菜单名称"`
	Permission    *string `json:"permission" dc:"权限标识"`
	Type          *int    `json:"type" dc:"菜单类型：1目录 2菜单 3按钮"`
	Sort          *int    `json:"sort" dc:"排序值"`
	ParentId      *int64  `json:"parent_id" dc:"父级ID"`
	Path          *string `json:"path" dc:"路由路径"`
	Icon          *string `json:"icon" dc:"图标"`
	Component     *string `json:"component" dc:"组件路径"`
	ComponentName *string `json:"component_name" dc:"组件名称"`
	Status        *int    `json:"status" dc:"状态：0正常 1停用"`
	Visible       *bool   `json:"visible" dc:"是否可见"`
	KeepAlive     *bool   `json:"keep_alive" dc:"是否缓存"`
	AlwaysShow    *bool   `json:"always_show" dc:"是否总是显示"`
}

type MenuUpdateRes struct{}

type MenuDeleteReq struct {
	g.Meta `path:"/menu/{id}" method:"delete" tags:"菜单" summary:"删除菜单" permission:"menu:delete"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"菜单ID"`
}

type MenuDeleteRes struct{}
