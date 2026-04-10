package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
	"exam/internal/middleware"
	membersvc "exam/internal/service/member"
)

func (c *ControllerV1) MemberUpdate(ctx context.Context, req *v1.MemberUpdateReq) (res *v1.MemberUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = membersvc.Member().MemberUpdate(ctx, req.Id, req.Password, req.Nickname, req.Email, req.Mobile, updater, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.MemberUpdateRes{}, nil
}
