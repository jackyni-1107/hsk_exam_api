package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MemberList(ctx context.Context, req *v1.MemberListReq) (res *v1.MemberListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := sysdao.SysMember.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Username != "" {
		model = model.WhereLike("username", "%"+req.Username+"%")
	}
	if req.Status == consts.StatusNormal || req.Status == consts.StatusDisabled {
		model = model.Where("status", req.Status)
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var rows []sysentity.SysMember
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&rows)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	list := make([]*v1.MemberItem, 0, len(rows))
	for _, u := range rows {
		item := &v1.MemberItem{
			Id: u.Id, Username: u.Username, Nickname: u.Nickname,
			Email: u.Email, Mobile: u.Mobile, Status: u.Status,
		}
		if u.CreateTime != nil {
			item.CreateTime = util.ToRFC3339UTC(u.CreateTime)
		}
		list = append(list, item)
	}
	return &v1.MemberListRes{List: list, Total: total}, nil
}
