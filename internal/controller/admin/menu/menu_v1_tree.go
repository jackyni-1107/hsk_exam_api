package menu

import (
	"context"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/model/bo"
	menusvc "exam/internal/service/sysmenu"
)

func (c *ControllerV1) MenuTree(ctx context.Context, req *v1.MenuTreeReq) (res *v1.MenuTreeRes, err error) {
	menuTree, err := menusvc.SysMenu().MenuTreeNodes(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MenuTreeRes{List: toAdminMenuTree(menuTree)}, nil
}

func toAdminMenuTree(nodes []*bo.MenuTreeNode) []*v1.MenuTreeNode {
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
			Children:      toAdminMenuTree(node.Children),
		}
		out = append(out, item)
	}
	return out
}
