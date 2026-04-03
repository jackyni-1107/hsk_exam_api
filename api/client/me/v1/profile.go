package v1

import "github.com/gogf/gf/v2/frame/g"

type ProfileReq struct {
	g.Meta `path:"/me/profile" method:"get" tags:"客户端-当前用户" summary:"我的信息"`
}

type ProfileRes struct {
	Id                 int64  `json:"id" dc:"用户ID"`
	Username           string `json:"username" dc:"账号"`
	Nickname           string `json:"nickname" dc:"昵称"`
	Avatar             string `json:"avatar" dc:"头像"`
	Email              string `json:"email" dc:"邮箱"`
	Mobile             string `json:"mobile" dc:"手机"`
	Status             int    `json:"status" dc:"状态：0正常 1停用"`
	MustChangePassword int    `json:"must_change_password" dc:"是否须改密：0否 1是"`
	PasswordChangedAt  string `json:"password_changed_at,omitempty" dc:"密码最后修改时间"`
	LoginIp            string `json:"login_ip,omitempty" dc:"最后登录IP"`
	LoginTime          string `json:"login_time,omitempty" dc:"最后登录时间"`
	CreateTime         string `json:"create_time,omitempty" dc:"注册时间"`
}
