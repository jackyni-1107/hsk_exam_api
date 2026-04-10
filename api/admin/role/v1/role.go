package v1

import "github.com/gogf/gf/v2/frame/g"

type RoleListReq struct {
	g.Meta `path:"/role/list" method:"get" tags:"角色" summary:"角色列表"`
	Page   int    `json:"page" dc:"页码"`
	Size   int    `json:"size" dc:"每页条数"`
	Name   string `json:"name" dc:"角色名称"`
	Status int    `json:"status" dc:"状态：0正常 1停用"`
}

type RoleListRes struct {
	List  []*RoleItem `json:"list" dc:"列表"`
	Total int         `json:"total" dc:"总数"`
}

type RoleItem struct {
	Id         int64   `json:"id" dc:"角色ID"`
	Name       string  `json:"name" dc:"角色名称"`
	Code       string  `json:"code" dc:"角色编码"`
	Status     int     `json:"status" dc:"状态：0正常 1停用"`
	Sort       int     `json:"sort" dc:"排序值"`
	Type       int     `json:"type" dc:"角色类型"`
	Remark     string  `json:"remark" dc:"备注"`
	MenuIds    []int64 `json:"menu_ids" dc:"已分配菜单ID列表"`
	CreateTime string  `json:"create_time" dc:"创建时间"`
}

type RoleCreateReq struct {
	g.Meta `path:"/role" method:"post" tags:"角色" summary:"新增角色"`
	Name   string `json:"name" v:"required#err.invalid_params" dc:"角色名称"`
	Code   string `json:"code" v:"required#err.invalid_params" dc:"角色编码"`
	Status int    `json:"status" dc:"状态：0正常 1停用"`
	Sort   int    `json:"sort" dc:"排序值"`
	Type   int    `json:"type" dc:"角色类型"`
	Remark string `json:"remark" dc:"备注"`
}

type RoleCreateRes struct {
	Id int64 `json:"id" dc:"角色ID"`
}

type RoleUpdateReq struct {
	g.Meta `path:"/role/{id}" method:"put" tags:"角色" summary:"更新角色"`
	Id     int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"角色ID"`
	Name   string `json:"name" dc:"角色名称"`
	Code   string `json:"code" dc:"角色编码"`
	Status int    `json:"status" dc:"状态：0正常 1停用"`
	Sort   int    `json:"sort" dc:"排序值"`
	Type   int    `json:"type" dc:"角色类型"`
	Remark string `json:"remark" dc:"备注"`
}

type RoleUpdateRes struct{}

type RoleDeleteReq struct {
	g.Meta `path:"/role/{id}" method:"delete" tags:"角色" summary:"删除角色"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"角色ID"`
}

type RoleDeleteRes struct{}

type RoleMenuAssignReq struct {
	g.Meta  `path:"/role/{id}/menus" method:"post" tags:"角色" summary:"分配菜单"`
	Id      int64   `json:"id" in:"path" v:"required#err.invalid_params" dc:"角色ID"`
	MenuIds []int64 `json:"menu_ids" dc:"菜单ID列表"`
}

type RoleMenuAssignRes struct{}
