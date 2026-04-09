package v1

import "github.com/gogf/gf/v2/frame/g"

type MenuTreeReq struct {
	g.Meta `path:"/menu/tree" method:"get" tags:"菜单" summary:"菜单树"`
}

type MenuTreeRes struct {
	List []*MenuTreeNode `json:"list"`
}

type MenuTreeNode struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name"`
	Permission    string          `json:"permission"`
	Type          int             `json:"type"`
	Sort          int             `json:"sort"`
	ParentId      int64           `json:"parent_id"`
	Path          string          `json:"path"`
	Icon          string          `json:"icon"`
	Component     string          `json:"component"`
	ComponentName string          `json:"component_name"`
	Status        int             `json:"status"`
	Visible       bool            `json:"visible"`
	KeepAlive     bool            `json:"keep_alive"`
	AlwaysShow    bool            `json:"always_show"`
	Children      []*MenuTreeNode `json:"children,omitempty"`
}

type MenuCreateReq struct {
	g.Meta        `path:"/menu" method:"post" tags:"菜单" summary:"新增菜单"`
	Name          string `json:"name" v:"required#err.invalid_params"`
	Permission    string `json:"permission"`
	Type          int    `json:"type"`
	Sort          int    `json:"sort"`
	ParentId      int64  `json:"parent_id"`
	Path          string `json:"path"`
	Icon          string `json:"icon"`
	Component     string `json:"component"`
	ComponentName string `json:"component_name"`
	Status        int    `json:"status"`
	Visible       bool   `json:"visible"`
	KeepAlive     bool   `json:"keep_alive"`
	AlwaysShow    bool   `json:"always_show"`
}

type MenuCreateRes struct {
	Id int64 `json:"id"`
}

type MenuUpdateReq struct {
	g.Meta        `path:"/menu/{id}" method:"put" tags:"菜单" summary:"更新菜单"`
	Id            int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	Name          string `json:"name"`
	Permission    string `json:"permission"`
	Type          int    `json:"type"`
	Sort          int    `json:"sort"`
	ParentId      int64  `json:"parent_id"`
	Path          string `json:"path"`
	Icon          string `json:"icon"`
	Component     string `json:"component"`
	ComponentName string `json:"component_name"`
	Status        int    `json:"status"`
	Visible       bool   `json:"visible"`
	KeepAlive     bool   `json:"keep_alive"`
	AlwaysShow    bool   `json:"always_show"`
}

type MenuUpdateRes struct{}

type MenuDeleteReq struct {
	g.Meta `path:"/menu/{id}" method:"delete" tags:"菜单" summary:"删除菜单"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type MenuDeleteRes struct{}
