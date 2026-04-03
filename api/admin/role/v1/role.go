package v1

import "github.com/gogf/gf/v2/frame/g"

type RoleListReq struct {
	g.Meta `path:"/role/list" method:"get" tags:"角色" summary:"角色列表"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type RoleListRes struct {
	List  []*RoleItem `json:"list"`
	Total int         `json:"total"`
}

type RoleItem struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Status     int     `json:"status"`
	Sort       int     `json:"sort"`
	Type       int     `json:"type"`
	Remark     string  `json:"remark"`
	MenuIds    []int64 `json:"menu_ids"`
	CreateTime string  `json:"create_time"`
}

type RoleCreateReq struct {
	g.Meta `path:"/role" method:"post" tags:"角色" summary:"新增角色"`
	Name   string `json:"name" v:"required#err.invalid_params"`
	Code   string `json:"code" v:"required#err.invalid_params"`
	Status int    `json:"status"`
	Sort   int    `json:"sort"`
	Type   int    `json:"type"`
	Remark string `json:"remark"`
}

type RoleCreateRes struct {
	Id int64 `json:"id"`
}

type RoleUpdateReq struct {
	g.Meta `path:"/role/{id}" method:"put" tags:"角色" summary:"更新角色"`
	Id     int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int    `json:"status"`
	Sort   int    `json:"sort"`
	Type   int    `json:"type"`
	Remark string `json:"remark"`
}

type RoleUpdateRes struct{}

type RoleDeleteReq struct {
	g.Meta `path:"/role/{id}" method:"delete" tags:"角色" summary:"删除角色"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type RoleDeleteRes struct{}

type RoleMenuAssignReq struct {
	g.Meta  `path:"/role/{id}/menus" method:"post" tags:"角色" summary:"分配菜单"`
	Id      int64   `json:"id" in:"path" v:"required#err.invalid_params"`
	MenuIds []int64 `json:"menu_ids"`
}

type RoleMenuAssignRes struct{}
