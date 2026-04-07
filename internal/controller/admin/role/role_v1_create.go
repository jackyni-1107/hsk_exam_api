package role

import (
	"context"

	v1 "exam/api/admin/role/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) RoleCreate(ctx context.Context, req *v1.RoleCreateReq) (res *v1.RoleCreateRes, err error) {
	var exist sysentity.SysRole
	_ = dao.SystemRole.Ctx(ctx).Where("code", req.Code).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.role_exists")
	}
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	status := req.Status
	if status != consts.StatusNormal && status != consts.StatusDisabled {
		status = consts.StatusNormal
	}
	id, err := dao.SystemRole.Ctx(ctx).InsertAndGetId(sysdo.SysRole{
		Name: req.Name, Code: req.Code, Status: status, Sort: req.Sort, Type: req.Type, Remark: req.Remark,
		Creator: creator, Updater: creator, DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.RoleCreateRes{Id: id}, nil
}
