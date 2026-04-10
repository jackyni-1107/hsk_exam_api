package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
	"exam/internal/middleware"
	membersvc "exam/internal/service/member"
)

func (c *ControllerV1) MemberCreate(ctx context.Context, req *v1.MemberCreateReq) (res *v1.MemberCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := membersvc.Member().MemberCreate(ctx, req.Username, req.Password, req.Nickname, req.Email, req.Mobile, creator, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.MemberCreateRes{Id: id}, nil
}
