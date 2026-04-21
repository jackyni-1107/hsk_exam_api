package user

import (
	"context"

	v1 "exam/api/admin/user/v1"
	usersvc "exam/internal/service/sysuser"
	"exam/internal/utility"
)

func (c *ControllerV1) UserList(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	users, total, err := usersvc.SysUser().UserList(ctx, req.Page, req.Size, req.Username, req.Status)
	if err != nil {
		return nil, err
	}
	userIDs := make([]int64, len(users))
	for i := range users {
		userIDs[i] = users[i].Id
	}
	roleIDsByUser, err := usersvc.SysUser().UserRoleIDsByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.UserItem, 0, len(users))
	for _, u := range users {
		item := &v1.UserItem{
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
		item.RoleIds = roleIDsByUser[u.Id]
		list = append(list, item)
	}

	return &v1.UserListRes{List: list, Total: total}, nil
}
