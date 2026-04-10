package menu

import menusvc "exam/internal/service/menu"

func init() {
	menusvc.RegisterMenu(new(sMenu))
}
