package me

import (
	"context"

	v1 "exam/api/admin/me/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	"exam/internal/model/bo"
	menusvc "exam/internal/service/sysmenu"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}

	menuTree, err := menusvc.SysMenu().VisibleMenuTreeForUser(ctx, d.UserId)
	if err != nil {
		return nil, err
	}

	return &v1.MenusRes{List: toV1MenuTree(menuTree)}, nil
}

func toV1MenuTree(nodes []*bo.MenuTreeNode) []*v1.MenuTreeNode {
	out := make([]*v1.MenuTreeNode, 0, len(nodes))
	for _, node := range nodes {
		item := &v1.MenuTreeNode{
			Id:            node.Id,
			Name:          node.Name,
			Permission:    node.Permission,
			Type:          node.Type,
			Sort:          node.Sort,
			ParentId:      node.ParentId,
			Path:          node.Path,
			Icon:          node.Icon,
			Component:     node.Component,
			ComponentName: node.ComponentName,
			Status:        node.Status,
			Visible:       node.Visible,
			KeepAlive:     node.KeepAlive,
			AlwaysShow:    node.AlwaysShow,
			Children:      toV1MenuTree(node.Children),
		}
		out = append(out, item)
	}
	return out
}
