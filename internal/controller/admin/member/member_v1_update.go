package member

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"exam/api/admin/member/v1"
	"exam/internal/consts"
	daosys "exam/internal/dao/sys"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MemberUpdate(ctx context.Context, req *v1.MemberUpdateReq) (res *v1.MemberUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := dosys.SysMember{Updater: updater}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
		}
		data.Password = string(hash)
	}
	if req.Nickname != "" {
		data.Nickname = req.Nickname
	}
	if req.Email != "" {
		data.Email = req.Email
	}
	if req.Mobile != "" {
		data.Mobile = req.Mobile
	}
	data.Status = req.Status
	_, err = daosys.SysMember.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.MemberUpdateRes{}, nil
}
