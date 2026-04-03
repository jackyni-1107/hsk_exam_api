package member

import (
	"context"

	"exam/api/admin/member/v1"
	"exam/internal/consts"
	daosys "exam/internal/dao/sys"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MemberDelete(ctx context.Context, req *v1.MemberDeleteReq) (res *v1.MemberDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = daosys.SysMember.Ctx(ctx).Where("id", req.Id).Data(dosys.SysMember{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.MemberDeleteRes{}, nil
}
