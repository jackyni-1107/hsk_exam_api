package sysmenu

import (
	"testing"

	"exam/internal/consts"
	sysentity "exam/internal/model/entity/sys"
)

func TestFilterActiveMenus(t *testing.T) {
	menus := []sysentity.SysMenu{
		{Id: 1, Status: consts.StatusNormal},
		{Id: 2, Status: consts.StatusDisabled},
		{Id: 3, Status: consts.StatusNormal},
	}

	active := filterActiveMenus(menus)
	if len(active) != 2 {
		t.Fatalf("unexpected active menu count: %d", len(active))
	}
	if active[0].Id != 1 || active[1].Id != 3 {
		t.Fatalf("unexpected active menus: %+v", active)
	}
}

func TestFilterVisibleMenusIncludesAncestors(t *testing.T) {
	menus := []sysentity.SysMenu{
		{Id: 1, ParentId: 0, Status: consts.StatusNormal},
		{Id: 2, ParentId: 1, Status: consts.StatusNormal},
		{Id: 3, ParentId: 2, Status: consts.StatusNormal},
		{Id: 4, ParentId: 0, Status: consts.StatusNormal},
	}

	filtered := filterVisibleMenus(menus, map[int64]struct{}{
		3: {},
	})
	if len(filtered) != 3 {
		t.Fatalf("unexpected visible menu count: %d", len(filtered))
	}
	if filtered[0].Id != 1 || filtered[1].Id != 2 || filtered[2].Id != 3 {
		t.Fatalf("unexpected visible menus: %+v", filtered)
	}
}

func TestFilterVisibleMenusSkipsMissingOrUnallowedNodes(t *testing.T) {
	menus := []sysentity.SysMenu{
		{Id: 1, ParentId: 0, Status: consts.StatusNormal},
		{Id: 2, ParentId: 1, Status: consts.StatusNormal},
	}

	filtered := filterVisibleMenus(menus, map[int64]struct{}{
		99: {},
	})
	if len(filtered) != 0 {
		t.Fatalf("unexpected visible menus: %+v", filtered)
	}
}

func TestBuildMenuTreeNodes(t *testing.T) {
	menus := []sysentity.SysMenu{
		{Id: 3, ParentId: 1, Sort: 2, Name: "child-b"},
		{Id: 1, ParentId: 0, Sort: 2, Name: "root-b"},
		{Id: 2, ParentId: 0, Sort: 1, Name: "root-a"},
		{Id: 4, ParentId: 1, Sort: 1, Name: "child-a"},
	}

	tree := buildMenuTreeNodes(menus, 0)
	if len(tree) != 2 {
		t.Fatalf("unexpected root node count: %d", len(tree))
	}
	if tree[0].Id != 2 || tree[1].Id != 1 {
		t.Fatalf("unexpected root order: %+v", tree)
	}
	if len(tree[1].Children) != 2 {
		t.Fatalf("unexpected child count: %d", len(tree[1].Children))
	}
	if tree[1].Children[0].Id != 4 || tree[1].Children[1].Id != 3 {
		t.Fatalf("unexpected child order: %+v", tree[1].Children)
	}
}
