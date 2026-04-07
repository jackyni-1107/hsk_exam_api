package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
	"exam/internal/consts"
	"exam/internal/logic/security"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) UserKickSessions(ctx context.Context, req *v1.UserKickSessionsReq) (res *v1.UserKickSessionsRes, err error) {
	if err := security.RevokeAllUserSessions(ctx, consts.UserTypeAdmin, req.Id); err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.UserKickSessionsRes{}, nil
}
