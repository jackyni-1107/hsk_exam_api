package sysmenu

import (
	"testing"

	"exam/internal/consts"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
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

func TestCollectMenuSubtreeIDs(t *testing.T) {
	menus := []sysentity.SysMenu{
		{Id: 1, ParentId: 0},
		{Id: 2, ParentId: 1},
		{Id: 3, ParentId: 1},
		{Id: 4, ParentId: 2},
		{Id: 5, ParentId: 0},
	}

	subtree := collectMenuSubtreeIDs(menus, 1)
	if len(subtree) != 4 {
		t.Fatalf("unexpected subtree length: %d", len(subtree))
	}

	idSet := make(map[int64]struct{}, len(subtree))
	for _, id := range subtree {
		idSet[id] = struct{}{}
	}
	for _, expectedID := range []int64{1, 2, 3, 4} {
		if _, ok := idSet[expectedID]; !ok {
			t.Fatalf("missing subtree id: %d", expectedID)
		}
	}
	if _, ok := idSet[5]; ok {
		t.Fatalf("unexpected unrelated id in subtree: %v", subtree)
	}
}

func TestBuildMenuUpdateData(t *testing.T) {
	before := menuState{
		Name:      "old",
		Sort:      1,
		Visible:   true,
		KeepAlive: false,
	}
	after := menuState{
		Name:      "",
		Sort:      0,
		Visible:   false,
		KeepAlive: true,
	}
	data := buildMenuUpdateData(before, after, "admin")

	if data.Updater != "admin" {
		t.Fatalf("unexpected updater: %+v", data)
	}
	if data.Name != "" {
		t.Fatalf("expected explicit empty name to be preserved: %+v", data)
	}
	if data.Sort != 0 {
		t.Fatalf("expected explicit zero sort to be preserved: %+v", data)
	}
	if data.Visible != 0 {
		t.Fatalf("expected false visible to map to 0: %+v", data)
	}
	if data.KeepAlive != 1 {
		t.Fatalf("expected true keep_alive to map to 1: %+v", data)
	}
	if data.Permission != nil || data.Path != nil || data.Status != nil {
		t.Fatalf("expected omitted fields to stay nil: %+v", data)
	}
}

func TestNewMenuCreateStateNormalizesButtonFields(t *testing.T) {
	state := newMenuCreateState(
		" Save ",
		" user:update ",
		" /save ",
		" Edit ",
		" views/system/user/index ",
		" UserPage ",
		consts.MenuTypeButton,
		1,
		2,
		consts.StatusNormal,
		1,
		1,
		1,
	)

	if state.Name != "Save" || state.Permission != "user:update" {
		t.Fatalf("unexpected trimmed state: %+v", state)
	}
	if state.Path != "" || state.Icon != "" || state.Component != "" || state.ComponentName != "" {
		t.Fatalf("expected button-only fields to be cleared: %+v", state)
	}
	if state.KeepAlive || state.AlwaysShow {
		t.Fatalf("expected button-only flags to be reset: %+v", state)
	}
}

func TestBuildMenuUpdateDataClearsButtonFields(t *testing.T) {
	before := menuState{
		Type:          consts.MenuTypeMenu,
		Path:          "/system/user",
		Icon:          "User",
		Component:     "views/system/user/index",
		ComponentName: "SystemUser",
		KeepAlive:     true,
		AlwaysShow:    true,
	}
	after := sanitizeMenuState(menuState{
		Type:          consts.MenuTypeButton,
		Path:          before.Path,
		Icon:          before.Icon,
		Component:     before.Component,
		ComponentName: before.ComponentName,
		KeepAlive:     before.KeepAlive,
		AlwaysShow:    before.AlwaysShow,
	})

	data := buildMenuUpdateData(before, after, "admin")
	if data.Type != consts.MenuTypeButton {
		t.Fatalf("expected button type update: %+v", data)
	}
	if data.Path != "" || data.Icon != "" || data.Component != "" || data.ComponentName != "" {
		t.Fatalf("expected irrelevant button fields to be cleared: %+v", data)
	}
	if data.KeepAlive != 0 || data.AlwaysShow != 0 {
		t.Fatalf("expected button flags to be reset: %+v", data)
	}
}

func TestValidateMenuStructure(t *testing.T) {
	menus := []sysentity.SysMenu{
		{Id: 1, ParentId: 0, Type: consts.MenuTypeDir},
		{Id: 2, ParentId: 1, Type: consts.MenuTypeMenu},
		{Id: 3, ParentId: 2, Type: consts.MenuTypeButton},
	}

	cases := []struct {
		name         string
		currentMenu  int64
		state        menuState
		expectedCode int
		wantErr      bool
	}{
		{
			name:        "valid root menu",
			currentMenu: 0,
			state: menuState{
				Name:     "system",
				Type:     consts.MenuTypeDir,
				Sort:     0,
				Status:   consts.StatusNormal,
				ParentID: 0,
			},
		},
		{
			name:        "invalid type",
			currentMenu: 0,
			state: menuState{
				Name:     "invalid",
				Type:     99,
				Sort:     0,
				Status:   consts.StatusNormal,
				ParentID: 0,
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "invalid status",
			currentMenu: 0,
			state: menuState{
				Name:      "menu",
				Type:      consts.MenuTypeMenu,
				Sort:      0,
				Status:    99,
				ParentID:  0,
				Path:      "/menu",
				Component: "views/system/menu/index",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "missing parent",
			currentMenu: 0,
			state: menuState{
				Name:      "menu",
				Type:      consts.MenuTypeMenu,
				Sort:      0,
				Status:    consts.StatusNormal,
				ParentID:  99,
				Path:      "/menu",
				Component: "views/system/menu/index",
			},
			expectedCode: consts.CodeMenuNotFound.Code(),
			wantErr:      true,
		},
		{
			name:        "button cannot be parent",
			currentMenu: 0,
			state: menuState{
				Name:      "menu",
				Type:      consts.MenuTypeMenu,
				Sort:      0,
				Status:    consts.StatusNormal,
				ParentID:  3,
				Path:      "/menu",
				Component: "views/system/menu/index",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "cannot parent self",
			currentMenu: 2,
			state: menuState{
				Name:      "menu",
				Type:      consts.MenuTypeMenu,
				Sort:      0,
				Status:    consts.StatusNormal,
				ParentID:  2,
				Path:      "/menu",
				Component: "views/system/menu/index",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "cannot move under descendant",
			currentMenu: 1,
			state: menuState{
				Name:     "dir",
				Type:     consts.MenuTypeDir,
				Sort:     0,
				Status:   consts.StatusNormal,
				ParentID: 2,
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "directory cannot carry permission",
			currentMenu: 0,
			state: menuState{
				Name:       "system",
				Type:       consts.MenuTypeDir,
				Sort:       0,
				Status:     consts.StatusNormal,
				ParentID:   0,
				Permission: "user:list",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "menu requires path",
			currentMenu: 0,
			state: menuState{
				Name:      "user",
				Type:      consts.MenuTypeMenu,
				Sort:      0,
				Status:    consts.StatusNormal,
				ParentID:  1,
				Component: "views/system/user/index",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "menu requires component",
			currentMenu: 0,
			state: menuState{
				Name:     "user",
				Type:     consts.MenuTypeMenu,
				Sort:     0,
				Status:   consts.StatusNormal,
				ParentID: 1,
				Path:     "/system/user",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "keep alive requires component name",
			currentMenu: 0,
			state: menuState{
				Name:      "user",
				Type:      consts.MenuTypeMenu,
				Sort:      0,
				Status:    consts.StatusNormal,
				ParentID:  1,
				Path:      "/system/user",
				Component: "views/system/user/index",
				KeepAlive: true,
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "button requires permission",
			currentMenu: 0,
			state: menuState{
				Name:     "user:create",
				Type:     consts.MenuTypeButton,
				Sort:     0,
				Status:   consts.StatusNormal,
				ParentID: 2,
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
		{
			name:        "button cannot be root",
			currentMenu: 0,
			state: menuState{
				Name:       "user:create",
				Type:       consts.MenuTypeButton,
				Sort:       0,
				Status:     consts.StatusNormal,
				ParentID:   0,
				Permission: "user:create",
			},
			expectedCode: consts.CodeInvalidParams.Code(),
			wantErr:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMenuStructure(menus, tc.currentMenu, tc.state)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				code := gerror.Code(err)
				if code == nil || code.Code() != tc.expectedCode {
					t.Fatalf("unexpected error code: got=%v want=%d", code, tc.expectedCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
