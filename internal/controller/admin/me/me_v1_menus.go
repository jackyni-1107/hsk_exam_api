package me

import (
	"context"
	"sort"

	v1 "exam/api/admin/me/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	sysentity "exam/internal/model/entity/sys"
	menusvc "exam/internal/service/menu"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}

	all, err := menusvc.Menu().MenuTree(ctx)
	if err != nil {
		return nil, err
	}

	active := make([]sysentity.SysMenu, 0, len(all))
	for _, m := range all {
		if m.Status == consts.StatusNormal {
			active = append(active, m)
		}
	}

	byID := make(map[int64]sysentity.SysMenu, len(active))
	for _, m := range active {
		byID[m.Id] = m
	}

	var allowed map[int64]struct{}
	if d.UserId == consts.SuperAdminUserId {
		allowed = make(map[int64]struct{}, len(active))
		for _, m := range active {
			allowed[m.Id] = struct{}{}
		}
	} else {
		allowed, err = menusvc.Menu().MenuIdsForUser(ctx, d.UserId)
		if err != nil {
			return nil, err
		}
	}

	visible := make(map[int64]struct{})
	for id := range allowed {
		for id != 0 {
			if _, ok := visible[id]; ok {
				break
			}
			visible[id] = struct{}{}
			m, ok := byID[id]
			if !ok {
				break
			}
			id = m.ParentId
		}
	}

	filtered := make([]sysentity.SysMenu, 0, len(visible))
	for _, m := range active {
		if _, ok := visible[m.Id]; ok {
			filtered = append(filtered, m)
		}
	}

	return &v1.MenusRes{List: buildMenuTreeV1(filtered, 0)}, nil
}

func buildMenuTreeV1(list []sysentity.SysMenu, parentId int64) []*v1.MenuTreeNode {
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
