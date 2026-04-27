package v1

import "github.com/gogf/gf/v2/frame/g"

type UserKickSessionsReq struct {
	g.Meta `path:"/user/{id}/kick-sessions" method:"post" tags:"用户管理" summary:"强制下线该用户全部会话" permission:"user:kick_sessions"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"用户ID"`
}

type UserKickSessionsRes struct{}
