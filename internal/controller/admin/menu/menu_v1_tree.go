package menu

import (
	"context"
	"sort"

	v1 "exam/api/admin/menu/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MenuTree(ctx context.Context, req *v1.MenuTreeReq) (res *v1.MenuTreeRes, err error) {
	var all []sysentity.SysMenu
	err = dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("sort").OrderAsc("id").
		Scan(&all)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.MenuTreeRes{List: buildAdminMenuTree(all, 0)}, nil
}

func buildAdminMenuTree(list []sysentity.SysMenu, parentId int64) []*v1.MenuTreeNode {
	children := make(map[int64][]sysentity.SysMenu)
	for _, m := range list {
		pid := m.ParentId
		children[pid] = append(children[pid], m)
	}
	for pid := range children {
		sort.Slice(children[pid], func(i, j int) bool {
			a, b := children[pid][i], children[pid][j]
			if a.Sort != b.Sort {
				return a.Sort < b.Sort
			}
			return a.Id < b.Id
		})
	}
	var walk func(pid int64) []*v1.MenuTreeNode
	walk = func(pid int64) []*v1.MenuTreeNode {
		slice := children[pid]
		out := make([]*v1.MenuTreeNode, 0, len(slice))
		for _, m := range slice {
			n := &v1.MenuTreeNode{
				Id: m.Id, Name: m.Name, Permission: m.Permission, Type: m.Type, Sort: m.Sort,
				ParentId: m.ParentId, Path: m.Path, Icon: m.Icon, Component: m.Component,
				ComponentName: m.ComponentName, Status: m.Status, Visible: m.Visible,
				KeepAlive: m.KeepAlive, AlwaysShow: m.AlwaysShow,
			}
			n.Children = walk(m.Id)
			out = append(out, n)
		}
		return out
	}
	return walk(parentId)
}
