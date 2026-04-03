package menu

import (
	"context"

	"exam/api/admin/menu/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	dosys "exam/internal/model/do/sys"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MenuUpdate(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := dosys.SysMenu{Updater: updater}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Permission != "" {
		data.Permission = req.Permission
	}
	data.Type = req.Type
	data.Sort = req.Sort
	data.ParentId = req.ParentId
	if req.Path != "" {
		data.Path = req.Path
	}
	if req.Icon != "" {
		data.Icon = req.Icon
	}
	if req.Component != "" {
		data.Component = req.Component
	}
	if req.ComponentName != "" {
		data.ComponentName = req.ComponentName
	}
	data.Status = req.Status
	data.Visible = req.Visible
	data.KeepAlive = req.KeepAlive
	data.AlwaysShow = req.AlwaysShow
	_, err = dao.SystemMenu.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.MenuUpdateRes{}, nil
}
