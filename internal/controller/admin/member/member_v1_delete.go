package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
	"exam/internal/middleware"
	membersvc "exam/internal/service/member"
)

func (c *ControllerV1) MemberDelete(ctx context.Context, req *v1.MemberDeleteReq) (res *v1.MemberDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = membersvc.Member().MemberDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.MemberDeleteRes{}, nil
}
