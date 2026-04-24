package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type MemberListReq struct {
	g.Meta   `path:"/member/list" method:"get" tags:"会员" summary:"会员列表"`
	Page     int    `json:"page" dc:"页码"`
	Size     int    `json:"size" dc:"每页条数"`
	Username string `json:"username" dc:"用户名"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
}

type MemberListRes struct {
	List  []*MemberItem `json:"list" dc:"列表"`
	Total int           `json:"total" dc:"总数"`
}

type MemberItem struct {
	Id         int64  `json:"id" dc:"会员ID"`
	Username   string `json:"username" dc:"用户名"`
	Nickname   string `json:"nickname" dc:"昵称"`
	Email      string `json:"email" dc:"邮箱"`
	Mobile     string `json:"mobile" dc:"手机号"`
	Status     int    `json:"status" dc:"状态：0正常 1停用"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type MemberCreateReq struct {
	g.Meta   `path:"/member" method:"post" tags:"会员" summary:"新增会员"`
	Username string `json:"username" v:"required#err.invalid_params" dc:"用户名"`
	Password string `json:"password" v:"required#err.invalid_params" dc:"密码"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" dc:"邮箱"`
	Mobile   string `json:"mobile" dc:"手机号"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
}

type MemberCreateRes struct {
	Id int64 `json:"id" dc:"会员ID"`
}

type MemberUpdateReq struct {
	g.Meta   `path:"/member/{id}" method:"put" tags:"会员" summary:"更新会员"`
	Id       int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"会员ID"`
	Password string `json:"password" dc:"密码"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" dc:"邮箱"`
	Mobile   string `json:"mobile" dc:"手机号"`
	Status   int    `json:"status" dc:"状态：0正常 1停用"`
}

type MemberUpdateRes struct{}

type MemberDeleteReq struct {
	g.Meta `path:"/member/{id}" method:"delete" tags:"会员" summary:"删除会员"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"会员ID"`
}

type MemberDeleteRes struct{}

type MemberImportReq struct {
	g.Meta `path:"/member/import" method:"post" tags:"会员" summary:"批量导入客户(CSV)"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"CSV 文件（表单字段名 file）"`
}

type MemberImportRes struct {
	Total   int      `json:"total" dc:"有效数据行数（不含空行）"`
	Success int      `json:"success" dc:"成功条数"`
	Failed  int      `json:"failed" dc:"失败条数"`
	Errors  []string `json:"errors" dc:"失败明细（含行号，最多返回若干条）"`
}
