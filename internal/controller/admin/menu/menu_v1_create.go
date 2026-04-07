package menu

import (
	"context"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MenuCreate(ctx context.Context, req *v1.MenuCreateReq) (res *v1.MenuCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := dao.SystemMenu.Ctx(ctx).InsertAndGetId(sysdo.SysMenu{
		Name: req.Name, Permission: req.Permission, Type: req.Type, Sort: req.Sort, ParentId: req.ParentId,
		Path: req.Path, Icon: req.Icon, Component: req.Component, ComponentName: req.ComponentName,
		Status: req.Status, Visible: req.Visible, KeepAlive: req.KeepAlive, AlwaysShow: req.AlwaysShow,
		Creator: creator, Updater: creator, DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.MenuCreateRes{Id: id}, nil
}
