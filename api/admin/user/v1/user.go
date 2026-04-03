package v1

import "github.com/gogf/gf/v2/frame/g"

type UserListReq struct {
	g.Meta   `path:"/user/list" method:"get" tags:"用户" summary:"用户列表"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Username string `json:"username"`
	Status   int    `json:"status"`
}

type UserListRes struct {
	List  []*UserItem `json:"list"`
	Total int         `json:"total"`
}

type UserItem struct {
	Id         int64   `json:"id"`
	Username   string  `json:"username"`
	Nickname   string  `json:"nickname"`
	Email      string  `json:"email"`
	Mobile     string  `json:"mobile"`
	Status     int     `json:"status"`
	RoleIds    []int64 `json:"role_ids"`
	CreateTime string  `json:"create_time"`
}

type UserCreateReq struct {
	g.Meta   `path:"/user" method:"post" tags:"用户" summary:"新增用户"`
	Username string  `json:"username" v:"required#err.invalid_params"`
	Password string  `json:"password" v:"required#err.invalid_params"`
	Nickname string  `json:"nickname"`
	Email    string  `json:"email"`
	Mobile   string  `json:"mobile"`
	Status   int     `json:"status"`
	RoleIds  []int64 `json:"role_ids"`
}

type UserCreateRes struct {
	Id int64 `json:"id"`
}

type UserUpdateReq struct {
	g.Meta   `path:"/user/{id}" method:"put" tags:"用户" summary:"更新用户"`
	Id       int64   `json:"id" in:"path" v:"required#err.invalid_params"`
	Password string  `json:"password"`
	Nickname string  `json:"nickname"`
	Email    string  `json:"email"`
	Mobile   string  `json:"mobile"`
	Status   int     `json:"status"`
}

type UserUpdateRes struct{}

type UserDeleteReq struct {
	g.Meta `path:"/user/{id}" method:"delete" tags:"用户" summary:"删除用户"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type UserDeleteRes struct{}

type UserRoleAssignReq struct {
	g.Meta  `path:"/user/{id}/roles" method:"post" tags:"用户" summary:"分配角色"`
	Id      int64   `json:"id" in:"path" v:"required#err.invalid_params"`
	RoleIds []int64 `json:"role_ids"`
}

type UserRoleAssignRes struct{}
