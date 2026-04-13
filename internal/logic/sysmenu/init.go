package sysmenu

import "exam/internal/service/sysmenu"

type sSysMenu struct{}

func init() {
	sysmenu.RegisterSysMenu(New())
}

func New() *sSysMenu {
	return &sSysMenu{}
}
