package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
	membersvc "exam/internal/service/member"
	"exam/internal/utility"
)

func (c *ControllerV1) MemberList(ctx context.Context, req *v1.MemberListReq) (res *v1.MemberListRes, err error) {
	rows, total, err := membersvc.Member().MemberList(ctx, req.Page, req.Size, req.Username, req.Status)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.MemberItem, 0, len(rows))
	for _, u := range rows {
		item := &v1.MemberItem{
			Id:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
			Mobile:   u.Mobile,
			Status:   u.Status,
		}
		if u.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(u.CreateTime)
		}
		list = append(list, item)
	}

	return &v1.MemberListRes{List: list, Total: total}, nil
}
