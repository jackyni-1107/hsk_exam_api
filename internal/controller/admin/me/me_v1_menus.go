package me

import (
	"context"
	"sort"

	v1 "exam/api/admin/me/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) Menus(ctx context.Context, req *v1.MenusReq) (res *v1.MenusRes, err error) {
	d := middleware.GetCtxData(ctx)
	if d == nil {
		return nil, gerror.NewCode(consts.CodeTokenRequired)
	}
	var all []sysentity.SysMenu
	err = dao.SystemMenu.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		OrderAsc("sort").OrderAsc("id").
		Scan(&all)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	byID := make(map[int64]sysentity.SysMenu, len(all))
	for _, m := range all {
		byID[m.Id] = m
	}
	var allowed map[int64]struct{}
	if d.UserId == consts.SuperAdminUserId {
		allowed = make(map[int64]struct{}, len(all))
		for _, m := range all {
			allowed[m.Id] = struct{}{}
		}
	} else {
		allowed, err = menuIdsForUser(ctx, d.UserId)
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
	for _, m := range all {
		if _, ok := visible[m.Id]; ok {
			filtered = append(filtered, m)
		}
	}
	return &v1.MenusRes{List: buildMenuTreeV1(filtered, 0)}, nil
}

func menuIdsForUser(ctx context.Context, userId int64) (map[int64]struct{}, error) {
	var userRoles []sysentity.SysUserRole
	if err := dao.SystemUserRole.Ctx(ctx).Where("user_id", userId).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&userRoles); err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if len(userRoles) == 0 {
		return map[int64]struct{}{}, nil
	}
	roleIds := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIds = append(roleIds, ur.RoleId)
	}
	var roleMenus []sysentity.SysRoleMenu
	if err := dao.SystemRoleMenu.Ctx(ctx).WhereIn("role_id", roleIds).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&roleMenus); err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	out := make(map[int64]struct{}, len(roleMenus))
	for _, rm := range roleMenus {
		out[rm.MenuId] = struct{}{}
	}
	return out, nil
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
