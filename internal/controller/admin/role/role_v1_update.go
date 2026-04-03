package role

import (
	"context"

	"exam/api/admin/role/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) RoleUpdate(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := dosys.SysRole{Updater: updater}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Code != "" {
		data.Code = req.Code
	}
	data.Status = req.Status
	data.Sort = req.Sort
	data.Type = req.Type
	if req.Remark != "" {
		data.Remark = req.Remark
	}
	_, err = dao.SystemRole.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.RoleUpdateRes{}, nil
}
