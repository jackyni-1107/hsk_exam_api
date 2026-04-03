package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type LogoutReq struct {
	g.Meta `path:"/auth/logout" method:"post" tags:"客户端认证" summary:"客户端登出"`
}

type LogoutRes struct {
}
