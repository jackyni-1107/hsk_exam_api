package member

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"exam/api/admin/member/v1"
	"exam/internal/consts"
	daosys "exam/internal/dao/sys"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"exam/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MemberCreate(ctx context.Context, req *v1.MemberCreateReq) (res *v1.MemberCreateRes, err error) {
	var exist entity.SystemMember
	_ = daosys.SysMember.Ctx(ctx).Where("username", req.Username).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.member_exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	status := req.Status
	if status != consts.StatusNormal && status != consts.StatusDisabled {
		status = consts.StatusNormal
	}
	id, err := daosys.SysMember.Ctx(ctx).InsertAndGetId(dosys.SysMember{
		Username: req.Username, Password: string(hash), Nickname: req.Nickname,
		Email: req.Email, Mobile: req.Mobile, Status: status,
		Creator: creator, Updater: creator, DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.MemberCreateRes{Id: id}, nil
}
