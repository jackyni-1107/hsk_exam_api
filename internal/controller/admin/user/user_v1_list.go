package user

import (
	"context"

	"exam/api/admin/user/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) UserList(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	model := dao.SystemUser.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Username != "" {
		model = model.WhereLike("username", "%"+req.Username+"%")
	}
	// status: -1=全部 0=正常 1=停用，仅当明确指定时过滤
	if req.Status == consts.StatusNormal || req.Status == consts.StatusDisabled {
		model = model.Where("status", req.Status)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}

	var users []entity.SystemUser
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&users)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}

	list := make([]*v1.UserItem, 0, len(users))
	for _, u := range users {
		item := &v1.UserItem{
			Id:         u.Id,
			Username:   u.Username,
			Nickname:   u.Nickname,
			Email:      u.Email,
			Mobile:     u.Mobile,
			Status:     u.Status,
			CreateTime: "",
		}
		if u.CreateTime != nil {
			item.CreateTime = u.CreateTime.Format("Y-m-d H:i:s")
		}
		// 查询角色
		var userRoles []entity.SystemUserRole
		_ = dao.SystemUserRole.Ctx(ctx).Where("user_id", u.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&userRoles)
		for _, ur := range userRoles {
			item.RoleIds = append(item.RoleIds, ur.RoleId)
		}
		list = append(list, item)
	}

	return &v1.UserListRes{List: list, Total: total}, nil
}
