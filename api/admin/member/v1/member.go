package v1

import "github.com/gogf/gf/v2/frame/g"

type MemberListReq struct {
	g.Meta   `path:"/member/list" method:"get" tags:"会员" summary:"会员列表"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Username string `json:"username"`
	Status   int    `json:"status"`
}

type MemberListRes struct {
	List  []*MemberItem `json:"list"`
	Total int           `json:"total"`
}

type MemberItem struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

type MemberCreateReq struct {
	g.Meta   `path:"/member" method:"post" tags:"会员" summary:"新增会员"`
	Username string `json:"username" v:"required#err.invalid_params"`
	Password string `json:"password" v:"required#err.invalid_params"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Status   int    `json:"status"`
}

type MemberCreateRes struct {
	Id int64 `json:"id"`
}

type MemberUpdateReq struct {
	g.Meta   `path:"/member/{id}" method:"put" tags:"会员" summary:"更新会员"`
	Id       int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Status   int    `json:"status"`
}

type MemberUpdateRes struct{}

type MemberDeleteReq struct {
	g.Meta `path:"/member/{id}" method:"delete" tags:"会员" summary:"删除会员"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type MemberDeleteRes struct{}
