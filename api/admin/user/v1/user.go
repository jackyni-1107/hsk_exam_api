package v1

import "github.com/gogf/gf/v2/frame/g"

type UserListReq struct {
	g.Meta   `path:"/user/list" method:"get" tags:"用户" summary:"用户列表"`
	Page     int    `json:"page" dc:"页码"`
	Size     int    `json:"size" dc:"每页条数"`
	Username string `json:"username" dc:"用户名"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
}

type UserListRes struct {
	List  []*UserItem `json:"list" dc:"列表"`
	Total int         `json:"total" dc:"总数"`
}

type UserItem struct {
	Id         int64   `json:"id" dc:"用户ID"`
	Username   string  `json:"username" dc:"用户名"`
	Nickname   string  `json:"nickname" dc:"昵称"`
	Email      string  `json:"email" dc:"邮箱"`
	Mobile     string  `json:"mobile" dc:"手机号"`
	Status     int     `json:"status" dc:"状态：0正常 1停用"`
	RoleIds    []int64 `json:"role_ids" dc:"角色ID列表"`
	CreateTime string  `json:"create_time" dc:"创建时间"`
}

type UserCreateReq struct {
	g.Meta   `path:"/user" method:"post" tags:"用户" summary:"新增用户"`
	Username string  `json:"username" v:"required#err.invalid_params" dc:"用户名"`
	Password string  `json:"password" v:"required#err.invalid_params" dc:"密码"`
	Nickname string  `json:"nickname" dc:"昵称"`
	Email    string  `json:"email" dc:"邮箱"`
	Mobile   string  `json:"mobile" dc:"手机号"`
	Status   int     `json:"status" dc:"状态：0正常 1停用"`
	RoleIds  []int64 `json:"role_ids" dc:"角色ID列表"`
}

type UserCreateRes struct {
	Id int64 `json:"id" dc:"用户ID"`
}

type UserUpdateReq struct {
	g.Meta   `path:"/user/{id}" method:"put" tags:"用户" summary:"更新用户"`
	Id       int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"用户ID"`
	Password string `json:"password" dc:"密码"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" dc:"邮箱"`
	Mobile   string `json:"mobile" dc:"手机号"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
}

type UserUpdateRes struct{}

type UserDeleteReq struct {
	g.Meta `path:"/user/{id}" method:"delete" tags:"用户" summary:"删除用户"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"用户ID"`
}

type UserDeleteRes struct{}

type UserRoleAssignReq struct {
	g.Meta  `path:"/user/{id}/roles" method:"post" tags:"用户" summary:"分配角色"`
	Id      int64   `json:"id" in:"path" v:"required#err.invalid_params" dc:"用户ID"`
	RoleIds []int64 `json:"role_ids" dc:"角色ID列表"`
}

type UserRoleAssignRes struct{}
